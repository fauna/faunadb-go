package query

import (
	"encoding/json"
	"testing"
)

func TestSerializeObject(t *testing.T) {
	assertJson(t,
		Obj{"key": "value"},
		`{"object":{"key":"value"}}`,
	)

	assertJson(t,
		Obj{"key": Obj{"nested": "value"}},
		`{"object":{"key":{"object":{"nested":"value"}}}}`,
	)

	assertJson(t,
		Obj{"key": map[string]string{"nested": "value"}},
		`{"object":{"key":{"object":{"nested":"value"}}}}`,
	)
}

func TestSerializeArray(t *testing.T) {
	assertJson(t,
		Arr{map[string]string{"key": "value"}},
		`[{"object":{"key":"value"}}]`,
	)

	assertJson(t,
		Arr{Obj{"key": "value"}},
		`[{"object":{"key":"value"}}]`,
	)
}

func TestSerializeStruct(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}

	assertJson(t,
		Obj{"data": user{"Jhon", 42}},
		`{"object":{"data":{"object":{"Age":42,"Name":"Jhon"}}}}`,
	)
}

func TestSerializeStructWithTags(t *testing.T) {
	type user struct {
		Name string `fauna:"name"`
		Age  int    `fauna:"age"`
	}

	assertJson(t,
		Obj{"data": user{"Jhon", 42}},
		`{"object":{"data":{"object":{"age":42,"name":"Jhon"}}}}`,
	)
}

func TestSerializeStructWithPointers(t *testing.T) {
	type user struct {
		Name string
		Age  *int
	}

	age := 42

	assertJson(t,
		Obj{"data": &user{"Jhon", &age}},
		`{"object":{"data":{"object":{"Age":42,"Name":"Jhon"}}}}`,
	)
}

func TestSerializeStructWithNestedExpressions(t *testing.T) {
	type user struct {
		Name string
	}

	type userInfo struct {
		User        user
		Credentials map[string]string
	}

	assertJson(t,
		Obj{"data": userInfo{user{"Jhon"}, map[string]string{"password": "1234"}}},
		`{"object":{"data":{"object":{"Credentials":{"object":{"password":"1234"}},"User":{"object":{"Name":"Jhon"}}}}}}`,
	)
}

func TestSerializeCreate(t *testing.T) {
	assertJson(t,
		Create(Ref("classes/spells"), Obj{
			"name": "fire",
		}),
		`{"create":{"@ref":"classes/spells"},"params":{"object":{"name":"fire"}}}`,
	)
}

func TestSerializeGet(t *testing.T) {
	assertJson(t,
		Get(Ref("classes/spells/42")),
		`{"get":{"@ref":"classes/spells/42"}}`,
	)
}

func TestSerializeDelete(t *testing.T) {
	assertJson(t,
		Delete(Ref("classes/spells/42")),
		`{"delete":{"@ref":"classes/spells/42"}}`,
	)
}

func assertJson(t *testing.T, expr interface{}, expected string) {
	if json, err := json.Marshal(wrap(expr)); err == nil {
		actual := string(json)

		if expected != actual {
			t.Errorf("\n%10s: %#v\n%10s: %#v", "Expected", expected, "got", actual)
		}
	} else {
		t.Error(err)
	}
}
