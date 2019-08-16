package faunadb

// LowerBound is a new optional parameter
func LowerBound(ref interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["lowerbound"] = wrap(ref)
	}
}

// UpperBound is a new optional parameter
func UpperBound(ref interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["upperbound"] = wrap(ref)
	}
}

// Range will provide the ability to limit a set based on lower and upper bounds of its natural order.
//
// Parameters:
//  set SetRef - A set reference
//
// Optional parameters:
//  Range(set, lowerBound, upperBound)
//
func Range(set interface{}, options ...OptionalParameter) Expr {
	return fn1("range", set, options...)
}

// Reduce function which may be used on arrays, pages or sets. This will behave similarly to foldLeft or reduce in functional languages.
// Reduce(set/array/page, init, fn)
func Reduce(coll, init, lambda interface{}) Expr {
	return fn3("reduce", lambda, "init", init, "collection", coll)
}
