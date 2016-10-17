package faunadb

import "encoding/json"

/*
Expr represents query language expressions.

Expressions are created by using one of the functions at query.go. Query functions are designed to compose with other
query function as well as with custom data structures. For example:

	type User struct {
		Name string
	}

	_, _ := client.Query(
		Create(
			Ref("classes/users"),
			Obj{"data": User{"Jhon"}},
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
func (inv invalidExpr) expr()  {}

func (inv invalidExpr) MarshalJSON() ([]byte, error) {
	return nil, inv.err
}

// Obj is a expression shortcut to represent any valid JSON object
type Obj map[string]interface{}

// Arr is a expression shortcut to represent any valid JSON array
type Arr []interface{}

func (obj Obj) expr() {}
func (arr Arr) expr() {}

// MarshalJSON implements json.Marshaler for Obj expression
func (obj Obj) MarshalJSON() ([]byte, error) { return json.Marshal(wrap(obj)) }

// MarshalJSON implements json.Marshaler for Arr expression
func (arr Arr) MarshalJSON() ([]byte, error) { return json.Marshal(wrap(arr)) }

// OptionalParameter describes optional parameters for query language functions
type OptionalParameter func(unescapedObj)

func applyOptionals(options []OptionalParameter, fn unescapedObj) Expr {
	for _, option := range options {
		option(fn)
	}
	return fn
}

func fn1(k1 string, v1 interface{}, options ...OptionalParameter) Expr {
	return applyOptionals(options, unescapedObj{
		k1: wrap(v1),
	})
}

func fn2(k1 string, v1 interface{}, k2 string, v2 interface{}, options ...OptionalParameter) Expr {
	return applyOptionals(options, unescapedObj{
		k1: wrap(v1),
		k2: wrap(v2),
	})
}

func fn3(k1 string, v1 interface{}, k2 string, v2 interface{}, k3 string, v3 interface{}, options ...OptionalParameter) Expr {
	return applyOptionals(options, unescapedObj{
		k1: wrap(v1),
		k2: wrap(v2),
		k3: wrap(v3),
	})
}

func fn4(k1 string, v1 interface{}, k2 string, v2 interface{}, k3 string, v3 interface{}, k4 string, v4 interface{}, options ...OptionalParameter) Expr {
	return applyOptionals(options, unescapedObj{
		k1: wrap(v1),
		k2: wrap(v2),
		k3: wrap(v3),
		k4: wrap(v4),
	})
}
