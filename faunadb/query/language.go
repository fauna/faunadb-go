package query

import (
	"encoding/json"
	"faunadb/values"
)

type Expr interface{}

type fn map[string]Expr

func (f fn) MarshalJSON() ([]byte, error) {
	return json.Marshal(wrap(f))
}

type Obj map[string]Expr

func (obj Obj) MarshalJSON() ([]byte, error) {
	return json.Marshal(wrap(obj))
}

type Arr []Expr

func (arr Arr) MarshalJSON() ([]byte, error) {
	return json.Marshal(wrap(arr))
}

func Ref(id string) Expr {
	return values.RefV{id}
}

func Create(class, params Expr) Expr {
	return fn{"create": class, "params": params}
}

func Delete(ref Expr) Expr {
	return fn{"delete": ref}
}

func Get(ref Expr) Expr {
	return fn{"get": ref}
}
