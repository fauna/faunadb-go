package faunadb

import (
	"testing"
)

// Test with e.g.:
// FAUNA_ROOT_KEY="dummy" go test -timeout 30s github.com/fauna/faunadb-go/faunadb -count=1 -run TestSerializeRange

// Range(set, lowerBound, upperBound)
//
// Range() will provide the ability to limit a set based on lower and upper bounds of its natural order.
func TestSerializeRange(t *testing.T) {
	assertJSON(t,
		Range(Ref("databases")),
		`{"range":{"@ref":"databases"}}`,
	)

	assertJSON(t,
		Range(Match("users_by_name"), LowerBound("Brown"), UpperBound("Smith")),
		`{"lowerbound":"Brown","range":{"match":"users_by_name"},"upperbound":"Smith"}`,
	)

	assertJSON(t,
		Range(Match("users_by_last_first"), LowerBound(Arr{"Brown", "A"}), UpperBound("Smith")),
		`{"lowerbound":["Brown","A"],"range":{"match":"users_by_last_first"},"upperbound":"Smith"}`,
	)
}

// Filter(set/array/page, predicate)
//
// Filter() currently takes an array or page and filters its elements based on a predicate function.
// It will be enhanced to work on sets, in order to enable more ergonomic pagination and ability to compose it with other set modifiers.
func TestSerializeFilterSet(t *testing.T) {
	assertJSON(t,
		Filter(SetRefV{ObjectV{"name": StringV("a")}}, Lambda("x", Var("x"))),
		`{"collection":{"@set":{"name":"a"}},"filter":{"expr":{"var":"x"},"lambda":"x"}}`,
	)
}

// Map(set/array/page, fn)
//
// Map will be enhanced to work on sets in addition to pages and arrays.
// This will allow for more ergonomic pagination and combination with functions like Take() and Drop().
func TestSerializeMapSet(t *testing.T) {
	assertJSON(t,
		Map(SetRefV{ObjectV{"name": StringV("a")}}, Lambda("x", Var("x"))),
		`{"collection":{"@set":{"name":"a"}},"map":{"expr":{"var":"x"},"lambda":"x"}}`,
	)
}
