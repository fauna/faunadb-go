package faunadb

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

var (
	exprType  = reflect.TypeOf((*Expr)(nil)).Elem()
	objType   = reflect.TypeOf((*Obj)(nil)).Elem()
	arrType   = reflect.TypeOf((*Arr)(nil)).Elem()
	timeType  = reflect.TypeOf((*time.Time)(nil)).Elem()
	intType   = reflect.TypeOf((*int64)(nil)).Elem()
	floatType = reflect.TypeOf((*float64)(nil)).Elem()

	errMapKeyMustBeString = invalidExpr{errors.New("Error while encoding map to json: All map keys must be of type string")}
)

func wrap(i interface{}) Expr {
	value, valueType := indirectValue(i)
	kind := value.Kind()

	if kind == reflect.Ptr && value.IsNil() {
		return NullV{}
	}

	// Is an expression but not a syntax sugar
	if valueType.Implements(exprType) && valueType != objType && valueType != arrType {
		return value.Interface().(Expr)
	}

	switch kind {
	case reflect.String:
		return StringV(value.Interface().(string))

	case reflect.Bool:
		return BooleanV(value.Interface().(bool))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return LongV(value.Convert(intType).Interface().(int64))

	case reflect.Float32, reflect.Float64:
		return DoubleV(value.Convert(floatType).Interface().(float64))

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

	case reflect.Slice:
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
