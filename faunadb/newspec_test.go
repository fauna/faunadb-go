package faunadb

import (
	"testing"
)

// Test with e.g.:
// FAUNA_ROOT_KEY="dummy" go test -timeout 30s github.com/fauna/faunadb-go/faunadb -count=1 -run TestSerializeRange

// Range(set, lowerBound, upperBound)
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

// Reduce(set/array/page, init, fn)
func TestSerializeReduce(t *testing.T) {
	assertJSON(t,
		Reduce(Arr{1, 2, 3}, 0, Lambda("x", Var("x"))),
		`{"collection":[1,2,3],"init":0,"reduce":{"expr":{"var":"x"},"lambda":"x"}}`,
	)
}

func TestReducerAliases(t *testing.T) {
	assertJSON(t,
		Min(Arr{1, 2, 3}),
		`{"min":[1,2,3]}`,
	)
	assertJSON(t,
		Max(Arr{1, 2, 3}),
		`{"max":[1,2,3]}`,
	)
	assertJSON(t,
		Count(Arr{1, 2, 3}),
		`{"count":[1,2,3]}`,
	)
	assertJSON(t,
		Average(Arr{1, 2, 3}),
		`{"average":[1,2,3]}`,
	)
	assertJSON(t,
		Sum(Arr{1, 2, 3}),
		`{"sum":[1,2,3]}`,
	)
}
