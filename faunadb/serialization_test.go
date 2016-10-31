package faunadb

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSerializeObjectV(t *testing.T) {
	assertJSON(t,
		ObjectV{
			"data": ObjectV{
				"name": StringV("test"),
			},
		},
		`{"object":{"data":{"object":{"name":"test"}}}}`,
	)
}

func TestSerializeArrayV(t *testing.T) {
	assertJSON(t,
		ArrayV{
			ObjectV{"name": StringV("a")},
			ObjectV{"name": StringV("b")},
		},
		`[{"object":{"name":"a"}},{"object":{"name":"b"}}]`,
	)
}

func TestSerializeSetRefV(t *testing.T) {
	assertJSON(t,
		SetRefV{
			ObjectV{"name": StringV("a")},
		},
		`{"@set":{"name":"a"}}`,
	)
}

func TestSerializeDateV(t *testing.T) {
	assertJSON(t,
		DateV(time.Unix(0, 0).UTC()),
		`{"@date":"1970-01-01"}`,
	)
}

func TestSerializeTimeV(t *testing.T) {
	assertJSON(t,
		TimeV(time.Unix(1, 2).UTC()),
		`{"@ts":"1970-01-01T00:00:01.000000002Z"}`,
	)
}

func TestSerializeUint(t *testing.T) {
	assertJSON(t,
		Obj{"x": uint(10)},
		`{"object":{"x":10}}`,
	)
}

func TestNotSerializeUintBiggerThanMaxInt(t *testing.T) {
	_, err := json.Marshal(Obj{"x": uint(math.MaxUint64)})
	require.Contains(t, err.Error(), "Error while encoding number to json: Uint value exceeds maximum int64")
}

func TestFailtToSerializeUnsupportedTypes(t *testing.T) {
	c := make(chan string)
	_, err := json.Marshal(Obj{"x": c})
	require.Contains(t, err.Error(), "Error while converting Expr to JSON: Non supported type chan")
}

func TestSerializeObject(t *testing.T) {
	assertJSON(t,
		Obj{"key": "value"},
		`{"object":{"key":"value"}}`,
	)
}

func TestSerializeNestedObjects(t *testing.T) {
	assertJSON(t,
		Obj{"key": Obj{"nested": "value"}},
		`{"object":{"key":{"object":{"nested":"value"}}}}`,
	)
}

func TestSerializeNestedMaps(t *testing.T) {
	assertJSON(t,
		Obj{"key": map[string]string{"nested": "value"}},
		`{"object":{"key":{"object":{"nested":"value"}}}}`,
	)
}

func TestSerializeInvalidMaps(t *testing.T) {
	_, err := json.Marshal(Obj{"key": map[int]string{1: "value"}})
	require.Contains(t, err.Error(), "Error while encoding map to json: All map keys must be of type string")
}

func TestSerializeArray(t *testing.T) {
	assertJSON(t,
		Arr{1, 2, 3},
		`[1,2,3]`,
	)
}

func TestSerializeWithNestedArrays(t *testing.T) {
	assertJSON(t,
		Arr{Arr{1, 2, 3}},
		`[[1,2,3]]`,
	)
}

func TestSerializeStruct(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}

	assertJSON(t,
		Obj{"data": user{"Jhon", 42}},
		`{"object":{"data":{"object":{"Age":42,"Name":"Jhon"}}}}`,
	)
}

func TestSerializeWithNonExportedFields(t *testing.T) {
	type user struct {
		Name string
		age  int
	}

	assertJSON(t,
		Obj{"data": user{"Jhon", 42}},
		`{"object":{"data":{"object":{"Name":"Jhon"}}}}`,
	)
}

