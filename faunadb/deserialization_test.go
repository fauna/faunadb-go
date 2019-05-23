package faunadb

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	var value Value

	require.NoError(t, UnmarshalJSON([]byte(`{"@ts":"2019-01-01T01:30:20.000000005Z"}`), &value))
	require.Equal(t, TimeV(time.Date(2019, time.January, 1, 1, 30, 20, 5, time.UTC)), value)

	require.NoError(t, UnmarshalJSON([]byte(`["str", 10]`), &value))
	require.Equal(t, ArrayV{StringV("str"), LongV(10)}, value)

	require.NoError(t, UnmarshalJSON([]byte(`{"x":10, "y": 20}`), &value))
	require.Equal(t, ObjectV{"x": LongV(10), "y": LongV(20)}, value)

	require.NoError(t, UnmarshalJSON([]byte(`{"x":10, "y": 20, "z": ["str", {"w": 30}]}`), &value))
	require.Equal(t, ObjectV{"x": LongV(10), "y": LongV(20), "z": ArrayV{StringV("str"), ObjectV{"w": LongV(30)}}}, value)
}

func TestDeserializeStringV(t *testing.T) {
	var str StringV

	require.NoError(t, decodeJSON(`"test"`, &str))
	require.Equal(t, StringV("test"), str)
}

func TestDeserializeString(t *testing.T) {
	var str string

	require.NoError(t, decodeJSON(`"test"`, &str))
	require.Equal(t, "test", str)
}

func TestDeserializeLongV(t *testing.T) {
	var num LongV

	require.NoError(t, decodeJSON("9223372036854775807", &num))
	require.Equal(t, LongV(math.MaxInt64), num)
}

func TestDeserializeLong(t *testing.T) {
	var num int64

	require.NoError(t, decodeJSON("9223372036854775807", &num))
	require.Equal(t, int64(math.MaxInt64), num)
}

func TestNotDeserializeUint(t *testing.T) {
	var num uint64

	require.EqualError(t,
		decodeJSON("18446744073709551615", &num),
		`strconv.ParseInt: parsing "18446744073709551615": value out of range`,
	)
}

func TestDeserializeDoubleV(t *testing.T) {
	var num DoubleV

	require.NoError(t, decodeJSON("10.64", &num))
	require.Equal(t, DoubleV(10.64), num)
}

func TestDeserializeDouble(t *testing.T) {
	var num float64

	require.NoError(t, decodeJSON("10.64", &num))
	require.Equal(t, 10.64, num)
}

func TestConvertNumbers(t *testing.T) {
	var num int
	var float float32

	require.NoError(t, decodeJSON("10", &num))
	require.Equal(t, 10, num)

	require.NoError(t, decodeJSON("10.32", &float))
	require.Equal(t, float32(10.32), float)
}

func TestDeserializeBooleanV(t *testing.T) {
	var boolean BooleanV

	require.NoError(t, decodeJSON("true", &boolean))
	require.Equal(t, BooleanV(true), boolean)
}

func TestDeserializeBooleanTrue(t *testing.T) {
	var boolean bool

	require.NoError(t, decodeJSON("true", &boolean))
	require.True(t, boolean)
}

func TestDeserializeBooleanFalse(t *testing.T) {
	var boolean bool

	require.NoError(t, decodeJSON("false", &boolean))
	require.False(t, boolean)
}

func TestDeserializeRefV(t *testing.T) {
	var ref RefV

	require.NoError(t, decodeJSON(`{"@ref":{"id":"42","class":{"@ref":{"id":"spells","class":{"@ref":{"id":"classes"}}}}}}`, &ref))
	require.Equal(t, RefV{"42", &RefV{"spells", NativeClasses(), nil}, nil}, ref)
}

func TestDeserializeDateV(t *testing.T) {
	var date DateV

	require.NoError(t, decodeJSON(`{ "@date": "1970-01-03" }`, &date))
	require.Equal(t, DateV(time.Date(1970, time.January, 3, 0, 0, 0, 0, time.UTC)), date)
}

