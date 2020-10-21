package faunadb

// Set

// Singleton returns the history of the document's presence of the provided ref.
//
// Parameters:
//  ref Ref - The reference of the document for which to retrieve the singleton set.
//
// Returns:
//  SetRef - The singleton SetRef.
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func Singleton(ref interface{}) Expr { return singletonFn{Singleton: wrap(ref)} }

type singletonFn struct {
	fnApply
	Singleton Expr `json:"singleton"`
}

// Events returns the history of document's data of the provided ref.
//
// Parameters:
//  refSet Ref|SetRef - A reference or set reference to retrieve an event set from.
//
// Returns:
//  SetRef - The events SetRef.
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func Events(refSet interface{}) Expr { return eventsFn{Events: wrap(refSet)} }

type eventsFn struct {
	fnApply
	Events Expr `json:"events"`
}

// Match returns the set of documents for the specified index.
//
// Parameters:
//  ref Ref - The reference of the index to match against.
//
// Returns:
//  SetRef
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func Match(ref interface{}) Expr { return matchFn{Match: wrap(ref), Terms: nil} }

type matchFn struct {
	fnApply
	Match Expr `json:"match"`
	Terms Expr `json:"terms,omitempty"`
}

// MatchTerm returns th set of documents that match the terms in an index.
//
// Parameters:
//  ref Ref - The reference of the index to match against.
//  terms []Value - A list of terms used in the match.
//
// Returns:
//  SetRef
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func MatchTerm(ref, terms interface{}) Expr { return matchFn{Match: wrap(ref), Terms: wrap(terms)} }

// Union returns the set of documents that are present in at least one of the specified sets.
//
// Parameters:
//  sets []SetRef - A list of SetRef to union together.
//
// Returns:
//  SetRef
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func Union(sets ...interface{}) Expr { return unionFn{Union: wrap(varargs(sets...))} }

type unionFn struct {
	fnApply
	Union Expr `json:"union" faunarepr:"varargs"`
}

// Merge two or more objects..
//
// Parameters:
//   merge merge the first object.
//   with the second object or a list of objects
//   lambda a lambda to resolve possible conflicts
//
// Returns:
// merged object
//
func Merge(merge interface{}, with interface{}, lambda ...OptionalParameter) Expr {
	fn := mergeFn{Merge: wrap(merge), With: wrap(with)}
	return applyOptionals(fn, lambda)
}

type mergeFn struct {
	fnApply
	Merge  Expr `json:"merge"`
	With   Expr `json:"with"`
	Lambda Expr `json:"lambda,omitempty" faunarepr:"optfn,name=ConflictResolver"`
}

// Reduce function applies a reducer Lambda function serially to each member of the collection to produce a single value.
//
// Parameters:
// lambda     Expr  - The accumulator function
// initial    Expr  - The initial value
// collection Expr  - The collection to be reduced
//
// Returns:
// Expr
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/reduce
func Reduce(lambda, initial interface{}, collection interface{}) Expr {
	return reduceFn{
		Reduce:     wrap(lambda),
		Initial:    wrap(initial),
		Collection: wrap(collection),
	}
}

type reduceFn struct {
	fnApply
	Reduce     Expr `json:"reduce"`
	Initial    Expr `json:"initial"`
	Collection Expr `json:"collection"`
}

// Intersection returns the set of documents that are present in all of the specified sets.
//
// Parameters:
//  sets []SetRef - A list of SetRef to intersect.
//
// Returns:
//  SetRef
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func Intersection(sets ...interface{}) Expr {
	return intersectionFn{Intersection: wrap(varargs(sets...))}
}

type intersectionFn struct {
	fnApply
	Intersection Expr `json:"intersection" faunarepr:"varargs"`
}

// Difference returns the set of documents that are present in the first set but not in
// any of the other specified sets.
//
// Parameters:
//  sets []SetRef - A list of SetRef to diff.
//
// Returns:
//  SetRef
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func Difference(sets ...interface{}) Expr { return differenceFn{Difference: wrap(varargs(sets...))} }

type differenceFn struct {
	fnApply
	Difference Expr `json:"difference" faunarepr:"varargs"`
}

// Distinct returns the set of documents with duplicates removed.
//
// Parameters:
//  set []SetRef - A list of SetRef to remove duplicates from.
//
// Returns:
//  SetRef
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func Distinct(set interface{}) Expr { return distinctFn{Distinct: wrap(set)} }

type distinctFn struct {
	fnApply
	Distinct Expr `json:"distinct"`
}

// Join derives a set of resources by applying each document in the source set to the target set.
//
// Parameters:
//  source SetRef - A SetRef of the source set.
//  target Lambda - A Lambda that will accept each element of the source Set and return a Set.
//
// Returns:
//  SetRef
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func Join(source, target interface{}) Expr { return joinFn{Join: wrap(source), With: wrap(target)} }

type joinFn struct {
	fnApply
	Join Expr `json:"join"`
	With Expr `json:"with"`
}

// Range filters the set based on the lower/upper bounds (inclusive).
//
// Parameters:
//  set SetRef - Set to be filtered.
//  from - lower bound.
//  to - upper bound
//
// Returns:
//  SetRef
//
// See: https://app.fauna.com/documentation/reference/queryapi#sets
func Range(set interface{}, from interface{}, to interface{}) Expr {
	return rangeFn{Range: wrap(set), From: wrap(from), To: wrap(to)}
}

type rangeFn struct {
	fnApply
	Range Expr `json:"range"`
	From  Expr `json:"from"`
	To    Expr `json:"to"`
}
