package faunadb

// Equals checks if all args are equivalents.
//
// Parameters:
//  args []Value - A collection of expressions to check for equivalence.
//
// Returns:
//  bool - true if all elements are equals, false otherwise.
//
// See: https://app.fauna.com/documentation/reference/queryapi#miscellaneous-functions
func Equals(args ...interface{}) Expr { return equalsFn{Equals: wrap(varargs(args...))} }

type equalsFn struct {
	fnApply
	Equals Expr `json:"equals" faunarepr:"varargs"`
}

// Any evaluates to true if any element of the collection is true.
//
// Parameters:
// collection  - the collection
//
// Returns:
// Expr
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/any
func Any(collection interface{}) Expr {
	return anyFn{Any: wrap(collection)}
}

type anyFn struct {
	fnApply
	Any Expr `json:"any"`
}

// All evaluates to true if all elements of the collection are true.
//
// Parameters:
// collection - the collection
//
// Returns:
// Expr
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/all
func All(collection interface{}) Expr {
	return allFn{All: wrap(collection)}
}

type allFn struct {
	fnApply
	All Expr `json:"all"`
}

// LT returns true if each specified value is less than all the subsequent values. Otherwise LT returns false.
//
// Parameters:
//  args []number - A collection of terms to compare.
//
// Returns:
//  bool - true if all elements are less than each other from left to right.
//
// See: https://app.fauna.com/documentation/reference/queryapi#miscellaneous-functions
func LT(args ...interface{}) Expr { return ltFn{LT: wrap(varargs(args...))} }

type ltFn struct {
	fnApply
	LT Expr `json:"lt" faunarepr:"varargs"`
}

// LTE returns true if each specified value is less than or equal to all subsequent values. Otherwise LTE returns false.
//
// Parameters:
//  args []number - A collection of terms to compare.
//
// Returns:
//  bool - true if all elements are less than of equals each other from left to right.
//
// See: https://app.fauna.com/documentation/reference/queryapi#miscellaneous-functions
func LTE(args ...interface{}) Expr { return lteFn{LTE: wrap(varargs(args...))} }

type lteFn struct {
	fnApply
	LTE Expr `json:"lte" faunarepr:"varargs"`
}

// GT returns true if each specified value is greater than all subsequent values. Otherwise GT returns false.
// and false otherwise.
//
// Parameters:
//  args []number - A collection of terms to compare.
//
// Returns:
//  bool - true if all elements are greather than to each other from left to right.
//
// See: https://app.fauna.com/documentation/reference/queryapi#miscellaneous-functions
func GT(args ...interface{}) Expr { return gtFn{GT: wrap(varargs(args...))} }

type gtFn struct {
	fnApply
	GT Expr `json:"gt" faunarepr:"varargs"`
}

// GTE returns true if each specified value is greater than or equal to all subsequent values. Otherwise GTE returns false.
//
// Parameters:
//  args []number - A collection of terms to compare.
//
// Returns:
//  bool - true if all elements are greather than or equals to each other from left to right.
//
// See: https://app.fauna.com/documentation/reference/queryapi#miscellaneous-functions
func GTE(args ...interface{}) Expr { return gteFn{GTE: wrap(varargs(args...))} }

type gteFn struct {
	fnApply
	GTE Expr `json:"gte" faunarepr:"varargs"`
}

// And returns the conjunction of a list of boolean values.
//
// Parameters:
//  args []bool - A collection to compute the conjunction of.
//
// Returns:
//  bool - true if all elements are true, false otherwise.
//
// See: https://app.fauna.com/documentation/reference/queryapi#miscellaneous-functions
func And(args ...interface{}) Expr { return andFn{And: wrap(varargs(args...))} }

type andFn struct {
	fnApply
	And Expr `json:"and" faunarepr:"varargs"`
}

// Or returns the disjunction of a list of boolean values.
//
// Parameters:
//  args []bool - A collection to compute the disjunction of.
//
// Returns:
//  bool - true if at least one element is true, false otherwise.
//
// See: https://app.fauna.com/documentation/reference/queryapi#miscellaneous-functions
func Or(args ...interface{}) Expr { return orFn{Or: wrap(varargs(args...))} }

type orFn struct {
	fnApply
	Or Expr `json:"or" faunarepr:"varargs"`
}

// Not returns the negation of a boolean value.
//
// Parameters:
//  boolean bool - A boolean to produce the negation of.
//
// Returns:
//  bool - The value negated.
//
// See: https://app.fauna.com/documentation/reference/queryapi#miscellaneous-functions
func Not(boolean interface{}) Expr { return notFn{Not: wrap(boolean)} }

type notFn struct {
	fnApply
	Not Expr `json:"not"`
}
