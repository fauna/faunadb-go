package faunadb

import "reflect"

func structToMap(aStruct reflect.Value) map[string]interface{} {
	res := make(map[string]interface{}, aStruct.NumField())

	if expStructFields, err := exportedStructFields(aStruct); err == nil {
		for key, value := range expStructFields {
			res[key] = value.Interface()
		}
	}
	return res
}

func exportedStructFields(aStruct reflect.Value) (map[string]reflect.Value, error) {
	fields := make(map[string]reflect.Value)
	aStructType := aStruct.Type()

	for i, size := 0, aStruct.NumField(); i < size; i++ {
		field := aStruct.Field(i)

		if !field.CanInterface() {
			continue
		}

		fieldName, ignore, omitempty, err := parseTag(aStructType.Field(i))
		if err != nil {
			return nil, err
		}

		if omitempty && isEmptyValue(field) {
			continue
		}

		if !ignore && fieldName != "" {
			fields[fieldName] = field
		}
	}

	return fields, nil
}

func indirectValue(i interface{}) (reflect.Value, reflect.Type) {
	var value reflect.Value

	if reflected, ok := i.(reflect.Value); ok {
		value = reflected
	} else {
		value = reflect.ValueOf(i)
	}

	for {
		if value.Kind() == reflect.Interface && !value.IsNil() {
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
			if value.CanSet() {
				value.Set(reflect.New(value.Type().Elem()))
			} else {
				break
			}
		}

		value = value.Elem()
	}

	return value, value.Type()
}

func allStructFields(aStruct reflect.Value) map[string]reflect.Value {
	fields := make(map[string]reflect.Value)
	aStructType := aStruct.Type()

	for i, size := 0, aStruct.NumField(); i < size; i++ {
		field := aStruct.Field(i)
		structTypeField := aStructType.Field(i)

		if !field.CanInterface() {
			continue
		}

		fieldName, ignore, _, _ := parseTag(aStructType.Field(i))

		if !ignore && fieldName != "" {
			fields[fieldName] = field
		}
		fields[structTypeField.Name] = field
	}

	return fields
}
