package faunadb

import "reflect"

func structToMap(aStruct reflect.Value) map[string]interface{} {
	res := make(map[string]interface{}, aStruct.NumField())

	for key, value := range exportedStructFields(aStruct) {
		res[key] = value.Interface()
	}

	return res
}

func exportedStructFields(aStruct reflect.Value) map[string]reflect.Value {
	fields := make(map[string]reflect.Value)
	aStructType := aStruct.Type()

	for i, size := 0, aStruct.NumField(); i < size; i++ {
		field := aStruct.Field(i)

		if !field.CanInterface() {
			continue
		}

		fields[fieldName(aStructType.Field(i))] = field
	}

	return fields
}

func indirectValue(i interface{}) (reflect.Value, reflect.Type) {
	var value reflect.Value

	if reflected, ok := i.(reflect.Value); ok {
		value = reflected
	} else {
		value = reflect.ValueOf(i)
	}

	for {
		if value.Kind() == reflect.Interface {
			elem := value.Elem()

			if elem.IsValid() {
				value = elem
				continue
			}
		}

		if value.Kind() != reflect.Ptr {
			break
		}

		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}

		value = value.Elem()
	}

	return value, value.Type()
}
