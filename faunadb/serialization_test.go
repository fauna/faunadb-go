package faunadb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSerializeObjectV(t *testing.T) {
	json, err := toJSON(ObjectV{
		"data": ObjectV{
			"name": StringV("test"),
		},
	})

	require.NoError(t, err)
	require.Equal(t, `{"object":{"data":{"object":{"name":"test"}}}}`, json)
}

func TestSerializeArrayV(t *testing.T) {
	json, err := toJSON(ArrayV{
		ObjectV{"name": StringV("a")},
		ObjectV{"name": StringV("b")},
	})

	require.NoError(t, err)
	require.Equal(t, `[{"object":{"name":"a"}},{"object":{"name":"b"}}]`, json)
}

func TestSerializeSetRefV(t *testing.T) {
	json, err := toJSON(SetRefV{
		ObjectV{"name": StringV("a")},
	})

	require.NoError(t, err)
	require.Equal(t, `{"@set":{"object":{"name":"a"}}}`, json)
}

func TestSerializeDateV(t *testing.T) {
	json, err := toJSON(DateV(time.Unix(0, 0).UTC()))

	require.NoError(t, err)
	require.Equal(t, `{"@date":"1970-01-01"}`, json)
}

func TestSerializeTimeV(t *testing.T) {
	json, err := toJSON(TimeV(time.Unix(1, 2).UTC()))

	require.NoError(t, err)
	require.Equal(t, `{"@ts":"1970-01-01T00:00:01.000000002Z"}`, json)
}

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

func TestSerializeInsert(t *testing.T) {
	json, err := toJSON(
		Insert(
			Ref("classes/spells/104979509696660483"),
			time.Unix(0, 0).UTC(),
			ActionCreate,
			Obj{"data": Obj{"name": "test"}},
		),
	)

	require.NoError(t, err)
	require.Equal(t, `{"action":"create","insert":{"@ref":"classes/spells/104979509696660483"},"params":{"object":{"data":{"object":{"name":"test"}}}},"ts":{"@ts":"1970-01-01T00:00:00Z"}}`, json)
}

func TestSerializeRemove(t *testing.T) {
	json, err := toJSON(
		Remove(
			Ref("classes/spells/104979509696660483"),
			time.Unix(0, 0).UTC(),
			ActionDelete,
		),
	)

	require.NoError(t, err)
	require.Equal(t, `{"action":"delete","remove":{"@ref":"classes/spells/104979509696660483"},"ts":{"@ts":"1970-01-01T00:00:00Z"}}`, json)
}

func TestSerializeCreateClass(t *testing.T) {
	json, err := toJSON(
		CreateClass(Obj{
			"name": "boons",
		}),
	)

	require.NoError(t, err)
	require.Equal(t, `{"create_class":{"object":{"name":"boons"}}}`, json)
}

func TestSerializeCreateDatabase(t *testing.T) {
	json, err := toJSON(
		CreateDatabase(Obj{
			"name": "db-next",
		}),
	)

	require.NoError(t, err)
	require.Equal(t, `{"create_database":{"object":{"name":"db-next"}}}`, json)
}

func TestSerializeCreateIndex(t *testing.T) {
	json, err := toJSON(
		CreateIndex(Obj{
			"name":   "new-index",
			"source": Ref("classes/spells"),
		}),
	)

	require.NoError(t, err)
	require.Equal(t, `{"create_index":{"object":{"name":"new-index","source":{"@ref":"classes/spells"}}}}`, json)
}

