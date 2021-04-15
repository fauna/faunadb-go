package faunadb

// Set

// Singleton returns the history of the document's presence of the
// provided ref.
//
// Parameters:
//  ref Ref - The reference of the document to include in the
//            singleton set.
//
// Returns:
//  SetRef - A reference to the set containing the single document.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/singleton?lang=go
func Singleton(ref interface{}) Expr { return singletonFn{Singleton: wrap(ref)} }

type singletonFn struct {
	fnApply
	Singleton Expr `json:"singleton"`
}

// Events returns the history of document's data of the provided ref.
//
// Parameters:
//  refSet Ref|SetRef - A reference or set reference to retrieve an
//                      event set from.
//
// Returns:
//  SetRef - A reference to the set of events.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/events?lang=go
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
//  SetRef - A reference to the set of matching documents.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/match?lang=go
func Match(ref interface{}) Expr { return matchFn{Match: wrap(ref), Terms: nil} }

type matchFn struct {
	fnApply
	Match Expr `json:"match"`
	Terms Expr `json:"terms,omitempty"`
}

// MatchTerm returns the set of documents that match the terms in an index.
//
// Parameters:
//  ref Ref       - The reference of the index to match against.
//  terms []Value - A list of terms used in the match.
//
// Returns:
//  SetRef - A reference to the set of matching documents.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/match?lang=go
func MatchTerm(ref, terms interface{}) Expr { return matchFn{Match: wrap(ref), Terms: wrap(terms)} }

// Union returns the set of documents that are present in at least one
// of the specified sets.
//
// Parameters:
//  sets []SetRef - A list of SetRef to union together.
//
// Returns:
//  SetRef - A reference to the set of unioned documents.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/union?lang=go
func Union(sets ...interface{}) Expr { return unionFn{Union: wrap(varargs(sets...))} }

type unionFn struct {
	fnApply
	Union Expr `json:"union" faunarepr:"varargs"`
}

// Merge two or more objects.
//
// Parameters:
//  merge  Object   - The first object to merge.
//  with   Object   - The second object, or a list of objects, to merge.
//  lambda Function - A lambda to resolve possible conflicts.
//
// Returns:
//  object - A new object representing the merged objects.
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

// Reduce function applies a reducer Lambda function serially to each
// member of the collection to produce a single value.
//
// Parameters:
//  lambda     Expr  - The accumulator function
//  initial    Expr  - The initial value
//  collection Expr  - The collection to be reduced
//
// Returns:
//  Expr - The result of the reducer.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/reduce?lang=go
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

// Intersection returns the set of documents that are present in all of
// the specified sets.
//
// Parameters:
//  sets []SetRef - A list of SetRef to intersect.
//
// Returns:
//  SetRef - A reference to the set of intersected documents.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/intersection?lang=go
func Intersection(sets ...interface{}) Expr {
	return intersectionFn{Intersection: wrap(varargs(sets...))}
}

type intersectionFn struct {
	fnApply
	Intersection Expr `json:"intersection" faunarepr:"varargs"`
}

// Difference returns the set of documents that are present in the first
// set but not in any of the other specified sets.
//
// Parameters:
//  sets []SetRef - A list of SetRef to diff.
//
// Returns:
//  SetRef - A reference to the set of differenced documents.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/difference?lang=go
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
//  SetRef - A reference to the set of distinct documents.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/distinct?lang=go
func Distinct(set interface{}) Expr { return distinctFn{Distinct: wrap(set)} }

type distinctFn struct {
	fnApply
	Distinct Expr `json:"distinct"`
}

// Join derives a set of resources by applying each document in the
// source set to the target set.
//
// Parameters:
//  source SetRef - A SetRef of the source set.
//  target Lambda - A Lambda that will accept each element of the source
//                  Set and return a Set.
//
// Returns:
//  SetRef - A reference to the set containing joined documents.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/join?lang=go
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
//  from       - The lower bound of the range.
//  to         - The upper bound of the range.
//
// Returns:
//  SetRef - A reference to the set containing documents in the specific
//           range.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/range?lang=go
func Range(set interface{}, from interface{}, to interface{}) Expr {
	return rangeFn{Range: wrap(set), From: wrap(from), To: wrap(to)}
}

type rangeFn struct {
	fnApply
	Range Expr `json:"range"`
	From  Expr `json:"from"`
	To    Expr `json:"to"`
}
