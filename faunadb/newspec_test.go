package faunadb

import (
	"testing"
)

// Test with e.g.:
// FAUNA_ROOT_KEY="dummy" go test -timeout 30s github.com/fauna/faunadb-go/faunadb -count=1 -run TestSerializeRange
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
