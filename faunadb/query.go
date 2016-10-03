package faunadb

func Ref(id string) Expr       { return RefV{id} }
func Null() Expr               { return NullV{} }
func Get(ref interface{}) Expr { return fn{"get": ref} }

func Create(ref, params interface{}) Expr  { return fn{"create": ref, "params": params} }
func Update(ref, params interface{}) Expr  { return fn{"update": ref, "params": params} }
func Replace(ref, params interface{}) Expr { return fn{"replace": ref, "params": params} }
func Delete(ref interface{}) Expr          { return fn{"delete": ref} }

func Exists(ref interface{}) Expr { return fn{"exists": ref} }
