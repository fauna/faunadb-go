package faunadb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractValueFromObject(t *testing.T) {
	var str string

	value := ObjectV{
		"data": ObjectV{
			"testField": StringV("A"),
		},
	}
	err := value.At(ObjKey("data", "testField")).Get(&str)

	require.NoError(t, err)
	require.Equal(t, "A", str)
}

func TestExtractValueFromArray(t *testing.T) {
	var num int

	value := ArrayV{
		ArrayV{
			LongV(1),
			LongV(2),
			LongV(3),
		},
	}
	err := value.At(ArrIndex(0, 2)).Get(&num)

	require.NoError(t, err)
	require.Equal(t, 3, num)
}

func TestFailToExtractFieldForNonTransversableValues(t *testing.T) {
	assertFailToExtractAnyFieldFor(t,
		StringV("test"),
		LongV(0),
		DoubleV(0),
		BooleanV(false),
		DateV(time.Now()),
		TimeV(time.Now()),
		RefV{"classes/spells"},
		SetRefV{map[string]Value{"any": StringV("set")}},
		NullV{},
	)
}

func TestReportKeyNotFound(t *testing.T) {
	assertFailToExtractField(t, ObjectV{"data": ObjectV{}}, ObjKey("data", "testField", "ref"),
		"Error while extracting path: data / testField / ref. Object key testField not found")
}

func TestReportIndexNotFound(t *testing.T) {
	assertFailToExtractField(t, ArrayV{}, ArrIndex(0).AtKey("ref"),
		"Error while extracting path: 0 / ref. Array index 0 not found")
}

func TestReportErrorPathWhenValueIsNotAnArray(t *testing.T) {
	value := ObjectV{
		"data": ObjectV{
			"testField": StringV("A"),
		},
	}

	assertFailToExtractField(t, value, ObjKey("data", "testField").AtIndex(1),
		"Error while extracting path: data / testField / 1. Expected value to be an array but was a faunadb.StringV")
}

func TestReportErrorPathWhenValueIsNotAnObject(t *testing.T) {
	assertFailToExtractField(t, ArrayV{ArrayV{}}, ArrIndex(0).AtKey("testField"),
		"Error while extracting path: 0 / testField. Expected value to be an object but was a faunadb.ArrayV")
}

func assertFailToExtractField(t *testing.T, value Value, field Field, message string) {
	_, err := value.At(field).GetValue()
	require.EqualError(t, err, message)

	var res Value
	require.EqualError(t, value.At(field).Get(&res), message)
}

func assertFailToExtractAnyFieldFor(t *testing.T, values ...Value) {
	key := ObjKey("anyField")
	index := ArrIndex(0)

	for _, value := range values {
		_, err := value.At(key).GetValue()
		assert.Contains(t, err.Error(), "Error while extracting path: anyField. Expected value to be an object but was a")

		_, err = value.At(index).GetValue()
		assert.Contains(t, err.Error(), "Error while extracting path: 0. Expected value to be an array but was a")
	}
}
