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

// Other range-like predicates to be added are RangeLT, RangeLTE, RangeGT, and RangeGTE,
// which let you bound the set only on one side, or can be combined to specify upper and lower bound exclusivity.

// RangeLT equals field_lt (less than)
func RangeLT(set, value interface{}) Expr {
	return fn2("range_lt", set, "value", value)
}

// RangeLTE equals field_lte (less than or equals)
func RangeLTE(set, value interface{}) Expr {
	return fn2("range_lte", set, "value", value)
}

// RangeGT equals field_gt (greater than)
func RangeGT(set, value interface{}) Expr {
	return fn2("range_gt", set, "value", value)
}

// RangeGTE equals field_gte (greater than or equals)
func RangeGTE(set, value interface{}) Expr {
	return fn2("range_gte", set, "value", value)
}

// Reduce function which may be used on arrays, pages or sets. This will behave similarly to foldLeft or reduce in functional languages.
// Reduce(set/array/page, init, fn)
func Reduce(coll, init, lambda interface{}) Expr {
	return fn3("reduce", lambda, "init", init, "collection", coll)
}

// Aliases for commonly used reducers, which may be used on arrays, pages or sets.
// Note: Min, Max already exist
func Count(coll interface{}) Expr   { return fn1("count", coll) }
func Average(coll interface{}) Expr { return fn1("average", coll) }
func Sum(coll interface{}) Expr     { return fn1("sum", coll) }

// Reverse(set/array/page)
// We will add a Reverse() function which can take an array, page, or set, and return the reversed version.
func Reverse(coll interface{}) Expr {
	return fn1("reduce", coll)
}

// Documents(set/array/page)
// We will add a Documents() built-in which will allow iterating through all of the documents in a collection. Combined with Filter(), Reduce(), Count(), etc.
// This will allow for arbitrary querying of a collection without the need for indexes. Initially this functionality will be backed by a scan of the collection.
func Documents(coll interface{}) Expr {
	return fn1("documents", coll)
}
