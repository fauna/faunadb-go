package query

import (
	"encoding/json"
	"faunadb/values"
)

type Expr struct {
	wrapped interface{}
}

func (e Expr) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.wrapped)
}

type Obj map[string]interface{}

type Arr []interface{}

type fn map[string]interface{}

func Ref(id string) Expr {
	return wrap(values.RefV{id})
}

func Create(class, params interface{}) Expr {
	return wrapFn(fn{"create": class, "params": params})
}

func Delete(ref interface{}) Expr {
	return wrapFn(fn{"delete": ref})
}

func Get(ref interface{}) Expr {
	return wrapFn(fn{"get": ref})
}
