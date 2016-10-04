package faunadb

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"
)

var (
	exprType              = reflect.TypeOf((*Expr)(nil)).Elem()
	timeType              = reflect.TypeOf((*time.Time)(nil)).Elem()
	errMapKeyMustBeString = errors.New("Error while encoding map to json: All map keys must be of type string")
)

func writeJSON(expr Expr) (bytes []byte, err error) {
	var escaped interface{}

	if escaped, err = expr.toJSON(); err == nil {
		bytes, err = json.Marshal(escaped)
	}

	return
}

func escapeValue(any interface{}) (interface{}, error) {
	value, valueType := indirectValue(any)

	if valueType.Implements(exprType) {
		return value.Interface().(Expr).toJSON()
	}

	if valueType == timeType {
		return TimeV(value.Interface().(time.Time)).toJSON()
	}

	switch value.Kind() {
	case reflect.Map:
		if valueType.Key().Kind() == reflect.String {
			return escapeMap(value)
		}
		return nil, errMapKeyMustBeString

	case reflect.Struct:
		return escapeMap(structToMap(value))
	case reflect.Slice:
		return escapeArray(value)
	default:
		return value.Interface(), nil
	}
}

func escapeMap(obj interface{}) (interface{}, error) {
	value, _ := indirectValue(obj)

	res := make(map[string]interface{}, value.Len())

	for _, key := range value.MapKeys() {
		escaped, err := escapeValue(value.MapIndex(key))

		if err != nil {
			return nil, err
		}

		res[key.String()] = escaped
	}

	return map[string]interface{}{"object": res}, nil
}

func escapeArray(arr interface{}) (interface{}, error) {
	value, _ := indirectValue(arr)

	res := make([]interface{}, value.Len())

	for i, size := 0, value.Len(); i < size; i++ {
		escaped, err := escapeValue(value.Index(i))

		if err != nil {
			return nil, err
		}

		res[i] = escaped
	}

	return res, nil
}
