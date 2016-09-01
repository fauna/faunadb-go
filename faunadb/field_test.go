package faunadb

import (
	"testing"

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

func TestReportKeyNotFound(t *testing.T) {
	assertFailToExtractField(t, ObjectV{"data": ObjectV{}}, ObjKey("data", "testField", "ref"),
		"Error while extrating path: data / testField / ref. Object key testField not found")
}

func TestReportIndexNotFound(t *testing.T) {
	assertFailToExtractField(t, ArrayV{}, ArrIndex(0).AtKey("ref"),
		"Error while extrating path: 0 / ref. Array index 0 not found")
}

func TestReportErrorPathWhenValueIsNotAnArray(t *testing.T) {
	value := ObjectV{
		"data": ObjectV{
			"testField": StringV("A"),
		},
	}

	assertFailToExtractField(t, value, ObjKey("data", "testField").AtIndex(1),
		"Error while extrating path: data / testField / 1. Expected value to be an array but was a faunadb.StringV")
}

func TestReportErrorPathWhenValueIsNotAnObject(t *testing.T) {
	assertFailToExtractField(t, ArrayV{ArrayV{}}, ArrIndex(0).AtKey("testField"),
		"Error while extrating path: 0 / testField. Expected value to be an object but was a faunadb.ArrayV")
}

func assertFailToExtractField(t *testing.T, value Value, field Field, message string) {
	_, err := value.At(field).GetValue()
	require.EqualError(t, err, message)

	var res Value
	require.EqualError(t, value.At(field).Get(&res), message)
}
