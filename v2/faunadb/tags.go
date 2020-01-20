package faunadb

import "reflect"

const faunaTag = "fauna"

func fieldName(field reflect.StructField) string {
	name := field.Tag.Get(faunaTag)

	if name == "" {
		name = field.Name
	}

	return name
}
