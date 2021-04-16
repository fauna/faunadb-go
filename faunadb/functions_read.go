package faunadb

// Read

// Get retrieves the document identified by the provided ref. Optional parameters: TS.
//
// Parameters:
//  ref Ref|SetRef - The reference to the object or a set reference.
//
// Optional parameters:
//  ts time - The snapshot time at which to get the document. See TS() function.
//
// Returns:
//  Object - The object requested.
//
// See: https://app.fauna.com/documentation/reference/queryapi#read-functions
func Get(ref interface{}, options ...OptionalParameter) Expr {
	fn := getFn{Get: wrap(ref)}
	return applyOptionals(fn, options)
}

type getFn struct {
	fnApply
	Get Expr `json:"get"`
	TS  Expr `json:"ts,omitempty" faunarepr:"optfn"`
}

// KeyFromSecret retrieves the key object from the given secret.
//
// Parameters:
//  secret string - The token secret.
//
// Returns:
//  Key - The key object related to the token secret.
//
// See: https://app.fauna.com/documentation/reference/queryapi#read-functions
func KeyFromSecret(secret interface{}) Expr { return keyFromSecretFn{KeyFromSecret: wrap(secret)} }

type keyFromSecretFn struct {
	fnApply
	KeyFromSecret Expr `json:"key_from_secret"`
}

// Exists returns boolean true if the provided ref exists (in the case of an document),
// or is non-empty (in the case of a set), and false otherwise. Optional parameters: TS.
//
// Parameters:
//  ref Ref - The reference to the object. It could be a document reference of a object reference like a collection.
//
// Optional parameters:
//  ts time - The snapshot time at which to check for the document's existence. See TS() function.
//
// Returns:
//  bool - true if the reference exists, false otherwise.
//
// See: https://app.fauna.com/documentation/reference/queryapi#read-functions
func Exists(ref interface{}, options ...OptionalParameter) Expr {
	fn := existsFn{Exists: wrap(ref)}
	return applyOptionals(fn, options)
}

type existsFn struct {
	fnApply
	Exists Expr `json:"exists"`
	TS     Expr `json:"ts,omitempty" faunarepr:"optfn"`
}

// Paginate retrieves a page from the provided set.
//
// Parameters:
//  set SetRef - A set reference to paginate over. See Match() or MatchTerm() functions.
//
// Optional parameters:
//  after Cursor - Return the next page of results after this cursor (inclusive). See After() function.
//  before Cursor - Return the previous page of results before this cursor (exclusive). See Before() function.
//  sources bool - If true, include the source sets along with each element. See Sources() function.
//  ts time - The snapshot time at which to get the document. See TS() function.
//
// Returns:
//  Page - A page of elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#read-functions
func Paginate(set interface{}, options ...OptionalParameter) Expr {
	fn := paginateFn{Paginate: wrap(set)}
	return applyOptionals(fn, options)
}

type paginateFn struct {
	fnApply
	Paginate Expr `json:"paginate"`
	Cursor   Expr `json:"cursor,omitempty" faunarepr:"optfn"`
	After    Expr `json:"after,omitempty" faunarepr:"optfn"`
	Before   Expr `json:"before,omitempty" faunarepr:"optfn"`
	Events   Expr `json:"events,omitempty" faunarepr:"fn=optfn,name=EventsOpt"`
	Size     Expr `json:"size,omitempty" faunarepr:"optfn"`
	Sources  Expr `json:"sources,omitempty" faunarepr:"optfn"`
	TS       Expr `json:"ts,omitempty" faunarepr:"optfn"`
}
