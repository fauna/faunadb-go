package faunadb

// Helper
func unescapedBindings(obj Obj) unescapedObj {
	res := make(unescapedObj, len(obj))

	for k, v := range obj {
		res[k] = wrap(v)
	}

	return res
}

// LetRef binds values to one or more variables as go pointer.
//
// Parameters:
//  bindings Object - An object binding a variable name to a value.
//  in Expr - An expression to be evaluated.
//
// Returns:
//  Value - The result of the given expression.
//
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func LetPtr(bindings Obj, in *Obj) Expr {

	return unescapedObj{
		"let": wrap(unescapedBindings(bindings)),
		"in":  in,
	}
}

// IfPtr is like if but preserves a pointer
func IfPtr(cond, then interface{}, elze *Obj) Expr {
	return unescapedObj{
		"if":   wrap(cond),
		"then": wrap(then),
		"else": elze,
	}
}
