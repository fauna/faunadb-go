package faunadb

import (
	"encoding/json"
)

/*
Expr is the base type for FaunaDB query language expressions.

Expressions are created by using the query language functions in query.go. Query functions are designed to compose with each other, as well as with
custom data structures. For example:

	type User struct {
		Name string
	}

	_, _ := client.Query(
		Create(
			Collection("users"),
			Obj{"data": User{"John"}},
		),
	)

*/
type Expr interface {
	expr() // Make sure only internal structures can be marked as valid expressions
}

type unescapedObj map[string]Expr
type unescapedArr []Expr
type invalidExpr struct{ err error }

func (obj unescapedObj) expr() {}

func (arr unescapedArr) expr() {}

func (inv invalidExpr) expr() {}

func (inv invalidExpr) MarshalJSON() ([]byte, error) {
	return nil, inv.err
}

// Obj is a expression shortcut to represent any valid JSON object
type Obj map[string]interface{}

func (obj Obj) expr() {}

// Arr is a expression shortcut to represent any valid JSON array
type Arr []interface{}

func (arr Arr) expr() {}

// MarshalJSON implements json.Marshaler for Obj expression
func (obj Obj) MarshalJSON() ([]byte, error) { return json.Marshal(wrap(obj)) }

// MarshalJSON implements json.Marshaler for Arr expression
func (arr Arr) MarshalJSON() ([]byte, error) { return json.Marshal(wrap(arr)) }
