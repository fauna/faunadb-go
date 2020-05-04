package faunadb

// Collections

// Map applies the lambda expression on each element of a collection or Page.
// It returns the result of each application on a collection of the same type.
//
// Parameters:
//  coll []Value - The collection of elements to iterate.
//  lambda Lambda - A lambda function to be applied to each element of the collection. See Lambda() function.
//
// Returns:
//  []Value - A new collection with elements transformed by the lambda function.
//
// See: https://app.fauna.com/documentation/reference/queryapi#collections
func Map(coll, lambda interface{}) Expr { return mapFn{Map: wrap(lambda), Collection: wrap(coll)} }

type mapFn struct {
	fnApply
	Map        Expr `json:"map"`
	Collection Expr `json:"collection"`
}

func (fn mapFn) String() string { return printFn(fn) }

// Foreach applies the lambda expression on each element of a collection or Page.
// The original collection is returned.
//
// Parameters:
//  coll []Value - The collection of elements to iterate.
//  lambda Lambda - A lambda function to be applied to each element of the collection. See Lambda() function.
//
// Returns:
//  []Value - The original collection.
//
// See: https://app.fauna.com/documentation/reference/queryapi#collections
func Foreach(coll, lambda interface{}) Expr {
	return foreachFn{Foreach: wrap(lambda), Collection: wrap(coll)}
}

type foreachFn struct {
	fnApply
	Foreach    Expr `json:"foreach"`
	Collection Expr `json:"collection"`
}

func (fn foreachFn) String() string { return printFn(fn) }

// Filter applies the lambda expression on each element of a collection or Page.
// It returns a new collection of the same type containing only the elements in which the
// function application returned true.
//
// Parameters:
//  coll []Value - The collection of elements to iterate.
//  lambda Lambda - A lambda function to be applied to each element of the collection. The lambda function must return a boolean value. See Lambda() function.
//
// Returns:
//  []Value - A new collection.
//
// See: https://app.fauna.com/documentation/reference/queryapi#collections
func Filter(coll, lambda interface{}) Expr {
	return filterFn{Filter: wrap(lambda), Collection: wrap(coll)}
}

type filterFn struct {
	fnApply
	Filter     Expr `json:"filter"`
	Collection Expr `json:"collection"`
}

func (fn filterFn) String() string { return printFn(fn) }

// Take returns a new collection containing num elements from the head of the original collection.
//
// Parameters:
//  num int64 - The number of elements to take from the collection.
//  coll []Value - The collection of elements.
//
// Returns:
//  []Value - A new collection.
//
// See: https://app.fauna.com/documentation/reference/queryapi#collections
func Take(num, coll interface{}) Expr { return takeFn{Take: wrap(num), Collection: wrap(coll)} }

type takeFn struct {
	fnApply
	Take       Expr `json:"take"`
	Collection Expr `json:"collection"`
}

func (fn takeFn) String() string { return printFn(fn) }

// Drop returns a new collection containing the remaining elements from the original collection
// after num elements have been removed.
//
// Parameters:
//  num int64 - The number of elements to drop from the collection.
//  coll []Value - The collection of elements.
//
// Returns:
//  []Value - A new collection.
//
// See: https://app.fauna.com/documentation/reference/queryapi#collections
func Drop(num, coll interface{}) Expr { return dropFn{Drop: wrap(num), Collection: wrap(coll)} }

type dropFn struct {
	fnApply
	Drop       Expr `json:"drop"`
	Collection Expr `json:"collection"`
}

func (fn dropFn) String() string { return printFn(fn) }

// Prepend returns a new collection that is the result of prepending elems to coll.
//
// Parameters:
//  elems []Value - Elements to add to the beginning of the other collection.
//  coll []Value - The collection of elements.
//
// Returns:
//  []Value - A new collection.
//
// See: https://app.fauna.com/documentation/reference/queryapi#collections
func Prepend(elems, coll interface{}) Expr {
	return prependFn{Prepend: wrap(elems), Collection: wrap(coll)}
}

type prependFn struct {
	fnApply
	Prepend    Expr `json:"prepend"`
	Collection Expr `json:"collection"`
}

func (fn prependFn) String() string { return printFn(fn) }

// Append returns a new collection that is the result of appending elems to coll.
//
// Parameters:
//  elems []Value - Elements to add to the end of the other collection.
//  coll []Value - The collection of elements.
//
// Returns:
//  []Value - A new collection.
//
// See: https://app.fauna.com/documentation/reference/queryapi#collections
func Append(elems, coll interface{}) Expr {
	return appendFn{Append: wrap(elems), Collection: wrap(coll)}
}

type appendFn struct {
	fnApply
	Append     Expr `json:"append"`
	Collection Expr `json:"collection"`
}

func (fn appendFn) String() string { return printFn(fn) }

// IsEmpty returns true if the collection is the empty set, else false.
//
// Parameters:
//  coll []Value - The collection of elements.
//
// Returns:
//   bool - True if the collection is empty, else false.
//
// See: https://app.fauna.com/documentation/reference/queryapi#collections
func IsEmpty(coll interface{}) Expr { return isEmptyFn{IsEmpty: wrap(coll)} }

type isEmptyFn struct {
	fnApply
	IsEmpty Expr `json:"is_empty"`
}

func (fn isEmptyFn) String() string { return printFn(fn) }

// IsNonEmpty returns false if the collection is the empty set, else true
//
// Parameters:
//  coll []Value - The collection of elements.
//
// Returns:
//   bool - True if the collection is not empty, else false.
//
// See: https://app.fauna.com/documentation/reference/queryapi#collections
func IsNonEmpty(coll interface{}) Expr { return isNonEmptyFn{IsNonEmpty: wrap(coll)} }

type isNonEmptyFn struct {
	fnApply
	IsNonEmpty Expr `json:"is_nonempty"`
}

func (fn isNonEmptyFn) String() string { return printFn(fn) }

// Contains checks if the provided value contains the path specified.
//
// Parameters:
//  path Path - An array representing a path to check for the existence of. Path can be either strings or ints.
//  value Object - An object to search against.
//
// Returns:
//  bool - true if the path contains any value, false otherwise.
//
// See: https://app.fauna.com/documentation/reference/queryapi#miscellaneous-functions
func Contains(path, value interface{}) Expr {
	return containsFn{Contains: wrap(path), Value: wrap(value)}
}

type containsFn struct {
	fnApply
	Contains Expr `json:"contains"`
	Value    Expr `json:"in"`
}

func (fn containsFn) String() string { return printFn(fn) }

// Count returns the number of elements in the collection.
//
// Parameters:
// collection Expr - the collection
//
// Returns:
// a new Expr instance
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/count
func Count(collection interface{}) Expr {
	return countFn{Count: wrap(collection)}
}

type countFn struct {
	fnApply
	Count Expr `json:"count"`
}

func (fn countFn) String() string { return printFn(fn) }

// Sum sums the elements in the collection.
//
// Parameters:
// collection Expr - the collection
//
// Returns:
// a new Expr instance
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/sum
func Sum(collection interface{}) Expr {
	return sumFn{Sum: wrap(collection)}
}

type sumFn struct {
	fnApply
	Sum Expr `json:"sum"`
}

func (fn sumFn) String() string { return printFn(fn) }

// Mean returns the mean of all elements in the collection.
//
// Parameters:
//
// collection Expr - the collection
//
// Returns:
// a new Expr instance
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/mean
func Mean(collection interface{}) Expr {
	return meanFn{Mean: wrap(collection)}
}

type meanFn struct {
	fnApply
	Mean Expr `json:"mean"`
}

func (fn meanFn) String() string { return printFn(fn) }