func TestDeserializeDate(t *testing.T) {
	var date time.Time

	require.NoError(t, decodeJSON(`{ "@date": "1970-01-03" }`, &date))
	require.Equal(t, time.Date(1970, time.January, 3, 0, 0, 0, 0, time.UTC), date)
}

func TestDeserializeTimeV(t *testing.T) {
	var localTime TimeV

	require.NoError(t, decodeJSON(`{ "@ts": "1970-01-01T00:00:00.000000005Z" }`, &localTime))
	require.Equal(t, TimeV(time.Date(1970, time.January, 1, 0, 0, 0, 5, time.UTC)), localTime)
}

func TestDeserializeTime(t *testing.T) {
	var localTime time.Time

	require.NoError(t, decodeJSON(`{ "@ts": "1970-01-01T00:00:00.000000005Z" }`, &localTime))
	require.Equal(t, time.Date(1970, time.January, 1, 0, 0, 0, 5, time.UTC), localTime)
}

func TestDeserializeBytesV(t *testing.T) {
	var bytes BytesV

	require.NoError(t, decodeJSON(`{"@bytes": "AQIDBA=="}`, &bytes))
	require.Equal(t, BytesV{1, 2, 3, 4}, bytes)
}

func TestDeserializeBytes(t *testing.T) {
	var bytes []byte

	require.NoError(t, decodeJSON(`{"@bytes": "AQIDBA=="}`, &bytes))
	require.Equal(t, []byte{1, 2, 3, 4}, bytes)
}

func TestDeserializeQueryV(t *testing.T) {
	var query QueryV

	lambda := json.RawMessage(`{"lambda": "x", "expr": {"var": "x"}}`)

	require.NoError(t, decodeJSON(`{"@query": {"lambda": "x", "expr": {"var": "x"}}}`, &query))
	require.Equal(t, QueryV{lambda}, query)
}

func TestDeserializeQueryVInsideObjectV(t *testing.T) {
	var object ObjectV

	lambda := json.RawMessage(`{"lambda": "x", "expr": {"var": "x"}}`)

	require.NoError(t, decodeJSON(`{"a": "a", "b": {"lambda": {"@query": {"lambda": "x", "expr": {"var": "x"}}}}, "c": "c"}`, &object))
	require.Equal(t, ObjectV{"a": StringV("a"), "b": ObjectV{"lambda": QueryV{lambda}}, "c": StringV("c")}, object)
}

func TestDeserializeInvalidQueryV(t *testing.T) {
	var object ObjectV

	require.EqualError(t,
		decodeJSON(`{"query":{"@query": {}, "invalid":"what?"}`, &object),
		`Expected end of object but got "invalid"`)
}

func TestDeserializeSetRefV(t *testing.T) {
	var setRef SetRefV

	json := `
	{
		"@set": {
			"match": {"@ref":{"id":"spells_by_element","class":{"@ref":{"id":"indexes"}}}},
			"terms": "fire"
		}
	}
	`

	require.NoError(t, decodeJSON(json, &setRef))

	require.Equal(t,
		SetRefV{ObjectV{
			"match": RefV{"spells_by_element", NativeIndexes(), nil},
			"terms": StringV("fire"),
		}},
		setRef,
	)
}

func TestDecodeEmptyValue(t *testing.T) {
	var str string
	var value StringV

	require.NoError(t, value.Get(&str))
	require.Equal(t, "", str)
}

func TestDeserializeArrayV(t *testing.T) {
	var array ArrayV

	require.NoError(t, decodeJSON("[1]", &array))
	require.Equal(t, ArrayV{LongV(1)}, array)
}

func TestDeserializeArray(t *testing.T) {
	var array []int64

	require.NoError(t, decodeJSON("[1, 2, 3]", &array))
	require.Equal(t, []int64{1, 2, 3}, array)
}

func TestDeserializeEmptyArray(t *testing.T) {
	var array []int64

	require.NoError(t, decodeJSON("[]", &array))
	require.Empty(t, array)
}

