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

func TestSerializeUpdate(t *testing.T) {
	json, err := toJSON(
		Update(Ref("classes/spells/123"), Obj{
			"name": "fire",
		}),
	)

	require.NoError(t, err)
	require.Equal(t, `{"params":{"object":{"name":"fire"}},"update":{"@ref":"classes/spells/123"}}`, json)
}

func TestSerializeReplace(t *testing.T) {
	json, err := toJSON(
		Replace(Ref("classes/spells/123"), Obj{
			"name": "fire",
		}),
	)

	require.NoError(t, err)
	require.Equal(t, `{"params":{"object":{"name":"fire"}},"replace":{"@ref":"classes/spells/123"}}`, json)
}

func TestSerializeDelete(t *testing.T) {
	json, err := toJSON(
		Delete(Ref("classes/spells/123")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"delete":{"@ref":"classes/spells/123"}}`, json)
}

func TestSerializeGet(t *testing.T) {
	json, err := toJSON(
		Get(Ref("classes/spells/42")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"get":{"@ref":"classes/spells/42"}}`, json)
}

func TestSerializeNull(t *testing.T) {
	json, err := toJSON(Null())

	require.NoError(t, err)
	require.Equal(t, `null`, json)
}

func TestSerializeNullOnObject(t *testing.T) {
	json, err := toJSON(Obj{"data": Null()})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"data":null}}`, json)
}

func TestSerializeNullOnStruct(t *testing.T) {
	type structWithNull struct {
		Null *string
	}

	json, err := toJSON(Obj{"data": structWithNull{nil}})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"data":{"object":{"Null":null}}}}`, json)
}

func TestSerializeLet(t *testing.T) {
	json, err := toJSON(
		Let(
			Obj{"v1": Ref("classes/spells/42")},
			Exists(Var("v1")),
		),
	)

	require.NoError(t, err)
	require.Equal(t, `{"in":{"exists":{"var":"v1"}},"let":{"v1":{"@ref":"classes/spells/42"}}}`, json)
}

func TestSerializeIf(t *testing.T) {
	json, err := toJSON(If(true, "exists", "does not exists"))

	require.NoError(t, err)
	require.Equal(t, `{"else":"does not exists","if":true,"then":"exists"}`, json)
}

func TestSerializeDo(t *testing.T) {
	json, err := toJSON(Do(Arr{
		Get(Ref("classes/spells/4")),
		Get(Ref("classes/spells/2")),
	}))

	require.NoError(t, err)
	require.Equal(t, `{"do":[{"get":{"@ref":"classes/spells/4"}},{"get":{"@ref":"classes/spells/2"}}]}`, json)
}

func TestSerializeDoWithVarargs(t *testing.T) {
	json, err := toJSON(Do(
		Get(Ref("classes/spells/4")),
		Get(Ref("classes/spells/2")),
	))

	require.NoError(t, err)
	require.Equal(t, `{"do":[{"get":{"@ref":"classes/spells/4"}},{"get":{"@ref":"classes/spells/2"}}]}`, json)
}

func TestSerializeLambda(t *testing.T) {
	json, err := toJSON(Lambda("x", Var("x")))

	require.NoError(t, err)
	require.Equal(t, `{"expr":{"var":"x"},"lambda":"x"}`, json)
}

func TestSerializeMap(t *testing.T) {
	json, err := toJSON(Map(Arr{1, 2, 3}, Lambda("x", Var("x"))))

	require.NoError(t, err)
	require.Equal(t, `{"collection":[1,2,3],"map":{"expr":{"var":"x"},"lambda":"x"}}`, json)
}

func TestSerializeForeach(t *testing.T) {
	json, err := toJSON(Foreach(Arr{1, 2, 3}, Lambda("x", Var("x"))))

	require.NoError(t, err)
	require.Equal(t, `{"collection":[1,2,3],"foreach":{"expr":{"var":"x"},"lambda":"x"}}`, json)
}

func TestSerializeFilter(t *testing.T) {
	json, err := toJSON(Filter(Arr{true, false}, Lambda("x", Var("x"))))

	require.NoError(t, err)
	require.Equal(t, `{"collection":[true,false],"filter":{"expr":{"var":"x"},"lambda":"x"}}`, json)
}

func TestSerializeTake(t *testing.T) {
	json, err := toJSON(Take(2, Arr{1, 2, 3}))

	require.NoError(t, err)
	require.Equal(t, `{"collection":[1,2,3],"take":2}`, json)
}

func TestSerializeExists(t *testing.T) {
	json, err := toJSON(
		Exists(Ref("classes/spells/42")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"exists":{"@ref":"classes/spells/42"}}`, json)
}

func toJSON(expr Expr) (string, error) {
	bytes, err := writeJSON(expr)
	return string(bytes), err
}
