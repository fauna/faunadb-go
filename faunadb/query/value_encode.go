package query

import (
	"encoding/json"
	"reflect"
)

type wrapped struct {
	value interface{}
}

func (w wrapped) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.value)
}

func wrap(expr Expr) wrapped {
	switch v := expr.(type) {
	case wrapped:
		return v
	case fn:
		return wrapFn(v)
	default:
		return reflectWrap(v)
	}
}

func wrapFn(fn fn) wrapped {
	call := make(map[string]wrapped)

	for key, value := range fn {
		call[key] = wrap(value)
	}

	return wrapped{call}
}

func reflectWrap(expr Expr) wrapped {
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
		return wrapped{expr}
	}
}

//FIXME: Validate key is an string
func wrapMap(value reflect.Value) wrapped {
	obj := make(map[string]wrapped)

	for _, key := range value.MapKeys() {
		obj[key.String()] = wrap(value.MapIndex(key).Interface())
	}

	return wrapped{map[string]Expr{"object": obj}}
}

func wrapArr(value reflect.Value) wrapped {
	var arr []wrapped

	for i, size := 0, value.Len(); i < size; i++ {
		arr = append(arr, wrap(value.Index(i).Interface()))
	}

	return wrapped{arr}
}

func wrapStruct(value reflect.Value) wrapped {
	obj := make(Obj)

	for i, size := 0, value.NumField(); i < size; i++ {
		key := getKeyName(value.Type().Field(i))
		value := value.Field(i).Interface()

		obj[key] = wrap(value)
	}

	return wrapped{obj}
}

func getKeyName(field reflect.StructField) string {
	if key := field.Tag.Get("fauna"); key != "" {
		return key
	}

	return field.Name
}
