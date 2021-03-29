package faunadb

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const faunaTag = "fauna"

func fieldName(field reflect.StructField) string {
	name := field.Tag.Get(faunaTag)

	if name == "" {
		name = field.Name
	}

	return name
}

// parseTag interprets fauna struct field tags
func parseTag(field reflect.StructField) (name string, ignore, omitempty bool, err error) {
	s := field.Tag.Get(faunaTag)
	parts := strings.Split(s, ",")
	if s == "" {
		return field.Name, false, false, nil
	}

	if parts[0] == "-" {
		return "", true, false, nil
	}

	if len(parts) > 1 {
		for _, p := range parts[1:] {
			switch p {
			case "omitempty":
				omitempty = true
			default:
				err = fmt.Errorf("fauna: struct tag has invalid option: %q", p)
				return "", false, false, err
			}
		}
	}
	if parts[0] != "" {
		name = parts[0]
	} else {
		name = field.Name
	}

	return
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			return t.IsZero()
		}
	}
	return false
}
