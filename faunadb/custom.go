package faunadb

// Helper
func unescapedBindings(obj Obj) unescapedObj {
	res := make(unescapedObj, len(obj))

	for k, v := range obj {
		res[k] = wrap(v)
	}

	return res
}

// InPtr is In with pointer bindings support
func (lb *LetBuilder) InPtr(in Expr) Expr {
	// return fn2("let", lb.bindings, "in", in)
	// return LetPtr(lb.bindings, &in)
	return unescapedObj{
		"let": wrap(lb.bindings),
		"in":  in,
	}
}

// LetPtr is the original Let implementation with pointer support in bindings
// Note: Has been updated to bindings in an array similar to builtin LetBuilder
// Consider using `Let().Bind().InPtr(&in)` instead
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
		"let": wrap(Arr{unescapedBindings(bindings)}),
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
