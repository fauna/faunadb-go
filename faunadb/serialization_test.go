package faunadb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerializeObject(t *testing.T) {
	json, err := toJSON(Obj{"key": "value"})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"key":"value"}}`, json)
}

func TestSerializeNestedObjects(t *testing.T) {
	json, err := toJSON(Obj{"key": Obj{"nested": "value"}})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"key":{"object":{"nested":"value"}}}}`, json)
}

func TestSerializeNestedMaps(t *testing.T) {
	json, err := toJSON(Obj{"key": map[string]string{"nested": "value"}})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"key":{"object":{"nested":"value"}}}}`, json)
}

func TestSerializeInvalidMaps(t *testing.T) {
	_, err := toJSON(Obj{"key": map[int]string{1: "value"}})
	require.EqualError(t, err, "Error while encoding map to json: All map keys must be of type string")
}

func TestSerializeArray(t *testing.T) {
	json, err := toJSON(Arr{1, 2, 3})

	require.NoError(t, err)
	require.Equal(t, `[1,2,3]`, json)
}

func TestSerializeWithNestedArrays(t *testing.T) {
	json, err := toJSON(Arr{Arr{1, 2, 3}})

	require.NoError(t, err)
	require.Equal(t, `[[1,2,3]]`, json)
}

func TestSerializeStruct(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}

	json, err := toJSON(Obj{"data": user{"Jhon", 42}})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"data":{"object":{"Age":42,"Name":"Jhon"}}}}`, json)
}

func TestSerializeWithNonExportedFields(t *testing.T) {
	type user struct {
		Name string
		age  int
	}

	json, err := toJSON(Obj{"data": user{"Jhon", 42}})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"data":{"object":{"Name":"Jhon"}}}}`, json)
}

func TestSerializeStructWithTags(t *testing.T) {
	type user struct {
		Name string `fauna:"name"`
		Age  int    `fauna:"age"`
	}

	json, err := toJSON(Obj{"data": user{"Jhon", 42}})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"data":{"object":{"age":42,"name":"Jhon"}}}}`, json)
}

func TestSerializeStructWithPointers(t *testing.T) {
	type user struct {
		Name string
		Age  *int
	}

	age := 42

	json, err := toJSON(Obj{"data": &user{"Jhon", &age}})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"data":{"object":{"Age":42,"Name":"Jhon"}}}}`, json)
}

func TestSerializeStructWithNestedExpressions(t *testing.T) {
	type user struct {
		Name string
	}

	type userInfo struct {
		User        user
		Credentials map[string]string
	}

	json, err := toJSON(Obj{"data": userInfo{user{"Jhon"}, map[string]string{"password": "1234"}}})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"data":{"object":{"Credentials":{"object":{"password":"1234"}},"User":{"object":{"Name":"Jhon"}}}}}}`, json)
}

func TestSerializeStructWithEmbeddedStructs(t *testing.T) {
	type Embedded struct {
		Str string
	}

	type Data struct {
		Int int
		Embedded
	}

	json, err := toJSON(Obj{"data": Data{42, Embedded{"a string"}}})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"data":{"object":{"Embedded":{"object":{"Str":"a string"}},"Int":42}}}}`, json)
}

func TestSerializeCreate(t *testing.T) {
	json, err := toJSON(
		Create(Ref("classes/spells"), Obj{
			"name": "fire",
		}),
	)

	require.NoError(t, err)
	require.Equal(t, `{"create":{"@ref":"classes/spells"},"params":{"object":{"name":"fire"}}}`, json)
}

func TestSerializeGet(t *testing.T) {
	json, err := toJSON(
		Get(Ref("classes/spells/42")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"get":{"@ref":"classes/spells/42"}}`, json)
}

func TestSerializeDelete(t *testing.T) {
	json, err := toJSON(
		Delete(Ref("classes/spells/42")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"delete":{"@ref":"classes/spells/42"}}`, json)
}

func toJSON(expr Expr) (string, error) {
	bytes, err := writeJSON(expr)
	return string(bytes), err
}
