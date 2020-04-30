package faunadb

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
)

var (
	exprType = reflect.TypeOf((*Expr)(nil)).Elem()
	objType  = reflect.TypeOf((*Obj)(nil)).Elem()
	arrType  = reflect.TypeOf((*Arr)(nil)).Elem()
	timeType = reflect.TypeOf((*time.Time)(nil)).Elem()

	maxSupportedUint = uint64(math.MaxInt64)

	errMapKeyMustBeString       = invalidExpr{errors.New("Error while encoding map to json: All map keys must be of type string")}
	errMaxSupportedUintExceeded = invalidExpr{errors.New("Error while encoding number to json: Uint value exceeds maximum int64")}
)

func wrap(i interface{}) Expr {
	if i == nil {
		return NullV{}
	}

	value, valueType := indirectValue(i)
	kind := value.Kind()

	if (kind == reflect.Ptr || kind == reflect.Interface) && value.IsNil() {
		return NullV{}
	}

	// Is an expression but not a syntax sugar
	if valueType.Implements(exprType) && valueType != objType && valueType != arrType {
		return value.Interface().(Expr)
	}

	switch kind {
	case reflect.String:
		return StringV(value.String())

	case reflect.Bool:
		return BooleanV(value.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return LongV(value.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num := value.Uint()

		if num > maxSupportedUint {
			return errMaxSupportedUintExceeded
		}

		return LongV(num)

	case reflect.Float32, reflect.Float64:
		return DoubleV(value.Float())

	case reflect.Map:
		if valueType.Key().Kind() != reflect.String {
			return errMapKeyMustBeString
		}

		return wrapMap(value)

	case reflect.Struct:
		if valueType == timeType {
			return TimeV(value.Interface().(time.Time))
		}

		value, _ = indirectValue(structToMap(value))
		return wrapMap(value)

	case reflect.Slice, reflect.Array:
		return wrapArray(value)

	default:
		return invalidExpr{fmt.Errorf("Error while converting Expr to JSON: Non supported type %v", kind)}
	}
}

func wrapMap(value reflect.Value) Expr {
	obj := make(unescapedObj, value.Len())

	for _, key := range value.MapKeys() {
		obj[key.String()] = wrap(value.MapIndex(key))
	}

	return unescapedObj{"object": obj}
}

func wrapArray(value reflect.Value) Expr {
	arr := make(unescapedArr, value.Len())

	for i, size := 0, value.Len(); i < size; i++ {
		arr[i] = wrap(value.Index(i))
	}

	return arr
}