func TestDeserializeArrayOnInvalidTarget(t *testing.T) {
	var wrongReference map[string]string

	require.EqualError(t,
		decodeJSON("[]", &wrongReference),
		"Error while decoding fauna value at: <root>. Can not decode array into a value of type \"map[string]string\"",
	)
}

func TestDeserializeObjectV(t *testing.T) {
	var object ObjectV

	require.NoError(t, decodeJSON(`{ "key": "value" }`, &object))
	require.Equal(t, ObjectV{"key": StringV("value")}, object)
}

func TestDeserializeObject(t *testing.T) {
	var object map[string]string

	require.NoError(t, decodeJSON(`{ "key": "value" }`, &object))
	require.Equal(t, map[string]string{"key": "value"}, object)
}

func TestDeserializeObjectLiteral(t *testing.T) {
	var object map[string]string

	require.NoError(t, decodeJSON(`{ "@obj": { "@name": "Test" } }`, &object))
	require.Equal(t, map[string]string{"@name": "Test"}, object)
}

func TestDeserializeEmptyObject(t *testing.T) {
	var object map[string]string

	require.NoError(t, decodeJSON(`{}`, &object))
	require.Empty(t, object)
}

func TestDeserializeObjectOnInvalidTarget(t *testing.T) {
	var wrongReference []string

	require.EqualError(t,
		decodeJSON("{}", &wrongReference),
		"Error while decoding fauna value at: <root>. Can not decode map into a value of type \"[]string\"",
	)
}

func TestDeserializeStruct(t *testing.T) {
	var object struct{ Name string }

	require.NoError(t, decodeJSON(`{ "Name": "Jhon" }`, &object))
	require.Equal(t, struct{ Name string }{"Jhon"}, object)
}

func TestDeserializeStructWithTags(t *testing.T) {
	type object struct {
		Name string `fauna:"name"`
		Age  int64  `fauna:"age"`
	}

	var obj object

	require.NoError(t, decodeJSON(`{ "name": "Jhon", "age": 10 }`, &obj))
	require.Equal(t, object{"Jhon", 10}, obj)
}

func TestDeserializeStructWithIgnoredFields(t *testing.T) {
	type object struct {
		Name string `fauna:"name"`
		Age  int64  `fauna:"-"`
	}

	var obj object

	require.NoError(t, decodeJSON(`{ "name": "Jhon", "age": 10 }`, &obj))
	require.Equal(t, object{"Jhon", 0}, obj)
}

func TestDeserializeStructWithPointers(t *testing.T) {
	type inner struct{ Name string }
	type object struct{ Inner *inner }

	var emptyObject, emptyInnerObject, obj *object

	require.NoError(t, decodeJSON(`{}`, &emptyObject))
	require.Equal(t, &object{}, emptyObject)

	require.NoError(t, decodeJSON(`{ "Inner": {} }`, &emptyInnerObject))
	require.Equal(t, &object{&inner{}}, emptyInnerObject)

	require.NoError(t, decodeJSON(`{ "Inner": { "Name": "Jhon"} }`, &obj))
	require.Equal(t, &object{&inner{"Jhon"}}, obj)
}

func TestDeserializeStructWithEmbeddedStructs(t *testing.T) {
	type Embedded struct {
		Str string
	}

	type Data struct {
		Int int
		Embedded
	}

	var data Data

	require.NoError(t, decodeJSON(`{"Int":42,"Embedded":{"Str":"a string"}}`, &data))
	require.Equal(t, Data{42, Embedded{"a string"}}, data)
}

func TestIgnoresUnmapedNamesInStruct(t *testing.T) {
	var object struct{ Name string }

	require.NoError(t, decodeJSON(`{ "Name": "Jhon", "SomeOtherThing": 42 }`, &object))
	require.Equal(t, struct{ Name string }{"Jhon"}, object)
}

func TestIgnoresPrivateMembersOfStruct(t *testing.T) {
	type object struct {
		name           string
		SomeOtherThing int64
	}

	var obj object

	require.NoError(t, decodeJSON(`{ "name": "Jhon", "SomeOtherThing": 42 }`, &obj))
	require.Equal(t, object{"", 42}, obj)
}