func TestSerializeStructWithTags(t *testing.T) {
	type user struct {
		Name string `fauna:"name"`
		Age  int    `fauna:"age"`
	}

	assertJSON(t,
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

	assertJSON(t,
		Obj{"data": &user{"Jhon", &age}},
		`{"object":{"data":{"object":{"Age":42,"Name":"Jhon"}}}}`,
	)

	assertJSON(t,
		Obj{"data": &user{Name: "Jhon"}},
		`{"object":{"data":{"object":{"Age":null,"Name":"Jhon"}}}}`,
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

	assertJSON(t,
		Obj{"data": userInfo{user{"Jhon"}, map[string]string{"password": "1234"}}},
		`{"object":{"data":{"object":{"Credentials":{"object":{"password":"1234"}},"User":{"object":{"Name":"Jhon"}}}}}}`,
	)
}

func TestSerializeStructWithEmbeddedStructs(t *testing.T) {
	type Embedded struct {
		Str string
	}

	type Data struct {
		Int int
		Embedded
	}

	assertJSON(t,
		Obj{"data": Data{42, Embedded{"a string"}}},
		`{"object":{"data":{"object":{"Embedded":{"object":{"Str":"a string"}},"Int":42}}}}`,
	)
}

func TestSerializeRef(t *testing.T) {
	assertJSON(t,
		RefClass(Ref("classes/spells"), "42"),
		`{"id":"42","ref":{"@ref":"classes/spells"}}`,
	)
}

func TestSerializeCreate(t *testing.T) {
	assertJSON(t,
		Create(Ref("classes/spells"), Obj{
			"name": "fire",
		}),
		`{"create":{"@ref":"classes/spells"},"params":{"object":{"name":"fire"}}}`,
	)
}

func TestSerializeUpdate(t *testing.T) {
	assertJSON(t,
		Update(Ref("classes/spells/123"), Obj{
			"name": "fire",
		}),
		`{"params":{"object":{"name":"fire"}},"update":{"@ref":"classes/spells/123"}}`,
	)
}

func TestSerializeReplace(t *testing.T) {
	assertJSON(t,
		Replace(Ref("classes/spells/123"), Obj{
			"name": "fire",
		}),
		`{"params":{"object":{"name":"fire"}},"replace":{"@ref":"classes/spells/123"}}`,
	)
}

func TestSerializeDelete(t *testing.T) {
	assertJSON(t,
		Delete(Ref("classes/spells/123")),
		`{"delete":{"@ref":"classes/spells/123"}}`,
	)
}

func TestSerializeInsert(t *testing.T) {
	assertJSON(t,
		Insert(
			Ref("classes/spells/104979509696660483"),
			time.Unix(0, 0).UTC(),
			ActionCreate,
			Obj{"data": Obj{"name": "test"}},
		),
		`{"action":"create","insert":{"@ref":"classes/spells/104979509696660483"},`+
			`"params":{"object":{"data":{"object":{"name":"test"}}}},"ts":{"@ts":"1970-01-01T00:00:00Z"}}`,
	)
}

func TestSerializeRemove(t *testing.T) {
	assertJSON(t,
		Remove(
			Ref("classes/spells/104979509696660483"),
			time.Unix(0, 0).UTC(),
			ActionDelete,
		),
		`{"action":"delete","remove":{"@ref":"classes/spells/104979509696660483"},"ts":{"@ts":"1970-01-01T00:00:00Z"}}`,
	)
}

func TestSerializeCreateClass(t *testing.T) {
	assertJSON(t,
		CreateClass(Obj{
			"name": "boons",
		}),
		`{"create_class":{"object":{"name":"boons"}}}`,
	)
}

func TestSerializeCreateDatabase(t *testing.T) {
	assertJSON(t,
		CreateDatabase(Obj{
			"name": "db-next",
		}),
		`{"create_database":{"object":{"name":"db-next"}}}`,
	)
}

func TestSerializeCreateIndex(t *testing.T) {
	assertJSON(t,
		CreateIndex(Obj{
			"name":   "new-index",
			"source": Ref("classes/spells"),
		}),
		`{"create_index":{"object":{"name":"new-index","source":{"@ref":"classes/spells"}}}}`,
	)
}

func TestSerializeCreateKey(t *testing.T) {
	assertJSON(t,
		CreateKey(Obj{
			"database": Ref("databases/prydain"),
			"role":     "server",
		}),
		`{"create_key":{"object":{"database":{"@ref":"databases/prydain"},"role":"server"}}}`,
	)
}

func TestSerializeNull(t *testing.T) {
	assertJSON(t, Null(), `null`)
}

func TestSerializeNullOnObject(t *testing.T) {
	assertJSON(t, Obj{"data": Null()}, `{"object":{"data":null}}`)
}

func TestSerializeNullOnStruct(t *testing.T) {
	type structWithNull struct {
		Null *string
	}

	assertJSON(t,
		Obj{"data": structWithNull{nil}},
		`{"object":{"data":{"object":{"Null":null}}}}`,
	)
}

func TestSerializeLet(t *testing.T) {
	assertJSON(t,
		Let(
			Obj{"v1": Ref("classes/spells/42")},
			Exists(Var("v1")),
		),
		`{"in":{"exists":{"var":"v1"}},"let":{"v1":{"@ref":"classes/spells/42"}}}`,
	)
}

func TestSerializeIf(t *testing.T) {
	assertJSON(t,
		If(true, "exists", "does not exists"),
		`{"else":"does not exists","if":true,"then":"exists"}`,
	)
}

func TestSerializeDo(t *testing.T) {
	assertJSON(t,
		Do(Arr{
			Get(Ref("classes/spells/4")),
			Get(Ref("classes/spells/2")),
		}),
		`{"do":[{"get":{"@ref":"classes/spells/4"}},{"get":{"@ref":"classes/spells/2"}}]}`,
	)

	assertJSON(t,
		Do(
			Get(Ref("classes/spells/4")),
			Get(Ref("classes/spells/2")),
		),
		`{"do":[{"get":{"@ref":"classes/spells/4"}},{"get":{"@ref":"classes/spells/2"}}]}`,
	)
}

func TestSerializeLambda(t *testing.T) {
	assertJSON(t,
		Lambda("x", Var("x")),
		`{"expr":{"var":"x"},"lambda":"x"}`,
	)
}

func TestSerializeMap(t *testing.T) {
	assertJSON(t,
		Map(Arr{1, 2, 3}, Lambda("x", Var("x"))),
		`{"collection":[1,2,3],"map":{"expr":{"var":"x"},"lambda":"x"}}`,
	)
}

func TestSerializeForeach(t *testing.T) {
	assertJSON(t,
		Foreach(Arr{1, 2, 3}, Lambda("x", Var("x"))),
		`{"collection":[1,2,3],"foreach":{"expr":{"var":"x"},"lambda":"x"}}`,
	)
}

func TestSerializeFilter(t *testing.T) {
	assertJSON(t,
		Filter(Arr{true, false}, Lambda("x", Var("x"))),
		`{"collection":[true,false],"filter":{"expr":{"var":"x"},"lambda":"x"}}`,
	)
}

func TestSerializeTake(t *testing.T) {
	assertJSON(t,
		Take(2, Arr{1, 2, 3}),
		`{"collection":[1,2,3],"take":2}`,
	)
}

func TestSerializeDrop(t *testing.T) {
	assertJSON(t,
		Drop(2, Arr{1, 2, 3}),
		`{"collection":[1,2,3],"drop":2}`,
	)
}

func TestSerializePrepend(t *testing.T) {
	assertJSON(t,
		Prepend(Arr{1, 2, 3}, Arr{4, 5, 6}),
		`{"collection":[4,5,6],"prepend":[1,2,3]}`,
	)
}

func TestSerializeAppend(t *testing.T) {
	assertJSON(t,
		Append(Arr{4, 5, 6}, Arr{1, 2, 3}),
		`{"append":[4,5,6],"collection":[1,2,3]}`,
	)
}

func TestSerializeGet(t *testing.T) {
	assertJSON(t,
		Get(Ref("classes/spells/42")),
		`{"get":{"@ref":"classes/spells/42"}}`,
	)
}

func TestSerializeGetWithTimestamp(t *testing.T) {
	assertJSON(t,
		Get(
			Ref("classes/spells/42"),
			TS(time.Unix(0, 0).UTC()),
		),
		`{"get":{"@ref":"classes/spells/42"},"ts":{"@ts":"1970-01-01T00:00:00Z"}}`,
	)
}

func TestSerializeExists(t *testing.T) {
	assertJSON(t,
		Exists(Ref("classes/spells/42")),
		`{"exists":{"@ref":"classes/spells/42"}}`,
	)
}

func TestSerializeExistsWithTimestamp(t *testing.T) {
	assertJSON(t,
		Exists(
			Ref("classes/spells/42"),
			TS(time.Unix(1, 1).UTC()),
		),
		`{"exists":{"@ref":"classes/spells/42"},"ts":{"@ts":"1970-01-01T00:00:01.000000001Z"}}`,
	)
}

func TestSerializePaginate(t *testing.T) {
	assertJSON(t,
		Paginate(Ref("databases")),
		`{"paginate":{"@ref":"databases"}}`,
	)
}

func TestSerializePaginateWithParameters(t *testing.T) {
	assertJSON(t,
		Paginate(
			Ref("databases"),
			Before(Ref("databases/test10")),
			After(Ref("databases/test")),
			Events(true),
			Sources(true),
			TS(10),
			Size(2),
		),
		`{"after":{"@ref":"databases/test"},"before":{"@ref":"databases/test10"},"events":true,`+
			`"paginate":{"@ref":"databases"},"size":2,"sources":true,"ts":10}`,
	)
}

func TestSerializeConcat(t *testing.T) {
	assertJSON(t,
		Concat(Arr{"a", "b"}),
		`{"concat":["a","b"]}`,
	)
}

func TestSerializeConcatWithSeparator(t *testing.T) {
	assertJSON(t,
		Concat(Arr{"a", "b"}, Separator("/")),
		`{"concat":["a","b"],"separator":"/"}`,
	)
}

func TestSerializeCasefold(t *testing.T) {
	assertJSON(t,
		Casefold("GET DOWN"),
		`{"casefold":"GET DOWN"}`,
	)
}

func TestSerializeTime(t *testing.T) {
	assertJSON(t,
		Time("1970-01-01T00:00:00+00:00"),
		`{"time":"1970-01-01T00:00:00+00:00"}`,
	)
}

func TestSerializeEpoch(t *testing.T) {
	assertJSON(t,
		Arr{
			Epoch(0, TimeUnitSecond),
			Epoch(0, TimeUnitMillisecond),
			Epoch(0, TimeUnitMicrosecond),
			Epoch(0, TimeUnitNanosecond),
		},
		`[{"epoch":0,"unit":"second"},{"epoch":0,"unit":"millisecond"},`+
			`{"epoch":0,"unit":"microsecond"},{"epoch":0,"unit":"nanosecond"}]`,
	)
}

func TestSerializeDate(t *testing.T) {
	assertJSON(t,
		Date("1970-01-01"),
		`{"date":"1970-01-01"}`,
	)
}

func TestSerializeMatch(t *testing.T) {
	assertJSON(t,
		Match(Ref("databases")),
		`{"match":{"@ref":"databases"}}`,
	)
}

func TestSerializeMatchWithTerms(t *testing.T) {
	assertJSON(t,
		MatchTerm(
			Ref("indexes/spells_by_name"),
			"magic missile",
		),
		`{"match":{"@ref":"indexes/spells_by_name"},"terms":"magic missile"}`,
	)
}

func TestSerializeUnion(t *testing.T) {
	assertJSON(t,
		Union(Arr{
			Ref("indexes/active_users"),
			Ref("indexes/vip_users"),
		}),
		`{"union":[{"@ref":"indexes/active_users"},{"@ref":"indexes/vip_users"}]}`,
	)

	assertJSON(t,
		Union(
			Ref("indexes/active_users"),
			Ref("indexes/vip_users"),
		),
		`{"union":[{"@ref":"indexes/active_users"},{"@ref":"indexes/vip_users"}]}`,
	)
}

func TestSerializeIntersection(t *testing.T) {
	assertJSON(t,
		Intersection(Arr{
			Ref("indexes/active_users"),
			Ref("indexes/vip_users"),
		}),
		`{"intersection":[{"@ref":"indexes/active_users"},{"@ref":"indexes/vip_users"}]}`,
	)

	assertJSON(t,
		Intersection(
			Ref("indexes/active_users"),
			Ref("indexes/vip_users"),
		),
		`{"intersection":[{"@ref":"indexes/active_users"},{"@ref":"indexes/vip_users"}]}`,
	)
}

func TestSerializeDifference(t *testing.T) {
	assertJSON(t,
		Difference(Arr{
			Ref("indexes/active_users"),
			Ref("indexes/vip_users"),
		}),
		`{"difference":[{"@ref":"indexes/active_users"},{"@ref":"indexes/vip_users"}]}`,
	)

	assertJSON(t,
		Difference(
			Ref("indexes/active_users"),
			Ref("indexes/vip_users"),
		),
		`{"difference":[{"@ref":"indexes/active_users"},{"@ref":"indexes/vip_users"}]}`,
	)
}

func TestSerializeDistinct(t *testing.T) {
	assertJSON(t,
		Distinct(Ref("indexes/active_users")),
		`{"distinct":{"@ref":"indexes/active_users"}}`,
	)
}

func TestSerializeJoin(t *testing.T) {
	assertJSON(t,
		Join(
			MatchTerm(Ref("indexes/spellbooks_by_owner"), Ref("classes/characters/104979509695139637")),
			Ref("indexes/spells_by_spellbook"),
		),
		`{"join":{"match":{"@ref":"indexes/spellbooks_by_owner"},"terms":{"@ref":"classes/characters/104979509695139637"}},`+
			`"with":{"@ref":"indexes/spells_by_spellbook"}}`,
	)
}

func TestSerializeLogin(t *testing.T) {
	assertJSON(t,
		Login(
			Ref("classes/characters/104979509695139637"),
			Obj{"password": "abracadabra"},
		),
		`{"login":{"@ref":"classes/characters/104979509695139637"},"params":{"object":{"password":"abracadabra"}}}`,
	)
}

func TestSerializeLogout(t *testing.T) {
	assertJSON(t,
		Logout(true),
		`{"logout":true}`,
	)
}

func TestSerializeIndentify(t *testing.T) {
	assertJSON(t,
		Identify(Ref("classes/characters/104979509695139637"), "abracadabra"),
		`{"identify":{"@ref":"classes/characters/104979509695139637"},"password":"abracadabra"}`,
	)
}

func TestSerializeNextID(t *testing.T) {
	assertJSON(t,
		NextID(),
		`{"next_id":null}`,
	)
}

func TestSerializeDatabase(t *testing.T) {
	assertJSON(t,
		Database("test-db"),
		`{"database":"test-db"}`,
	)
}

func TestSerializeIndex(t *testing.T) {
	assertJSON(t,
		Index("test-index"),
		`{"index":"test-index"}`,
	)
}

func TestSerializeClass(t *testing.T) {
	assertJSON(t,
		Class("test-class"),
		`{"class":"test-class"}`,
	)
}

func TestSerializeEquals(t *testing.T) {
	assertJSON(t,
		Equals(Arr{"fire", "fire"}),
		`{"equals":["fire","fire"]}`,
	)

	assertJSON(t,
		Equals("fire", "air"),
		`{"equals":["fire","air"]}`,
	)
}

func TestSerializeContains(t *testing.T) {
	assertJSON(t,
		Contains(
			Arr{"favorites", "foods"},
			Obj{"favorites": Obj{
				"foods": Arr{"stake"},
			}},
		),
		`{"contains":["favorites","foods"],"in":{"object":{"favorites":{"object":{"foods":["stake"]}}}}}`,
	)
}

func TestSerializeSelect(t *testing.T) {
	assertJSON(t,
		Select(
			Arr{"favorites", "foods", 0},
			Obj{"favorites": Obj{
				"foods": Arr{"stake"},
			}},
		),
		`{"from":{"object":{"favorites":{"object":{"foods":["stake"]}}}},"select":["favorites","foods",0]}`,
	)

	assertJSON(t,
		Select(
			Arr{"favorites", "foods", 0},
			Obj{"favorites": Obj{
				"foods": Arr{"stake"},
			}},
			Default("no food"),
		),
		`{"default":"no food","from":{"object":{"favorites":{"object":{"foods":["stake"]}}}},"select":["favorites","foods",0]}`,
	)
}

func TestSerializeAdd(t *testing.T) {
	assertJSON(t,
		Add(Arr{1, 2}),
		`{"add":[1,2]}`,
	)

	assertJSON(t,
		Add(3, 4),
		`{"add":[3,4]}`,
	)
}

func TestSerializeMultiply(t *testing.T) {
	assertJSON(t,
		Multiply(Arr{1, 2}),
		`{"multiply":[1,2]}`,
	)

	assertJSON(t,
		Multiply(3, 4),
		`{"multiply":[3,4]}`,
	)
}

func TestSerializeSubtract(t *testing.T) {
	assertJSON(t,
		Subtract(Arr{1, 2}),
		`{"subtract":[1,2]}`,
	)

	assertJSON(t,
		Subtract(3, 4),
		`{"subtract":[3,4]}`,
	)
}

func TestSerializeDivide(t *testing.T) {
	assertJSON(t,
		Divide(Arr{1, 2}),
		`{"divide":[1,2]}`,
	)

	assertJSON(t,
		Divide(3, 4),
		`{"divide":[3,4]}`,
	)
}

func TestSerializeModulo(t *testing.T) {
	assertJSON(t,
		Modulo(Arr{1, 2}),
		`{"modulo":[1,2]}`,
	)

	assertJSON(t,
		Modulo(3, 4),
		`{"modulo":[3,4]}`,
	)
}

func TestSerializeLT(t *testing.T) {
	assertJSON(t,
		LT(Arr{1, 2}),
		`{"lt":[1,2]}`,
	)

	assertJSON(t,
		LT(3, 4),
		`{"lt":[3,4]}`,
	)
}

func TestSerializeLTE(t *testing.T) {
	assertJSON(t,
		LTE(Arr{1, 2}),
		`{"lte":[1,2]}`,
	)

	assertJSON(t,
		LTE(3, 4),
		`{"lte":[3,4]}`,
	)
}

func TestSerializeGT(t *testing.T) {
	assertJSON(t,
		GT(Arr{1, 2}),
		`{"gt":[1,2]}`,
	)

	assertJSON(t,
		GT(3, 4),
		`{"gt":[3,4]}`,
	)
}

func TestSerializeGTE(t *testing.T) {
	assertJSON(t,
		GTE(Arr{1, 2}),
		`{"gte":[1,2]}`,
	)

	assertJSON(t,
		GTE(3, 4),
		`{"gte":[3,4]}`,
	)
}

func TestSerializeAnd(t *testing.T) {
	assertJSON(t,
		And(Arr{true, false}),
		`{"and":[true,false]}`,
	)

	assertJSON(t,
		And(true, false),
		`{"and":[true,false]}`,
	)
}

func TestSerializeOr(t *testing.T) {
	assertJSON(t,
		Or(Arr{true, false}),
		`{"or":[true,false]}`,
	)

	assertJSON(t,
		Or(true, false),
		`{"or":[true,false]}`,
	)
}

func TestSerializeNot(t *testing.T) {
	assertJSON(t,
		Not(false),
		`{"not":false}`,
	)
}

func assertJSON(t *testing.T, expr Expr, expected string) {
	bytes, err := json.Marshal(expr)

	require.NoError(t, err)
	require.Equal(t, expected, string(bytes))
}
