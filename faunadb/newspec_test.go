package faunadb

import (
	"testing"
)

// Test with e.g.:
// FAUNA_ROOT_KEY="dummy" go test -timeout 30s github.com/fauna/faunadb-go/faunadb -count=1 -run TestSerializeRange
// FAUNA_ROOT_KEY="dummy" go test github.com/fauna/faunadb-go/faunadb -run TestSerialize

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

func TestSerializeRangeComparisons(t *testing.T) {
	assertJSON(t,
		RangeLT(Match("coll_by_x"), 1),
		`{"range_lt":{"match":"coll_by_x"},"value":1}`,
	)
	assertJSON(t,
		RangeLTE(Match("coll_by_x"), 1),
		`{"range_lte":{"match":"coll_by_x"},"value":1}`,
	)
	assertJSON(t,
		RangeGT(Match("coll_by_x"), 1),
		`{"range_gt":{"match":"coll_by_x"},"value":1}`,
	)
	assertJSON(t,
		RangeGTE(Match("coll_by_x"), 1),
		`{"range_gte":{"match":"coll_by_x"},"value":1}`,
	)
}

// Filter(set/array/page, predicate)
// Filter() currently takes an array or page and filters its elements based on a predicate function.
// It will be enhanced to work on sets, in order to enable more ergonomic pagination and ability to compose it with other set modifiers.
func TestSerializeFilterSet(t *testing.T) {
	assertJSON(t,
		Filter(SetRefV{ObjectV{"name": StringV("a")}}, Lambda("x", Var("x"))),
		`{"collection":{"@set":{"name":"a"}},"filter":{"expr":{"var":"x"},"lambda":"x"}}`,
	)
}

// Map(set/array/page, fn)
// Map will be enhanced to work on sets in addition to pages and arrays.
// This will allow for more ergonomic pagination and combination with functions like Take() and Drop().
func TestSerializeMapSet(t *testing.T) {
	assertJSON(t,
		Map(SetRefV{ObjectV{"name": StringV("a")}}, Lambda("x", Var("x"))),
		`{"collection":{"@set":{"name":"a"}},"map":{"expr":{"var":"x"},"lambda":"x"}}`,
	)
}

// Drop(set/array/page, num)
// We will enhance the Drop() function to be able to take a set and return a set-like object which excludes the first N elements. This is equivalent to OFFSET in MySQL.
func TestSerializeDropSet(t *testing.T) {
	assertJSON(t,
		Drop(2, SetRefV{ObjectV{"name": StringV("a")}}),
		`{"collection":{"@set":{"name":"a"}},"drop":2}`,
	)
}

// Take(set/array/page, num)
// We will enhance the Take() function to be able to take a set and return an array of the first N elements.
// Combined with take() when used with drop(), can be used to simulate offset/limit style pagination.
func TestSerializeTakeSet(t *testing.T) {
	assertJSON(t,
		Take(2, SetRefV{ObjectV{"name": StringV("a")}}),
		`{"collection":{"@set":{"name":"a"}},"take":2}`,
	)
}

// Reduce(set/array/page, init, fn)
func TestSerializeReduce(t *testing.T) {
	assertJSON(t,
		Reduce(Arr{1, 2, 3}, 0, Lambda("x", Var("x"))),
		`{"collection":[1,2,3],"init":0,"reduce":{"expr":{"var":"x"},"lambda":"x"}}`,
	)

	assertJSON(t,
		Reduce(SetRefV{ObjectV{"name": StringV("a")}}, 0, Lambda("x", Var("x"))),
		`{"collection":{"@set":{"name":"a"}},"init":0,"reduce":{"expr":{"var":"x"},"lambda":"x"}}`,
	)
}

// Reduce(set/array/page, init, fn)
func TestSerializeReverse(t *testing.T) {
	assertJSON(t,
		Reverse(Arr{1, 2, 3}),
		`{"reduce":[1,2,3]}`,
	)

	assertJSON(t,
		Reverse(SetRefV{ObjectV{"name": StringV("a")}}),
		`{"reduce":{"@set":{"name":"a"}}}`,
	)
}

// Count(), Average(), Sum(), Min(), Max()
func TestSerializeReducerAliases(t *testing.T) {
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

func TestSerializeDocuments(t *testing.T) {
	assertJSON(t,
		Documents(Arr{1, 2, 3}),
		`{"documents":[1,2,3]}`,
	)

	assertJSON(t,
		Documents(SetRefV{ObjectV{"name": StringV("a")}}),
		`{"documents":{"@set":{"name":"a"}}}`,
	)
}
