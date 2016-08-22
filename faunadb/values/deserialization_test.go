package values

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestDeserializeString(t *testing.T) {
	var str string
	assert(t,
		fromJson(`"test"`, &str),
		equal(str, "test"),
	)
}

func TestDeserializeNumbers(t *testing.T) {
	var num int64
	assert(t,
		fromJson("10", &num),
		equal(num, int64(10)),
	)
}

func TestDeserializeFloat(t *testing.T) {
	var num float64
	assert(t,
		fromJson("10.254", &num),
		equal(num, 10.254),
	)
}

func TestConvertNumbersIfPossible(t *testing.T) {
	var i int
	assert(t,
		fromJson("42", &i),
		equal(i, 42),
	)

	var f float32
	assert(t,
		fromJson("42.24", &f),
		equal(f, float32(42.24)),
	)
}

func TestDeserializeBooleanTrue(t *testing.T) {
	var boolean bool
	assert(t,
		fromJson("true", &boolean),
		equal(boolean, true),
	)
}

func TestDeserializeBooleanFalse(t *testing.T) {
	var boolean bool
	assert(t,
		fromJson("false", &boolean),
		equal(boolean, false),
	)
}

func TestDeserializeRefV(t *testing.T) {
	var ref RefV
	assert(t,
		fromJson(`{ "@ref": "classes/spells/42" }`, &ref),
		equal(ref, RefV{"classes/spells/42"}),
	)
}

func TestDeserializeDateV(t *testing.T) {
	var date DateV
	assert(t,
		fromJson(`{ "@date": "1970-01-03" }`, &date),
		equal(date, DateV{time.Date(1970, time.January, 3, 0, 0, 0, 0, time.UTC)}),
	)
}

func TestDeserializeTimeV(t *testing.T) {
	var timeV TimeV
	assert(t,
		fromJson(`{ "@ts": "1970-01-01T00:00:00.000000005Z" }`, &timeV),
		equal(timeV, TimeV{time.Date(1970, time.January, 1, 0, 0, 0, 5, time.UTC)}),
	)
}

func TestDeserializeSetRefV(t *testing.T) {
	json := `
	{
		"@set": {
			"match": { "@ref": "indexes/spells/spells_by_element" },
			"terms": "fire"
		}
	}
	`

	var setRef SetRefV

	assert(t,
		fromJson(json, &setRef),
		deepEqual(setRef, SetRefV{map[string]Value{
			"match": Value{RefV{"indexes/spells/spells_by_element"}},
			"terms": Value{"fire"},
		}}),
	)
}

func TestDeserializeValue(t *testing.T) {
	var value Value
	var ptr *Value

	assert(t,
		fromJson(`"a value string"`, &value),
		deepEqual(value, Value{"a value string"}),
	)

	assert(t,
		fromJson(`"a pointer string"`, &ptr),
		deepEqual(ptr, &Value{"a pointer string"}),
	)
}

func TestDecodeEmptyValue(t *testing.T) {
	var str string
	value := Value{}
	assert(t, value.Get(&str), equal(str, ""))
}

func TestDeserializeArray(t *testing.T) {
	var array []int64
	assert(t,
		fromJson("[1, 2, 3]", &array),
		deepEqual(array, []int64{1, 2, 3}),
	)
}

func TestDeserializeEmptyArray(t *testing.T) {
	var array []int64
	assert(t,
		fromJson("[]", &array),
		deepEqual(array, []int64{}),
	)
}

func TestDeserializeArrayOnInvalidTarget(t *testing.T) {
	var wrongReference map[string]string
	assert(t,
		fromJson("[]", &wrongReference),
		failWith("Error while decoding fauna value at: root. Expected internal value to be a json object but was []values.Value"),
	)
}

func TestDeserializeObject(t *testing.T) {
	var object map[string]string
	assert(t,
		fromJson(`{ "key": "value" }`, &object),
		deepEqual(object, map[string]string{"key": "value"}),
	)
}

func TestDeserializeObjectLiteral(t *testing.T) {
	var object map[string]string
	assert(t,
		fromJson(`{ "@obj": { "@name": "Test" } }`, &object),
		deepEqual(object, map[string]string{"@name": "Test"}),
	)
}

func TestDeserializeEmptyObject(t *testing.T) {
	var object map[string]string
	assert(t,
		fromJson("{}", &object),
		deepEqual(object, map[string]string{}),
	)
}

func TestDeserializeObjectOnInvalidTarget(t *testing.T) {
	var wrongReference []string
	assert(t,
		fromJson("{}", &wrongReference),
		failWith("Error while decoding fauna value at: root. Expected internal value to be a json array but was map[string]values.Value"),
	)
}

func TestDeserializeStruct(t *testing.T) {
	var object struct{ Name string }
	assert(t,
		fromJson(`{ "Name": "Jhon" }`, &object),
		deepEqual(object, struct{ Name string }{"Jhon"}),
	)
}