func TestSerializeCreateKey(t *testing.T) {
	json, err := toJSON(
		CreateKey(Obj{
			"database": Ref("databases/prydain"),
			"role":     "server",
		}),
	)

	require.NoError(t, err)
	require.Equal(t, `{"create_key":{"object":{"database":{"@ref":"databases/prydain"},"role":"server"}}}`, json)
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

func TestSerializeDrop(t *testing.T) {
	json, err := toJSON(Drop(2, Arr{1, 2, 3}))

	require.NoError(t, err)
	require.Equal(t, `{"collection":[1,2,3],"drop":2}`, json)
}

func TestSerializePrepend(t *testing.T) {
	json, err := toJSON(Prepend(Arr{1, 2, 3}, Arr{4, 5, 6}))

	require.NoError(t, err)
	require.Equal(t, `{"collection":[4,5,6],"prepend":[1,2,3]}`, json)
}

func TestSerializeAppend(t *testing.T) {
	json, err := toJSON(Append(Arr{4, 5, 6}, Arr{1, 2, 3}))

	require.NoError(t, err)
	require.Equal(t, `{"append":[4,5,6],"collection":[1,2,3]}`, json)
}

func TestSerializeGet(t *testing.T) {
	json, err := toJSON(
		Get(Ref("classes/spells/42")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"get":{"@ref":"classes/spells/42"}}`, json)
}

func TestSerializeGetWithTimestamp(t *testing.T) {
	json, err := toJSON(
		Get(
			Ref("classes/spells/42"),
			TS(time.Unix(0, 0).UTC()),
		),
	)

	require.NoError(t, err)
	require.Equal(t, `{"get":{"@ref":"classes/spells/42"},"ts":{"@ts":"1970-01-01T00:00:00Z"}}`, json)
}

func TestSerializeExists(t *testing.T) {
	json, err := toJSON(
		Exists(Ref("classes/spells/42")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"exists":{"@ref":"classes/spells/42"}}`, json)
}

func TestSerializeExistsWithTimestamp(t *testing.T) {
	json, err := toJSON(
		Exists(
			Ref("classes/spells/42"),
			TS(time.Unix(1, 1).UTC()),
		),
	)

	require.NoError(t, err)
	require.Equal(t, `{"exists":{"@ref":"classes/spells/42"},"ts":{"@ts":"1970-01-01T00:00:01.000000001Z"}}`, json)
}

func TestSerializeCount(t *testing.T) {
	json, err := toJSON(Count(Ref("databases")))

	require.NoError(t, err)
	require.Equal(t, `{"count":{"@ref":"databases"}}`, json)
}

func TestSerializeCountEvents(t *testing.T) {
	json, err := toJSON(
		Count(Ref("databases"), Events(true)),
	)

	require.NoError(t, err)
	require.Equal(t, `{"count":{"@ref":"databases"},"events":true}`, json)
}

func TestSerializePaginate(t *testing.T) {
	json, err := toJSON(
		Paginate(Ref("databases")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"paginate":{"@ref":"databases"}}`, json)
}

func TestSerializePaginateWithParameters(t *testing.T) {
	json, err := toJSON(
		Paginate(
			Ref("databases"),
			Before(Ref("databases/test10")),
			After(Ref("databases/test")),
			Events(true),
			Sources(true),
			TS(10),
			Size(2),
		),
	)

	require.NoError(t, err)
	require.Equal(t, `{"after":{"@ref":"databases/test"},"before":{"@ref":"databases/test10"},"events":true,"paginate":{"@ref":"databases"},"size":2,"sources":true,"ts":10}`, json)
}

func TestSerializeConcat(t *testing.T) {
	json, err := toJSON(
		Concat(Arr{"a", "b"}),
	)

	require.NoError(t, err)
	require.Equal(t, `{"concat":["a","b"]}`, json)
}

func TestSerializeConcatWithSeparator(t *testing.T) {
	json, err := toJSON(
		Concat(Arr{"a", "b"}, Separator("/")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"concat":["a","b"],"separator":"/"}`, json)
}

func TestSerializeCasefold(t *testing.T) {
	json, err := toJSON(
		Casefold("GET DOWN"),
	)

	require.NoError(t, err)
	require.Equal(t, `{"casefold":"GET DOWN"}`, json)
}

func TestSerializeTime(t *testing.T) {
	json, err := toJSON(
		Time("1970-01-01T00:00:00+00:00"),
	)

	require.NoError(t, err)
	require.Equal(t, `{"time":"1970-01-01T00:00:00+00:00"}`, json)
}

func TestSerializeEpoch(t *testing.T) {
	json, err := toJSON(Arr{
		Epoch(0, TimeUnitSecond),
		Epoch(0, TimeUnitMillisecond),
		Epoch(0, TimeUnitMicrosecond),
		Epoch(0, TimeUnitNanosecond),
	})

	require.NoError(t, err)
	require.Equal(t, `[{"epoch":0,"unit":"second"},{"epoch":0,"unit":"millisecond"},{"epoch":0,"unit":"microsecond"},{"epoch":0,"unit":"nanosecond"}]`, json)
}

func TestSerializeDate(t *testing.T) {
	json, err := toJSON(
		Date("1970-01-01"),
	)

	require.NoError(t, err)
	require.Equal(t, `{"date":"1970-01-01"}`, json)
}

func TestSerializeMatch(t *testing.T) {
	json, err := toJSON(
		Match(Ref("databases")),
	)

	require.NoError(t, err)
	require.Equal(t, `{"match":{"@ref":"databases"}}`, json)
}

func TestSerializeMatchWithTerms(t *testing.T) {
	json, err := toJSON(
		MatchTerm(
			Ref("indexes/spells_by_name"),
			"magic missile",
		),
	)

	require.NoError(t, err)
	require.Equal(t, `{"match":{"@ref":"indexes/spells_by_name"},"terms":"magic missile"}`, json)
}

func TestSerializeUnion(t *testing.T) {
	json, err := toJSON(
		Union(
			Ref("indexes/active_users"),
			Ref("indexes/vip_users"),
		),
	)

	require.NoError(t, err)
	require.Equal(t, `{"union":[{"@ref":"indexes/active_users"},{"@ref":"indexes/vip_users"}]}`, json)
}

func toJSON(expr Expr) (string, error) {
	bytes, err := writeJSON(expr)
	return string(bytes), err
}
