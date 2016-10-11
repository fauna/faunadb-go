package faunadb

import "encoding/json"

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

type Obj map[string]interface{}
type Arr []interface{}

func (obj Obj) expr() {}
func (arr Arr) expr() {}

func (obj Obj) MarshalJSON() ([]byte, error) {
	asMap := map[string]interface{}(obj)
	return json.Marshal(wrap(asMap))
}

func (arr Arr) MarshalJSON() ([]byte, error) {
	asArr := []interface{}(arr)
	return json.Marshal(wrap(asArr))
}

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
