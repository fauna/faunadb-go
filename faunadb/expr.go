package faunadb

import (
	"encoding/json"
	"strconv"
	"strings"
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
	String() string
	expr() // Make sure only internal structures can be marked as valid expressions
}

type unescapedObj map[string]Expr
type unescapedArr []Expr
type invalidExpr struct{ err error }

func (obj unescapedObj) expr() {}
func (obj unescapedObj) String() string {
	if len(obj) == 1 && obj["object"] != nil {
		return obj["object"].String()
	}
	i := 0
	var sb strings.Builder
	sb.WriteString("Obj{")
	for k, v := range obj {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(strconv.Quote(k))
		sb.WriteString(": ")
		sb.WriteString(v.String())
		i++
	}
	sb.WriteString("}")
	return sb.String()
}

func (arr unescapedArr) expr() {}
func (arr unescapedArr) String() string {
	var sb strings.Builder
	sb.WriteString("Arr{")
	for i, v := range arr {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("}")
	return sb.String()
}

func (inv invalidExpr) expr() {}
func (inv invalidExpr) String() string {
	return "Invalid Expression: " + inv.err.Error()
}

func (inv invalidExpr) MarshalJSON() ([]byte, error) {
	return nil, inv.err
}

// Obj is a expression shortcut to represent any valid JSON object
type Obj map[string]interface{}

func (obj Obj) expr() {}
func (obj Obj) String() string {
	if len(obj) == 1 && obj["object"] != nil {
		return wrap(obj["object"]).String()
	}
	i := 0
	var sb strings.Builder
	sb.WriteString("Obj{")
	for k, v := range obj {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(strconv.Quote(k))
		sb.WriteString(": ")
		sb.WriteString(wrap(v).String())
		i++
	}
	sb.WriteString("}")
	return sb.String()
}

// Arr is a expression shortcut to represent any valid JSON array
type Arr []interface{}

func (arr Arr) expr() {}
func (arr Arr) String() string {
	var sb strings.Builder
	sb.WriteString("Arr{")
	for i, v := range arr {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(wrap(v).String())
	}
	sb.WriteString("}")
	return sb.String()
}

// MarshalJSON implements json.Marshaler for Obj expression
func (obj Obj) MarshalJSON() ([]byte, error) { return json.Marshal(wrap(obj)) }

// MarshalJSON implements json.Marshaler for Arr expression
func (arr Arr) MarshalJSON() ([]byte, error) { return json.Marshal(wrap(arr)) }