func TestReportErrorPath(t *testing.T) {
	var obj struct{ Arr []int }
	var aMap map[string]int

	require.EqualError(t,
		decodeJSON(`{ "Arr": [1, "right"] }`, &obj),
		"Error while decoding fauna value at: Arr / 1. Can not assign value of type \"faunadb.StringV\" to a value of type \"int\"",
	)

	require.EqualError(t,
		decodeJSON(`{ "One": 1, "Two": "2" }`, &aMap),
		"Error while decoding fauna value at: Two. Can not assign value of type \"faunadb.StringV\" to a value of type \"int\"",
	)
}

func TestDeserializeNullV(t *testing.T) {
	var null NullV

	require.NoError(t, decodeJSON(`null`, &null))
	require.Equal(t, NullV{}, null)
}

func TestDeserializeNull(t *testing.T) {
	var null string
	var pointer *string

	require.NoError(t, decodeJSON(`null`, &null))
	require.NoError(t, decodeJSON(`null`, &pointer))

	require.Equal(t, "", null)
	require.Nil(t, pointer)
}

func TestDeserializeComplexStruct(t *testing.T) {
	type nestedStruct struct {
		Nested string
	}

	type complexStruct struct {
		NonExistingField int
		nonPublicField   int
		TaggedString     string `fauna:"tagged"`
		Any              Value
		Ref              RefV
		Date             time.Time
		Time             time.Time
		LiteralObj       map[string]string
		Str              string
		Num              int
		Float            float64
		Boolean          bool
		IntArr           []int
		ObjArr           []nestedStruct
		Matrix           [][]int
		Map              map[string]string
		Object           nestedStruct
		Null             *nestedStruct
	}

	json := `
	{
		"Ref": {
			"@ref":{"id":"42","class":{"@ref":{"id":"spells","class":{"@ref":{"id":"classes"}}}}}
		},
		"Any": "any value",
		"Date": { "@date": "1970-01-03" },
		"Time":  { "@ts": "1970-01-01T00:00:00.000000005Z" },
		"LiteralObj":  { "@obj": {"@name": "@Jhon" } },
		"tagged": "TaggedString",
		"Str": "Jhon Knows",
		"Num": 31,
		"Float": 31.1,
		"Boolean": true,
		"IntArr": [1, 2, 3],
		"ObjArr": [{"Nested": "object1"}, {"Nested": "object2"}],
		"Matrix": [[1, 2], [3, 4]],
		"Map": {
			"key": "value"
		},
		"Object": {
			"Nested": "object"
		},
		"Null": null
	}
	`
	expected := complexStruct{
		TaggedString: "TaggedString",
		Ref:          RefV{"42", &RefV{"spells", NativeClasses(), nil}, nil},
		Any:          StringV("any value"),
		Date:         time.Date(1970, time.January, 3, 0, 0, 0, 0, time.UTC),
		Time:         time.Date(1970, time.January, 1, 0, 0, 0, 5, time.UTC),
		LiteralObj:   map[string]string{"@name": "@Jhon"},
		Str:          "Jhon Knows",
		Num:          31,
		Float:        31.1,
		Boolean:      true,
		IntArr: []int{
			1, 2, 3,
		},
		ObjArr: []nestedStruct{
			{"object1"},
			{"object2"},
		},
		Matrix: [][]int{
			{1, 2},
			{3, 4},
		},
		Map:    map[string]string{"key": "value"},
		Object: nestedStruct{"object"},
		Null:   nil,
	}

	var object complexStruct

	require.NoError(t, decodeJSON(json, &object))
	require.Equal(t, expected, object)
}

func decodeJSON(raw string, target interface{}) (err error) {
	buffer := []byte(raw)

	var value Value

	if value, err = parseJSON(bytes.NewReader(buffer)); err == nil {
		err = value.Get(&target)
	}

	return
}
