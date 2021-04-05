package faunadb

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	assertMarshal(t, LongV(10), "10")
	assertMarshal(t, StringV("str"), `"str"`)
	assertMarshal(t, ArrayV{StringV("str"), LongV(10)}, `["str",10]`)
	assertMarshal(t, ObjectV{"x": LongV(10), "y": LongV(20)}, `{"x":10,"y":20}`)
	assertMarshal(t,
		ObjectV{"x": LongV(10), "y": LongV(20), "z": ArrayV{StringV("str"), ObjectV{"w": LongV(10)}}},
		`{"x":10,"y":20,"z":["str",{"w":10}]}`,
	)
	assertMarshal(t,
		TimeV(time.Date(2019, time.January, 1, 1, 30, 20, 5, time.UTC)),
		`{"@ts":"2019-01-01T01:30:20.000000005Z"}`,
	)
	assertMarshal(t,
		ObjectV{"ts": TimeV(time.Date(2019, time.January, 1, 1, 30, 20, 5, time.UTC))},
		`{"ts":{"@ts":"2019-01-01T01:30:20.000000005Z"}}`,
	)
}

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

func TestSerializeBytesV(t *testing.T) {
	assertJSON(t,
		BytesV{1, 2, 3, 4},
		`{"@bytes":"AQIDBA=="}`,
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

func TestSerializeStructWithIgnoredFields(t *testing.T) {
	type user struct {
		Name string `fauna:"name"`
		Age  int    `fauna:"-"`
	}

	assertJSON(t,
		Obj{"data": user{"Jhon", 42}},
		`{"object":{"data":{"object":{"name":"Jhon"}}}}`,
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
		RefCollection(Ref("collections/spells"), "42"),
		`{"id":"42","ref":{"@ref":"collections/spells"}}`,
	)
}

func TestSerializeCreate(t *testing.T) {
	assertJSON(t,
		Create(Ref("collections/spells"), Obj{
			"name": "fire",
		}),
		`{"create":{"@ref":"collections/spells"},"params":{"object":{"name":"fire"}}}`,
	)
}

func TestSerializeUpdate(t *testing.T) {
	assertJSON(t,
		Update(Ref("collections/spells/123"), Obj{
			"name": "fire",
		}),
		`{"params":{"object":{"name":"fire"}},"update":{"@ref":"collections/spells/123"}}`,
	)
}

func TestSerializeReplace(t *testing.T) {
	assertJSON(t,
		Replace(Ref("collections/spells/123"), Obj{
			"name": "fire",
		}),
		`{"params":{"object":{"name":"fire"}},"replace":{"@ref":"collections/spells/123"}}`,
	)
}

func TestSerializeDelete(t *testing.T) {
	assertJSON(t,
		Delete(Ref("collections/spells/123")),
		`{"delete":{"@ref":"collections/spells/123"}}`,
	)
}

func TestSerializeInsert(t *testing.T) {
	assertJSON(t,
		Insert(
			Ref("collections/spells/104979509696660483"),
			time.Unix(0, 0).UTC(),
			ActionCreate,
			Obj{"data": Obj{"name": "test"}},
		),
		`{"action":"create","insert":{"@ref":"collections/spells/104979509696660483"},`+
			`"params":{"object":{"data":{"object":{"name":"test"}}}},"ts":{"@ts":"1970-01-01T00:00:00Z"}}`,
	)
}

func TestSerializeRemove(t *testing.T) {
	assertJSON(t,
		Remove(
			Ref("collections/spells/104979509696660483"),
			time.Unix(0, 0).UTC(),
			ActionDelete,
		),
		`{"action":"delete","remove":{"@ref":"collections/spells/104979509696660483"},"ts":{"@ts":"1970-01-01T00:00:00Z"}}`,
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

func TestSerializeCreateCollection(t *testing.T) {
	assertJSON(t,
		CreateCollection(Obj{
			"name": "boons",
		}),
		`{"create_collection":{"object":{"name":"boons"}}}`,
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
			"source": Ref("collections/spells"),
		}),
		`{"create_index":{"object":{"name":"new-index","source":{"@ref":"collections/spells"}}}}`,
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

func TestSerializeCreateRole(t *testing.T) {
	assertJSON(t,
		CreateRole(Obj{
			"name": "a_role",
			"privileges": Arr{Obj{
				"resource": Ref("databases"),
				"actions":  Obj{"read": true},
			}},
		}),
		`{"create_role":{"object":{"name":"a_role","privileges":[`+
			`{"object":{"actions":{"object":{"read":true}},"resource":{"@ref":"databases"}}}]}}}`,
	)

	assertJSON(t,
		CreateRole(Obj{
			"name": "a_role",
			"privileges": Obj{
				"resource": Ref("databases"),
				"actions":  Obj{"read": true},
			},
		}),
		`{"create_role":{"object":{"name":"a_role","privileges":`+
			`{"object":{"actions":{"object":{"read":true}},"resource":{"@ref":"databases"}}}}}}`,
	)
}

func TestSerializeMoveDatabase(t *testing.T) {
	assertJSON(t,
		MoveDatabase(Database("source"), Database("dest")),
		`{"move_database":{"database":"source"},"to":{"database":"dest"}}`,
	)
}

func TestSerializeNull(t *testing.T) {
	assertJSON(t, Null(), `null`)
}

func TestSerializeNil(t *testing.T) {
	assertJSON(t, nil, `null`)
	assertJSON(t, Arr{nil, 0}, `[null,0]`)
	assertJSON(t, Arr{Obj{"hey": nil}, nil, NullV{}, Null()}, `[{"object":{"hey":null}},null,null,null]`)
	assertJSON(t, Obj{"we_have": nil}, `{"object":{"we_have":null}}`)
	assertJSON(t,
		Filter(nil, Obj{"we_have": nil, "they_have": Arr{nil, 0}}),
		`{"collection":null,"filter":{"object":{"they_have":[null,0],"we_have":null}}}`,
	)
	assertJSON(t, ToString(nil), `{"to_string":null}`)
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

func TestSerializeAt(t *testing.T) {
	assertJSON(t,
		At(1, Paginate(Match(Index("all_things")))),
		`{"at":1,"expr":{"paginate":{"match":{"index":"all_things"}}}}`,
	)

	assertJSON(t,
		At(Time("1970-01-01T00:00:00+00:00"), Paginate(Match(Index("all_things")))),
		`{"at":{"time":"1970-01-01T00:00:00+00:00"},"expr":{"paginate":{"match":{"index":"all_things"}}}}`,
	)

	assertJSON(t,
		At(TimeV(time.Unix(1, 2).UTC()), Paginate(Match(Index("all_things")))),
		`{"at":{"@ts":"1970-01-01T00:00:01.000000002Z"},"expr":{"paginate":{"match":{"index":"all_things"}}}}`,
	)
}

func TestSerializeLet(t *testing.T) {
	query := Let().Bind("v1", Ref("collections/spells/42")).Bind("v2", Index("spells")).Bind("a1", Index("all_things")).In(Exists(Var("v1")))
	assertJSON(t, query,
		`{"in":{"exists":{"var":"v1"}},"let":[{"v1":{"@ref":"collections/spells/42"}},{"v2":{"index":"spells"}},{"a1":{"index":"all_things"}}]}`,
	)
}

func TestSerializeIf(t *testing.T) {
	assertJSON(t,
		If(true, "exists", "does not exists"),
		`{"else":"does not exists","if":true,"then":"exists"}`,
	)
}

func TestSerializeAbort(t *testing.T) {
	assertJSON(t,
		Abort("abort message"),
		`{"abort":"abort message"}`,
	)
}

func TestSerializeDo(t *testing.T) {
	assertJSON(t,
		Do(Arr{
			Get(Ref("collections/spells/4")),
			Get(Ref("collections/spells/2")),
		}),
		`{"do":[[{"get":{"@ref":"collections/spells/4"}},{"get":{"@ref":"collections/spells/2"}}]]}`,
	)

	assertJSON(t,
		Do(
			Get(Ref("collections/spells/4")),
			Get(Ref("collections/spells/2")),
		),
		`{"do":[{"get":{"@ref":"collections/spells/4"}},{"get":{"@ref":"collections/spells/2"}}]}`,
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

func TestSerializeReverse(t *testing.T) {
	assertJSON(t,
		Reverse(Arr{1, 2, 3}),
		`{"reverse":[1,2,3]}`,
	)
}

func TestSerializeGet(t *testing.T) {
	assertJSON(t,
		Get(Ref("collections/spells/42")),
		`{"get":{"@ref":"collections/spells/42"}}`,
	)
}

func TestSerializeGetWithTimestamp(t *testing.T) {
	assertJSON(t,
		Get(
			Ref("collections/spells/42"),
			TS(time.Unix(0, 0).UTC()),
		),
		`{"get":{"@ref":"collections/spells/42"},"ts":{"@ts":"1970-01-01T00:00:00Z"}}`,
	)
}

func TestSerializeKeyFromSecret(t *testing.T) {
	assertJSON(t,
		KeyFromSecret("s3cr3t"),
		`{"key_from_secret":"s3cr3t"}`,
	)
}

func TestSerializeExists(t *testing.T) {
	assertJSON(t,
		Exists(Ref("collections/spells/42")),
		`{"exists":{"@ref":"collections/spells/42"}}`,
	)
}

func TestSerializeExistsWithTimestamp(t *testing.T) {
	assertJSON(t,
		Exists(
			Ref("collections/spells/42"),
			TS(time.Unix(1, 1).UTC()),
		),
		`{"exists":{"@ref":"collections/spells/42"},"ts":{"@ts":"1970-01-01T00:00:01.000000001Z"}}`,
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
			EventsOpt(true),
			Sources(true),
			TS(10),
			Size(2),
		),
		`{"after":{"@ref":"databases/test"},"before":{"@ref":"databases/test10"},"events":true,`+
			`"paginate":{"@ref":"databases"},"size":2,"sources":true,"ts":10}`,
	)
}

func TestSerializeFormat(t *testing.T) {
	assertJSON(t,
		Format("You have %d points left", 89),
		`{"format":"You have %d points left","values":89}`,
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

	assertJSON(t,
		Casefold("GET DOWN", Normalizer("NFK")),
		`{"casefold":"GET DOWN","normalizer":"NFK"}`,
	)
}

func TestSerializeStartsWith(t *testing.T) {
	assertJSON(t, StartsWith("faunadb", "fauna"), `{"search":"fauna","startswith":"faunadb"}`)
}

func TestSerializeEndsWith(t *testing.T) {
	assertJSON(t, EndsWith("faunadb", "db"), `{"endswith":"faunadb","search":"db"}`)
}

func TestSerializeContainsStr(t *testing.T) {
	assertJSON(t, ContainsStr("faunadb", "db"), `{"containsstr":"faunadb","search":"db"}`)
}

func TestSerializeContainsStrRegex(t *testing.T) {
	assertJSON(t, ContainsStrRegex("faunadb", "f(.*)db"), `{"containsstrregex":"faunadb","pattern":"f(.*)db"}`)
}

func TestSerializeRegexEscape(t *testing.T) {
	assertJSON(t, RegexEscape("f[a](.*)db"), `{"regexescape":"f[a](.*)db"}`)
}

func TestSerializeFindStr(t *testing.T) {
	assertJSON(t,
		FindStr("GET DOWN", "DOWN"),
		`{"find":"DOWN","findstr":"GET DOWN"}`,
	)
}

func TestSerializeFindStrRegex(t *testing.T) {
	assertJSON(t,
		FindStrRegex("GET DOWN", "DOWN"),
		`{"findstrregex":"GET DOWN","pattern":"DOWN"}`,
	)
}

func TestSerializeLength(t *testing.T) {
	assertJSON(t,
		Length("0123456789"),
		`{"length":"0123456789"}`,
	)
}

func TestSerializeLowerCase(t *testing.T) {
	assertJSON(t,
		LowerCase("0123456789"),
		`{"lowercase":"0123456789"}`,
	)
}

func TestSerializeLTrim(t *testing.T) {
	assertJSON(t,
		LTrim("0123456789"),
		`{"ltrim":"0123456789"}`,
	)
}

func TestSerializeRepeat(t *testing.T) {
	assertJSON(t,
		Repeat("0123456789", Number(2)),
		`{"number":2,"repeat":"0123456789"}`,
	)
}

func TestSerializeReplaceStr(t *testing.T) {
	assertJSON(t,
		ReplaceStr("0123456789", "12", "34"),
		`{"find":"12","replace":"34","replacestr":"0123456789"}`,
	)
}

func TestSerializeReplaceStrRegex(t *testing.T) {
	assertJSON(t,
		ReplaceStrRegex("0123456789", "12", "34"),
		`{"pattern":"12","replace":"34","replacestrregex":"0123456789"}`,
	)
}

func TestSerializeRTrim(t *testing.T) {
	assertJSON(t,
		RTrim("0123456789"),
		`{"rtrim":"0123456789"}`,
	)
}

func TestSerializeSpace(t *testing.T) {
	assertJSON(t,
		Space("0123456789"),
		`{"space":"0123456789"}`,
	)
}

func TestSerializeSubString(t *testing.T) {
	assertJSON(t,
		SubString("0123456789", 1),
		`{"start":1,"substring":"0123456789"}`,
	)
}

func TestSerializeTitleCase(t *testing.T) {
	assertJSON(t,
		TitleCase("0123456789"),
		`{"titlecase":"0123456789"}`,
	)
}

func TestSerializeTrim(t *testing.T) {
	assertJSON(t,
		Trim("0123456789"),
		`{"trim":"0123456789"}`,
	)
}

func TestSerializeUpperCase(t *testing.T) {
	assertJSON(t,
		UpperCase("0123456789"),
		`{"uppercase":"0123456789"}`,
	)
}

func TestSerializeTime(t *testing.T) {
	assertJSON(t,
		Time("1970-01-01T00:00:00+00:00"),
		`{"time":"1970-01-01T00:00:00+00:00"}`,
	)
}

func TestSerializeTimeAdd(t *testing.T) {
	assertJSON(t,
		TimeAdd("1970-01-01T00:00:00+00:00", 1, TimeUnitHour),
		`{"offset":1,"time_add":"1970-01-01T00:00:00+00:00","unit":"hour"}`,
	)
}

func TestSerializeTimeSubtract(t *testing.T) {
	assertJSON(t,
		TimeSubtract("1970-01-01T00:00:00+00:00", 1, TimeUnitDay),
		`{"offset":1,"time_subtract":"1970-01-01T00:00:00+00:00","unit":"day"}`,
	)
}

func TestSerializeTimeDiff(t *testing.T) {
	assertJSON(t,
		TimeDiff("1970-01-01T00:00:00+00:00", Epoch(1, TimeUnitSecond), TimeUnitSecond),
		`{"other":{"epoch":1,"unit":"second"},"time_diff":"1970-01-01T00:00:00+00:00","unit":"second"}`,
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

func TestSerializeNow(t *testing.T) {
	assertJSON(t, Now(), `{"now":null}`)
}

func TestSerializeDate(t *testing.T) {
	assertJSON(t,
		Date("1970-01-01"),
		`{"date":"1970-01-01"}`,
	)
}

func TestSerializeSingleton(t *testing.T) {
	assertJSON(t,
		Singleton(Collection("widgets")),
		`{"singleton":{"collection":"widgets"}}`,
	)
}

func TestSerializeEvents(t *testing.T) {
	assertJSON(t,
		Events(Collection("widgets")),
		`{"events":{"collection":"widgets"}}`,
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

func TestSerializeMerge(t *testing.T) {
	assertJSON(t, Merge(Obj{"x": 24}, Obj{"y": 25}), "{\"merge\":{\"object\":{\"x\":24}},\"with\":{\"object\":{\"y\":25}}}")
	assertJSON(t, Merge(Obj{"points": 900}, Obj{"name": "Trevor"}), "{\"merge\":{\"object\":{\"points\":900}},\"with\":{\"object\":{\"name\":\"Trevor\"}}}")
	assertJSON(t, Merge(Obj{}, Obj{"id": 9}), "{\"merge\":{\"object\":{}},\"with\":{\"object\":{\"id\":9}}}")
	assertJSON(t, Merge(Obj{"x": 24}, Obj{"y": 25}, ConflictResolver(Lambda(Arr{"key", "left", "right"}, Var("right")))), "{\"lambda\":{\"expr\":{\"var\":\"right\"},\"lambda\":[\"key\",\"left\",\"right\"]},\"merge\":{\"object\":{\"x\":24}},\"with\":{\"object\":{\"y\":25}}}")
}

func TestSerializeReduce(t *testing.T) {
	assertJSON(t,
		Reduce(Lambda(Arr{"accum", "value"}, Add(Var("accum"), Var("value"))), 0, []int{10, 20, 30}),
		`{"collection":[10,20,30],"initial":0,"reduce":{"expr":{"add":[{"var":"accum"},{"var":"value"}]},"lambda":["accum","value"]}}`,
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
			MatchTerm(Ref("indexes/spellbooks_by_owner"), Ref("collections/characters/104979509695139637")),
			Ref("indexes/spells_by_spellbook"),
		),
		`{"join":{"match":{"@ref":"indexes/spellbooks_by_owner"},"terms":{"@ref":"collections/characters/104979509695139637"}},`+
			`"with":{"@ref":"indexes/spells_by_spellbook"}}`,
	)
}

func TestSerializeRange(t *testing.T) {
	assertJSON(t,
		Range(
			MatchTerm(Ref("indexes/spellbooks_by_owner"), Ref("collections/characters/104979509695139637")),
			10,
			50,
		),
		`{"from":10,"range":{"match":{"@ref":"indexes/spellbooks_by_owner"},"terms":{"@ref":"collections/characters/104979509695139637"}},"to":50}`,
	)
}

func TestSerializeLogin(t *testing.T) {
	assertJSON(t,
		Login(
			Ref("collections/characters/104979509695139637"),
			Obj{"password": "abracadabra"},
		),
		`{"login":{"@ref":"collections/characters/104979509695139637"},"params":{"object":{"password":"abracadabra"}}}`,
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
		Identify(Ref("collections/characters/104979509695139637"), "abracadabra"),
		`{"identify":{"@ref":"collections/characters/104979509695139637"},"password":"abracadabra"}`,
	)
}

func TestSerializeIdentity(t *testing.T) {
	assertJSON(t,
		Identity(),
		`{"identity":null}`,
	)
}

func TestSerializeHasIdentity(t *testing.T) {
	assertJSON(t,
		HasIdentity(),
		`{"has_identity":null}`,
	)
}

func TestSerializeNewId(t *testing.T) {
	assertJSON(t,
		NewId(),
		`{"new_id":null}`,
	)
}

func TestSerializeDatabase(t *testing.T) {
	assertJSON(t,
		Database("test-db"),
		`{"database":"test-db"}`,
	)

	assertJSON(t,
		ScopedDatabase("test-db", Database("scope")),
		`{"database":"test-db","scope":{"database":"scope"}}`,
	)
}

func TestSerializeIndex(t *testing.T) {
	assertJSON(t,
		Index("test-index"),
		`{"index":"test-index"}`,
	)

	assertJSON(t,
		ScopedIndex("test-db", Database("scope")),
		`{"index":"test-db","scope":{"database":"scope"}}`,
	)
}

func TestSerializeClass(t *testing.T) {
	assertJSON(t,
		Class("test-class"),
		`{"class":"test-class"}`,
	)

	assertJSON(t,
		ScopedClass("test-db", Database("scope")),
		`{"class":"test-db","scope":{"database":"scope"}}`,
	)
}

func TestSerializeCollection(t *testing.T) {
	assertJSON(t,
		Collection("test-collection"),
		`{"collection":"test-collection"}`,
	)

	assertJSON(t,
		ScopedCollection("test-db", Database("scope")),
		`{"collection":"test-db","scope":{"database":"scope"}}`,
	)
}

func TestSerializeFunction(t *testing.T) {
	assertJSON(t,
		Function("test-function"),
		`{"function":"test-function"}`,
	)

	assertJSON(t,
		ScopedFunction("test-db", Database("scope")),
		`{"function":"test-db","scope":{"database":"scope"}}`,
	)
}

func TestSerializeRole(t *testing.T) {
	assertJSON(t,
		Role("test-role"),
		`{"role":"test-role"}`,
	)

	assertJSON(t,
		ScopedRole("test-role", Database("scope")),
		`{"role":"test-role","scope":{"database":"scope"}}`,
	)
}

func TestSerializeClasses(t *testing.T) {
	assertJSON(t,
		Classes(),
		`{"classes":null}`,
	)

	assertJSON(t,
		ScopedClasses(Database("scope")),
		`{"classes":{"database":"scope"}}`,
	)
}

func TestSerializeCollections(t *testing.T) {
	assertJSON(t,
		Collections(),
		`{"collections":null}`,
	)

	assertJSON(t,
		ScopedCollections(Database("scope")),
		`{"collections":{"database":"scope"}}`,
	)
}

func TestSerializeDocuments(t *testing.T) {

	assertJSON(t,
		Documents(Collection("users")),
		`{"documents":{"collection":"users"}}`,
	)
}

func TestSerializeIndexes(t *testing.T) {
	assertJSON(t,
		Indexes(),
		`{"indexes":null}`,
	)

	assertJSON(t,
		ScopedIndexes(Database("scope")),
		`{"indexes":{"database":"scope"}}`,
	)
}

func TestSerializeDatabases(t *testing.T) {
	assertJSON(t,
		Databases(),
		`{"databases":null}`,
	)

	assertJSON(t,
		ScopedDatabases(Database("scope")),
		`{"databases":{"database":"scope"}}`,
	)
}

func TestSerializeFunctions(t *testing.T) {
	assertJSON(t,
		Functions(),
		`{"functions":null}`,
	)

	assertJSON(t,
		ScopedFunctions(Database("scope")),
		`{"functions":{"database":"scope"}}`,
	)
}

func TestSerializeRoles(t *testing.T) {
	assertJSON(t,
		Roles(),
		`{"roles":null}`,
	)

	assertJSON(t,
		ScopedRoles(Database("scope")),
		`{"roles":{"database":"scope"}}`,
	)
}

func TestSerializeKeys(t *testing.T) {
	assertJSON(t,
		Keys(),
		`{"keys":null}`,
	)

	assertJSON(t,
		ScopedKeys(Database("scope")),
		`{"keys":{"database":"scope"}}`,
	)
}

func TestSerializeTokens(t *testing.T) {
	assertJSON(t,
		Tokens(),
		`{"tokens":null}`,
	)

	assertJSON(t,
		ScopedTokens(Database("scope")),
		`{"tokens":{"database":"scope"}}`,
	)
}

func TestSerializeCredentials(t *testing.T) {
	assertJSON(t,
		Credentials(),
		`{"credentials":null}`,
	)

	assertJSON(t,
		ScopedCredentials(Database("scope")),
		`{"credentials":{"database":"scope"}}`,
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

func TestSerializeContainsPath(t *testing.T) {
	assertJSON(t,
		ContainsPath(
			Arr{"favorites", "foods"},
			Obj{"favorites": Obj{
				"foods": Arr{"stake"},
			}},
		),
		`{"contains_path":["favorites","foods"],"in":{"object":{"favorites":{"object":{"foods":["stake"]}}}}}`,
	)
}

func TestSerializeContainsValue(t *testing.T) {
	assertJSON(t,
		ContainsValue(
			"steak",
			Obj{"favorites": Obj{
				"foods": Arr{"steak"},
			}},
		),
		`{"contains_value":"steak","in":{"object":{"favorites":{"object":{"foods":["steak"]}}}}}`,
	)
}
func TestSerializeContainsField(t *testing.T) {
	assertJSON(t,
		ContainsField(
			"favorites",
			Obj{"favorites": Obj{
				"foods": Arr{"steak"},
			}},
		),
		`{"contains_field":"favorites","in":{"object":{"favorites":{"object":{"foods":["steak"]}}}}}`,
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

func TestSelializeSelectAll(t *testing.T) {
	assertJSON(t,
		SelectAll(
			"foo",
			Arr{
				Obj{"foo": "bar"},
				Obj{"foo": "baz"},
			},
		),
		`{"from":[{"object":{"foo":"bar"}},{"object":{"foo":"baz"}}],"select_all":"foo"}`,
	)
}

func TestSerializeAbs(t *testing.T) {
	assertJSON(t,
		Abs(1),
		`{"abs":1}`,
	)

	assertJSON(t,
		Abs(-1),
		`{"abs":-1}`,
	)
}

func TestSerializeAcos(t *testing.T) {
	assertJSON(t,
		Acos(1),
		`{"acos":1}`,
	)

	assertJSON(t,
		Acos(1.23),
		`{"acos":1.23}`,
	)
}

func TestSerializeAsin(t *testing.T) {
	assertJSON(t,
		Asin(1),
		`{"asin":1}`,
	)

	assertJSON(t,
		Asin(1.23),
		`{"asin":1.23}`,
	)
}

func TestSerializeAtan(t *testing.T) {
	assertJSON(t,
		Atan(1),
		`{"atan":1}`,
	)

	assertJSON(t,
		Atan(1.23),
		`{"atan":1.23}`,
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

func TestSerializeBitAnd(t *testing.T) {
	assertJSON(t,
		BitAnd(Arr{1, 2}),
		`{"bitand":[1,2]}`,
	)
}

func TestSerializeBitNot(t *testing.T) {
	assertJSON(t,
		BitNot(1),
		`{"bitnot":1}`,
	)

	assertJSON(t,
		BitNot(3),
		`{"bitnot":3}`,
	)
}

func TestSerializeBitOr(t *testing.T) {
	assertJSON(t,
		BitOr(Arr{1, 2}),
		`{"bitor":[1,2]}`,
	)
}

func TestSerializeBitXor(t *testing.T) {
	assertJSON(t,
		BitXor(Arr{1, 2}),
		`{"bitxor":[1,2]}`,
	)
}

func TestSerializeCeil(t *testing.T) {
	assertJSON(t,
		Ceil(1.8),
		`{"ceil":1.8}`,
	)
}

func TestSerializeCos(t *testing.T) {
	assertJSON(t,
		Cos(1),
		`{"cos":1}`,
	)
}

func TestSerializeCosh(t *testing.T) {
	assertJSON(t,
		Cosh(1),
		`{"cosh":1}`,
	)
}

func TestSerializeDegrees(t *testing.T) {
	assertJSON(t,
		Degrees(1),
		`{"degrees":1}`,
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

func TestSerializeExp(t *testing.T) {
	assertJSON(t,
		Exp(1),
		`{"exp":1}`,
	)
}

func TestSerializeFloor(t *testing.T) {
	assertJSON(t,
		Floor(1),
		`{"floor":1}`,
	)
}

func TestSerializeHypot(t *testing.T) {
	assertJSON(t,
		Hypot(1, 2),
		`{"b":2,"hypot":1}`,
	)
}

func TestSerializeLn(t *testing.T) {
	assertJSON(t,
		Ln(1),
		`{"ln":1}`,
	)
}

func TestSerializeLog(t *testing.T) {
	assertJSON(t,
		Log(1),
		`{"log":1}`,
	)
}

func TestSerializeMax(t *testing.T) {
	assertJSON(t,
		Max(1),
		`{"max":1}`,
	)
}

func TestSerializeMin(t *testing.T) {
	assertJSON(t,
		Min(1),
		`{"min":1}`,
	)
}

func TestSerializePow(t *testing.T) {
	assertJSON(t,
		Pow(1, 2),
		`{"exp":2,"pow":1}`,
	)
}

func TestSerializeRadians(t *testing.T) {
	assertJSON(t,
		Radians(2),
		`{"radians":2}`,
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

func TestSerializeRound(t *testing.T) {
	assertJSON(t,
		Round(1.2345678),
		`{"round":1.2345678}`,
	)

	assertJSON(t,
		Round(3),
		`{"round":3}`,
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

func TestSerializeSign(t *testing.T) {
	assertJSON(t,
		Sign(1),
		`{"sign":1}`,
	)

	assertJSON(t,
		Sign(0),
		`{"sign":0}`,
	)
}

func TestSerializeSin(t *testing.T) {
	assertJSON(t,
		Sin(1),
		`{"sin":1}`,
	)

	assertJSON(t,
		Sin(0),
		`{"sin":0}`,
	)
}

func TestSerializeSinh(t *testing.T) {
	assertJSON(t,
		Sinh(1),
		`{"sinh":1}`,
	)

	assertJSON(t,
		Sinh(0),
		`{"sinh":0}`,
	)
}

func TestSerializeSqrt(t *testing.T) {
	assertJSON(t,
		Sqrt(1),
		`{"sqrt":1}`,
	)

	assertJSON(t,
		Sqrt(0),
		`{"sqrt":0}`,
	)
}

func TestSerializeTan(t *testing.T) {
	assertJSON(t,
		Tan(1),
		`{"tan":1}`,
	)

	assertJSON(t,
		Tan(0),
		`{"tan":0}`,
	)
}

func TestSerializeTanh(t *testing.T) {
	assertJSON(t,
		Tanh(1),
		`{"tanh":1}`,
	)

	assertJSON(t,
		Tanh(0),
		`{"tanh":0}`,
	)
}

func TestSerializeTrunc(t *testing.T) {
	assertJSON(t,
		Trunc(1.2345678),
		`{"trunc":1.2345678}`,
	)

	assertJSON(t,
		Trunc(3),
		`{"trunc":3}`,
	)
}

func TestSerializeAny(t *testing.T) {
	assertJSON(t, Any([]bool{true, true, true}), `{"any":[true,true,true]}`)
}

func TestSerializeAll(t *testing.T) {
	assertJSON(t, All([]bool{true, true, true}), `{"all":[true,true,true]}`)
}

func TestSerializeCount(t *testing.T) {
	expected := `{"count":[1,2,3,4,5]}`

	assertJSON(t,
		Count(Arr{1, 2, 3, 4, 5}),
		expected,
	)
}

func TestSerializeSum(t *testing.T) {
	expected := `{"sum":[1,2,3,4,5]}`

	assertJSON(t,
		Sum(Arr{1, 2, 3, 4, 5}),
		expected,
	)
}

func TestSerializeMean(t *testing.T) {
	expected := `{"mean":[1,2,3,4,5]}`

	assertJSON(t,
		Mean(Arr{1, 2, 3, 4, 5}),
		expected,
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

func TestSerializeToString(t *testing.T) {
	assertJSON(t,
		ToString(42),
		`{"to_string":42}`,
	)
}

func TestSerializeToNumber(t *testing.T) {
	assertJSON(t,
		ToNumber("42"),
		`{"to_number":"42"}`,
	)
}

func TestSerializeToDouble(t *testing.T) {
	assertJSON(t,
		ToDouble(42),
		`{"to_double":42}`,
	)
}

func TestSerializeToInteger(t *testing.T) {
	assertJSON(t,
		ToInteger(3.14159),
		`{"to_integer":3.14159}`,
	)
}

func TestSerializeToObject(t *testing.T) {
	assertJSON(t,
		ToObject(Arr{Arr{"x", 90}}),
		`{"to_object":[["x", 90]]}`,
	)
}

func TestSerializeToArray(t *testing.T) {
	assertJSON(t,
		ToArray(Obj{"x": 1}),
		`{"to_array":{"object":{"x":1}}}`,
	)
}

func TestSerializeToTime(t *testing.T) {
	assertJSON(t,
		ToTime("1970-01-01T00:00:00Z"),
		`{"to_time":"1970-01-01T00:00:00Z"}`,
	)
}

func TestSerializeToDate(t *testing.T) {
	assertJSON(t,
		ToDate("1970-01-01"),
		`{"to_date":"1970-01-01"}`,
	)
}

func assertJSON(t *testing.T, expr Expr, expected string) {
	bytes, err := json.Marshal(expr)

	require.NoError(t, err)
	require.JSONEq(t, expected, string(bytes))
}

func assertMarshal(t *testing.T, value Value, expected string) {
	bytes, err := MarshalJSON(value)
	require.NoError(t, err)

	var valueUnmarshal, expectedUnmarshal interface{}

	require.NoError(t, json.Unmarshal(bytes, &valueUnmarshal))
	require.NoError(t, json.Unmarshal([]byte(expected), &expectedUnmarshal))

	require.Equal(t, expectedUnmarshal, valueUnmarshal)
}

func TestSerializeCreateAccessProvider(t *testing.T) {
	assertJSON(t,
		CreateAccessProvider(Obj{
			"name":     "a_provider",
			"issuer":   "supported_issuer",
			"jwks_uri": "https://xxxx.auth0.com",
		}),
		`{"create_access_provider":{"object":{"name":"a_provider","issuer":`+
			`"supported_issuer","jwks_uri":"https://xxxx.auth0.com"}}}`,
	)

	assertJSON(t,
		CreateAccessProvider(Obj{
			"name":     "a_provider",
			"issuer":   "supported_issuer",
			"jwks_uri": "https://xxxx.auth0.com",
			"roles":    Arr{"roles"},
			"data":     Obj{"key": "value"},
		}),
		`{"create_access_provider":{"object":{"name":"a_provider","issuer":`+
			`"supported_issuer","jwks_uri":"https://xxxx.auth0.com",`+
			`"roles":["roles"],"data":{"object":{"key":"value"}}}}}`,
	)
}

func TestSerializeAccessProvider(t *testing.T) {
	assertJSON(t,
		AccessProvider("auth0"),
		`{"access_provider":"auth0"}`,
	)

	assertJSON(t,
		ScopedAccessProvider("auth0", Database("auth_db")),
		`{"access_provider":"auth0","scope":{"database":"auth_db"}}`,
	)
}

func TestSerializeAccessProviders(t *testing.T) {
	assertJSON(t,
		AccessProviders(),
		`{"access_providers":null}`,
	)

	assertJSON(t,
		ScopedAccessProviders(Database("auth_db")),
		`{"access_providers":{"database":"auth_db"}}`,
	)
}

func TestSerializeCurrentIdentity(t *testing.T) {
	assertJSON(t,
		CurrentIdentity(),
		`{"current_identity":null}`,
	)
}

func TestSerializeCurrentToken(t *testing.T) {
	assertJSON(t,
		CurrentToken(),
		`{"current_token":null}`,
	)
}

func TestSerializeHasCurrentIdentity(t *testing.T) {
	assertJSON(t,
		HasCurrentIdentity(),
		`{"has_current_identity":null}`,
	)
}

func TestSerializeHasCurrentToken(t *testing.T) {
	assertJSON(t,
		HasCurrentToken(),
		`{"has_current_token":null}`,
	)
}

func TestSerializeStructWithOmitEmptyTags(t *testing.T) {
	type TestStruct struct{}
	type OmitStruct struct {
		Name           string      `fauna:"name,omitempty"`
		Age            int         `fauna:"age,omitempty"`
		Payment        float64     `fauna:"payment,omitempty"`
		AgePointer     *int        `fauna:"agePointer,omitempty"`
		PaymentPointer *float64    `fauna:"paymentPointer,omitempty"`
		Struct         *TestStruct `fauna:"struct,omitempty"`
	}
	x := 10
	y := 42.42
	z := TestStruct{}
	tests := []struct {
		name string
		expr Expr
		want string
	}{
		{
			name: "Empty Int",
			expr: Obj{"data": OmitStruct{Name: "John", Age: 0}},
			want: `{"object":{"data":{"object":{"name":"John"}}}}`,
		},
		{
			name: "Empty String",
			expr: Obj{"data": OmitStruct{Name: "", Age: 30}},
			want: `{"object":{"data":{"object":{"age":30}}}}`,
		},
		{
			name: "Empty Float",
			expr: Obj{"data": OmitStruct{Name: "John", Payment: 0.0}},
			want: `{"object":{"data":{"object":{"name":"John"}}}}`,
		},
		{
			name: "Int pointer",
			expr: Obj{"data": OmitStruct{Name: "", AgePointer: &x}},
			want: `{"object":{"data":{"object":{"agePointer":10}}}}`,
		},
		{
			name: "Empty Int pointer",
			expr: Obj{"data": OmitStruct{Name: "John", AgePointer: nil}},
			want: `{"object":{"data":{"object":{"name":"John"}}}}`,
		},
		{
			name: "Float pointer",
			expr: Obj{"data": OmitStruct{Name: "", PaymentPointer: &y}},
			want: `{"object":{"data":{"object":{"paymentPointer":42.42}}}}`,
		},
		{
			name: "Empty Float pointer",
			expr: Obj{"data": OmitStruct{Name: "John", PaymentPointer: nil}},
			want: `{"object":{"data":{"object":{"name":"John"}}}}`,
		},
		{
			name: "All data is empty",
			expr: Obj{"data": OmitStruct{Name: "", Age: 0, Payment: 0.0, AgePointer: nil, PaymentPointer: nil}},
			want: `{"object":{"data":{"object":{}}}}`,
		},
		{
			name: "Non-empty struct",
			expr: Obj{"data": OmitStruct{Struct: &z}},
			want: `{"object":{"data":{"object":{"struct":{"object":{}}}}}}`,
		},
		{
			name: "Empty struct",
			expr: Obj{"data": OmitStruct{Struct: nil}},
			want: `{"object":{"data":{"object":{}}}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.expr)
			require.NoError(t, err)
			require.JSONEq(t, tt.want, string(got))
		})
	}
}
