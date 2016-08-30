package faunadb

func Ref(id string) Expr {
	return RefV{id}
}

func Get(ref interface{}) Expr {
	return fn{"get": ref}
}

func Create(ref, params interface{}) Expr {
	return fn{"create": ref, "params": params}
}

func Delete(ref interface{}) Expr {
	return fn{"delete": ref}
}