func TestDeserializeStructWithTags(t *testing.T) {
	type object struct {
		Name string `fauna:"name"`
		Age  int64  `fauna:"age"`
	}

	var obj object

	assert(t,
		fromJson(`{ "name": "Jhon", "age": 10 }`, &obj),
		deepEqual(obj, object{"Jhon", 10}),
	)
}

func TestDeserializeStructWithPointers(t *testing.T) {
	type inner struct{ Name string }
	type object struct{ Inner *inner }

	var emptyObject, emptyInnerObject, obj *object

	assert(t,
		fromJson(`{}`, &emptyObject),
		deepEqual(emptyObject, &object{}),
	)

	assert(t,
		fromJson(`{ "Inner": {} }`, &emptyInnerObject),
		deepEqual(emptyInnerObject, &object{&inner{}}),
	)

	assert(t,
		fromJson(`{ "Inner": { "Name": "Jhon"} }`, &obj),
		deepEqual(obj, &object{&inner{"Jhon"}}),
	)
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

	assert(t,
		fromJson(`{"Int":42,"Embedded":{"Str":"a string"}}`, &data),
		equal(data, Data{42, Embedded{"a string"}}),
	)
}

func TestIgnoresUnmapedNamesInStruct(t *testing.T) {
	var object struct{ Name string }
	assert(t,
		fromJson(`{ "Name": "Jhon", "SomeOtherThing": 42 }`, &object),
		deepEqual(object, struct{ Name string }{"Jhon"}),
	)
}

func TestIgnoresPrivateMembersOfStrct(t *testing.T) {
	var object, expected struct {
		name           string
		SomeOtherThing int64
	}

	json := `{ "name": "Jhon", "SomeOtherThing": 42 }`

	expected = struct {
		name           string
		SomeOtherThing int64
	}{SomeOtherThing: 42}

	assert(t,
		fromJson(json, &object),
		deepEqual(object, expected),
	)
}

func TestReportErrorPath(t *testing.T) {
	var obj struct{ Arr []int }
	assert(t,
		fromJson(`{ "Arr": [1, "right"] }`, &obj),
		failWith("Error while decoding fauna value at: Arr / 1. Can not assign value of type string to a value of type int"),
	)
}

func TestDeserializeComplexStruct(t *testing.T) {
	type nestedStruct struct {
		Nested string
	}

	type complexStruct struct {
		NonExistingField int
		nonPublicField   int
		TaggedString     string `fauna:"tagged"`
		Ref              RefV
		Any              Value
		Str              string
		Num              int64
		Float            float64
		Boolean          bool
		IntArr           []int64
		ObjArr           []nestedStruct
		Matrix           [][]int64
		Map              map[string]string
		Object           nestedStruct
	}

	var object, expected complexStruct

	json := `
	{
		"Ref": {
			"@ref": "classes/spells/42"
		},
		"Any": "any value",
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
		}
	}
	`

	expected = complexStruct{
		TaggedString: "TaggedString",
		Ref:          RefV{"classes/spells/42"},
		Any:          Value{"any value"},
		Str:          "Jhon Knows",
		Num:          31,
		Float:        31.1,
		Boolean:      true,
		IntArr: []int64{
			1, 2, 3,
		},
		ObjArr: []nestedStruct{
			nestedStruct{"object1"},
			nestedStruct{"object2"},
		},
		Matrix: [][]int64{
			{1, 2},
			{3, 4},
		},
		Map:    map[string]string{"key": "value"},
		Object: nestedStruct{"object"},
	}

	assert(t,
		fromJson(json, &object),
		deepEqual(object, expected),
	)
}

func fromJson(raw string, pointer interface{}) (err error) {
	bytes := []byte(raw)

	var value *Value

	if err = json.Unmarshal(bytes, &value); err == nil {
		err = value.Get(&pointer)
	}

	return
}

type assertion func(error) error

func assert(t *testing.T, previousError error, check assertion) (err error) {
	err = check(previousError)

	if err != nil {
		t.Error(err)
	}

	return
}

func equal(actual, expected interface{}) assertion {
	return checkThat(actual == expected, actual, expected)
}

func deepEqual(actual, expected interface{}) assertion {
	return checkThat(reflect.DeepEqual(actual, expected), actual, expected)
}

func failWith(message string) assertion {
	return func(err error) error {
		if err == nil {
			return fmt.Errorf("Should have faild with message: %s", message)
		}

		errorMessage := err.Error()
		return checkThat(errorMessage == message, errorMessage, message)(nil)
	}
}

func checkThat(passed bool, actual, expected interface{}) assertion {
	return func(err error) error {
		if err != nil {
			return err
		}

		if !passed {
			return fmt.Errorf("\n%10s: %#v\n%10s: %#v", "Expected", expected, "got", actual)
		}

		return nil
	}
}
