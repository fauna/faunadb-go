package query

import "reflect"

func wrap(value interface{}) Expr {
	switch v := value.(type) {
	case Expr:
		return v
	case fn:
		return wrapFn(v)
	default:
		return reflectWrap(v)
	}
}

func wrapFn(fn fn) Expr {
	call := make(map[string]Expr)

	for key, value := range fn {
		call[key] = wrap(value)
	}

	return Expr{call}
}

func reflectWrap(expr interface{}) Expr {
	value := reflect.Indirect(reflect.ValueOf(expr))

	switch value.Kind() {
	case reflect.Map:
		return wrapMap(value)
	case reflect.Slice:
		return wrapArr(value)
	case reflect.Struct:
		if value.Type().PkgPath() != "faunadb/values" {
			return wrapStruct(value)
		}
		fallthrough
	default:
		return Expr{expr}
	}
}

//FIXME: Validate key is an string
func wrapMap(value reflect.Value) Expr {
	obj := make(map[string]Expr)

	for _, key := range value.MapKeys() {
		obj[key.String()] = wrap(value.MapIndex(key).Interface())
	}

	return Expr{map[string]interface{}{"object": obj}}
}

func wrapArr(value reflect.Value) Expr {
	var arr []Expr

	for i, size := 0, value.Len(); i < size; i++ {
		arr = append(arr, wrap(value.Index(i).Interface()))
	}

	return Expr{arr}
}

func wrapStruct(value reflect.Value) Expr {
	obj := make(map[string]interface{})

	for i, size := 0, value.NumField(); i < size; i++ {
		key := getKeyName(value.Type().Field(i))
		value := value.Field(i).Interface()

		obj[key] = value
	}

	return wrap(obj)
}

func getKeyName(field reflect.StructField) string {
	if key := field.Tag.Get("fauna"); key != "" {
		return key
	}

	return field.Name
}
