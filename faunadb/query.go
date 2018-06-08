package faunadb

// Event's action types. Usually used as a parameter for Insert or Remove functions.
//
// See: https://fauna.com/documentation/queries#values-events
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionAdd    = "add"
	ActionRemove = "remove"
)

// Time unit. Usually used as a parameter for Epoch function.
//
// See: https://fauna.com/documentation/queries#time_functions-epoch_num_unit_unit
const (
	TimeUnitSecond      = "second"
	TimeUnitMillisecond = "millisecond"
	TimeUnitMicrosecond = "microsecond"
	TimeUnitNanosecond  = "nanosecond"
)

// Normalizers for Casefold
//
// See: https://fauna.com/documentation/queries#string_functions
const (
	NormalizerNFKCCaseFold = "NFKCCaseFold"
	NormalizerNFC          = "NFC"
	NormalizerNFD          = "NFD"
	NormalizerNFKC         = "NFKC"
	NormalizerNFKD         = "NFKD"
)

// Helper functions

func varargs(expr ...interface{}) interface{} {
	if len(expr) == 1 {
		return expr[0]
	}

	return expr
}

func unescapedBindings(obj Obj) unescapedObj {
	res := make(unescapedObj, len(obj))

	for k, v := range obj {
		res[k] = wrap(v)
	}

	return res
}

// Optional parameters

// EventsOpt is an boolean optional parameter that describes if the query should include historical events.
// For more information about events, check https://fauna.com/documentation/queries#values-events.
//
// Functions that accept this optional parameter are: Paginate.
//
// Deprecated: The Events function was renamed to EventsOpt to support the new history API.
// EventsOpt is provided here for backwards compatibility. Instead of using Paginate with the EventsOpt parameter,
// you should use the new Events function.
func EventsOpt(events interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["events"] = wrap(events)
	}
}

// TS is a timestamp optional parameter that specifies in which timestamp a query should be executed.
//
// Functions that accept this optional parameter are: Get, Insert, Remove, Exists, and Paginate.
func TS(timestamp interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["ts"] = wrap(timestamp)
	}
}

// After is an optional parameter used when cursoring that refers to the specified cursor's the next page, inclusive.
// For more information about pages, check https://fauna.com/documentation/queries#values-pages.
//
// Functions that accept this optional parameter are: Paginate.
func After(ref interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["after"] = wrap(ref)
	}
}

// Before is an optional parameter used when cursoring that refers to the specified cursor's previous page, exclusive.
// For more information about pages, check https://fauna.com/documentation/queries#values-pages.
//
// Functions that accept this optional parameter are: Paginate.
func Before(ref interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["before"] = wrap(ref)
	}
}

// Size is a numeric optional parameter that specifies the size of a pagination cursor.
//
// Functions that accept this optional parameter are: Paginate.
func Size(size interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["size"] = wrap(size)
	}
}

// Sources is a boolean optional parameter that specifies if a pagination cursor should include
// the source sets along with each element.
//
// Functions that accept this optional parameter are: Paginate.
func Sources(sources interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["sources"] = wrap(sources)
	}
}

// Default is an optional parameter that specifies the default value for a select operation when
// the desired value path is absent.
//
// Functions that accept this optional parameter are: Select.
func Default(value interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["default"] = wrap(value)
	}
}

// Separator is a string optional parameter that specifies the separator for a concat operation.
//
// Functions that accept this optional parameter are: Concat.
func Separator(sep interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["separator"] = wrap(sep)
	}
}

// Normalizer is a string optional parameter that specifies the normalization function for casefold operation.
//
// Functions that accept this optional parameter are: Casefold.
func Normalizer(norm interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["normalizer"] = wrap(norm)
	}
}

// Values

// Ref creates a new RefV value with the provided ID.
//
// Parameters:
//  id string - A string representation of a reference type.
//
// Returns:
//  Ref - A new reference type.
//
// See: https://fauna.com/documentation/queries#values-special_types
func Ref(id string) Expr { return fn1("@ref", id) }

// RefClass creates a new Ref based on the provided class and ID.
//
// Parameters:
//  classRef Ref - A class reference.
//  id string|int64 - The instance ID.
//
// Returns:
//  Ref - A new reference type.
//
// See: https://fauna.com/documentation/queries#values-special_types
func RefClass(classRef, id interface{}) Expr { return fn2("ref", classRef, "id", id) }

// Null creates a NullV value.
//
// Returns:
//  Value - A null value.
//
// See: https://fauna.com/documentation/queries#values
func Null() Expr { return NullV{} }

// Basic forms

// Abort aborts the execution of the query
//
// Parameters:
//  msg string - An error message.
//
// Returns:
//  Error
//
// See: https://fauna.com/documentation/queries#basic_forms
func Abort(msg interface{}) Expr { return fn1("abort", msg) }

// Do sequentially evaluates its arguments, and returns the last expression.
// If no expressions are provided, do returns an error.
//
// Parameters:
//  exprs []Expr - A variable number of expressions to be evaluated.
//
// Returns:
//  Value - The result of the last expression in the list.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Do(exprs ...interface{}) Expr { return fn1("do", varargs(exprs...)) }

// If evaluates and returns then or elze depending on the value of cond.
// If cond evaluates to anything other than a boolean, if returns an “invalid argument” error
//
// Parameters:
//  cond bool - A boolean expression.
//  then Expr - The expression to run if condition is true.
//  elze Expr - The expression to run if condition is false.
//
// Returns:
//  Value - The result of either then or elze expression.
//
// See: https://fauna.com/documentation/queries#basic_forms
func If(cond, then, elze interface{}) Expr { return fn3("if", cond, "then", then, "else", elze) }

// Lambda creates an anonymous function. Mostly used with Collection functions.
//
// Parameters:
//  varName string|[]string - A string or an array of strings of arguments name to be bound in the body of the lambda.
//  expr Expr - An expression used as the body of the lambda.
//
// Returns:
//  Value - The result of the body expression.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Lambda(varName, expr interface{}) Expr { return fn2("lambda", varName, "expr", expr) }

// At execute an expression at a given timestamp.
//
// Parameters:
//  timestamp time - The timestamp in which the expression will be evaluated.
//  expr Expr - An expression to be evaluated.
//
// Returns:
//  Value - The result of the given expression.
//
// See: https://fauna.com/documentation/queries#basic_forms
func At(timestamp, expr interface{}) Expr { return fn2("at", timestamp, "expr", expr) }

// Let binds values to one or more variables.
//
// Parameters:
//  bindings Object - An object binding a variable name to a value.
//  in Expr - An expression to be evaluated.
//
// Returns:
//  Value - The result of the given expression.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Let(bindings Obj, in interface{}) Expr { return fn2("let", unescapedBindings(bindings), "in", in) }

// Var refers to a value of a variable on the current lexical scope.
//
// Parameters:
//  name string - The variable name.
//
// Returns:
//  Value - The variable value.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Var(name string) Expr { return fn1("var", name) }

// Call invokes the specified function passing in a variable number of arguments
//
// Parameters:
//  ref Ref - The reference to the user defined functions to call.
//  args []Value - A series of values to pass as arguments to the user defined function.
//
// Returns:
//  Value - The return value of the user defined function.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Call(ref interface{}, args ...interface{}) Expr {
	return fn2("call", ref, "arguments", varargs(args...))
}

// Query creates an instance of the @query type with the specified lambda
//
// Parameters:
//  lambda Lambda - A lambda representation. See Lambda() function.
//
// Returns:
//  Query - The lambda wrapped in a @query type.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Query(lambda interface{}) Expr { return fn1("query", lambda) }

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
// See: https://fauna.com/documentation/queries#collection_functions
func Map(coll, lambda interface{}) Expr { return fn2("map", lambda, "collection", coll) }

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
// See: https://fauna.com/documentation/queries#collection_functions
func Foreach(coll, lambda interface{}) Expr { return fn2("foreach", lambda, "collection", coll) }

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
// See: https://fauna.com/documentation/queries#collection_functions
func Filter(coll, lambda interface{}) Expr { return fn2("filter", lambda, "collection", coll) }

// Take returns a new collection containing num elements from the head of the original collection.
//
// Parameters:
//  num int64 - The number of elements to take from the collection.
//  coll []Value - The collection of elements.
//
// Returns:
//  []Value - A new collection.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Take(num, coll interface{}) Expr { return fn2("take", num, "collection", coll) }

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
// See: https://fauna.com/documentation/queries#collection_functions
func Drop(num, coll interface{}) Expr { return fn2("drop", num, "collection", coll) }

// Prepend returns a new collection that is the result of prepending elems to coll.
//
// Parameters:
//  elems []Value - Elements to add to the beginning of the other collection.
//  coll []Value - The collection of elements.
//
// Returns:
//  []Value - A new collection.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Prepend(elems, coll interface{}) Expr { return fn2("prepend", elems, "collection", coll) }

// Append returns a new collection that is the result of appending elems to coll.
//
// Parameters:
//  elems []Value - Elements to add to the end of the other collection.
//  coll []Value - The collection of elements.
//
// Returns:
//  []Value - A new collection.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Append(elems, coll interface{}) Expr { return fn2("append", elems, "collection", coll) }

// Read

// Get retrieves the instance identified by the provided ref. Optional parameters: TS.
//
// Parameters:
//  ref Ref|SetRef - The reference to the object or a set reference.
//
// Optional parameters:
//  ts time - The snapshot time at which to get the instance. See TS() function.
//
// Returns:
//  Object - The object requested.
//
// See: https://fauna.com/documentation/queries#read_functions
func Get(ref interface{}, options ...OptionalParameter) Expr { return fn1("get", ref, options...) }

// KeyFromSecret retrieves the key object from the given secret.
//
// Parameters:
//  secret string - The token secret.
//
// Returns:
//  Key - The key object related to the token secret.
//
// See: https://fauna.com/documentation/queries#read_functions
func KeyFromSecret(secret interface{}) Expr { return fn1("key_from_secret", secret) }

// Exists returns boolean true if the provided ref exists (in the case of an instance),
// or is non-empty (in the case of a set), and false otherwise. Optional parameters: TS.
//
// Parameters:
//  ref Ref - The reference to the object. It could be an instance reference of a object reference like a class.
//
// Optional parameters:
//  ts time - The snapshot time at which to check for the instance's existence. See TS() function.
//
// Returns:
//  bool - true if the reference exists, false otherwise.
//
// See: https://fauna.com/documentation/queries#read_functions
func Exists(ref interface{}, options ...OptionalParameter) Expr { return fn1("exists", ref, options...) }

// Paginate retrieves a page from the provided set.
//
// Parameters:
//  set SetRef - A set reference to paginate over. See Match() or MatchTerm() functions.
//
// Optional parameters:
//  after Cursor - Return the next page of results after this cursor (inclusive). See After() function.
//  before Cursor - Return the previous page of results before this cursor (exclusive). See Before() function.
//  sources bool - If true, include the source sets along with each element. See Sources() function.
//  ts time - The snapshot time at which to get the instance. See TS() function.
//
// Returns:
//  Page - A page of elements.
//
// See: https://fauna.com/documentation/queries#read_functions
func Paginate(set interface{}, options ...OptionalParameter) Expr {
	return fn1("paginate", set, options...)
}

// Write

// Create creates an instance of the specified class.
//
// Parameters:
//  ref Ref - A class reference.
//  params Object - An object with attributes of the instance created.
//
// Returns:
//  Object - A new instance of the class referenced.
//
// See: https://fauna.com/documentation/queries#write_functions
func Create(ref, params interface{}) Expr { return fn2("create", ref, "params", params) }

// CreateClass creates a new class.
//
// Parameters:
//  params Object - An object with attributes of the class.
//
// Returns:
//  Object - The new created class object.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateClass(params interface{}) Expr { return fn1("create_class", params) }

// CreateDatabase creates an new database.
//
// Parameters:
//  params Object - An object with attributes of the database.
//
// Returns:
//  Object - The new created database object.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateDatabase(params interface{}) Expr { return fn1("create_database", params) }

// CreateIndex creates a new index.
//
// Parameters:
//  params Object - An object with attributes of the index.
//
// Returns:
//  Object - The new created index object.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateIndex(params interface{}) Expr { return fn1("create_index", params) }

// CreateKey creates a new key.
//
// Parameters:
//  params Object - An object with attributes of the key.
//
// Returns:
//  Object - The new created key object.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateKey(params interface{}) Expr { return fn1("create_key", params) }

// CreateFunction creates a new function.
//
// Parameters:
//  params Object - An object with attributes of the function.
//
// Returns:
//  Object - The new created function object.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateFunction(params interface{}) Expr { return fn1("create_function", params) }

// Update updates the provided instance.
//
// Parameters:
//  ref Ref - The reference to update.
//  params Object - An object representing the parameters of the instance.
//
// Returns:
//  Object - The updated object.
//
// See: https://fauna.com/documentation/queries#write_functions
func Update(ref, params interface{}) Expr { return fn2("update", ref, "params", params) }

// Replace replaces the provided instance.
//
// Parameters:
//  ref Ref - The reference to replace.
//  params Object - An object representing the parameters of the instance.
//
// Returns:
//  Object - The replaced object.
//
// See: https://fauna.com/documentation/queries#write_functions
func Replace(ref, params interface{}) Expr { return fn2("replace", ref, "params", params) }

// Delete deletes the provided instance.
//
// Parameters:
//  ref Ref - The reference to delete.
//
// Returns:
//  Object - The deleted object.
//
// See: https://fauna.com/documentation/queries#write_functions
func Delete(ref interface{}) Expr { return fn1("delete", ref) }

// Insert adds an event to the provided instance's history.
//
// Parameters:
//  ref Ref - The reference to insert against.
//  ts time - The valid time of the inserted event.
//  action string - Whether the event shoulde be a ActionCreate, ActionUpdate or ActionDelete.
//
// Returns:
//  Object - The deleted object.
//
// See: https://fauna.com/documentation/queries#write_functions
func Insert(ref, ts, action, params interface{}) Expr {
	return fn4("insert", ref, "ts", ts, "action", action, "params", params)
}

// Remove deletes an event from the provided instance's history.
//
// Parameters:
//  ref Ref - The reference of the instance whose event should be removed.
//  ts time - The valid time of the inserted event.
//  action string - The event action (ActionCreate, ActionUpdate or ActionDelete) that should be removed.
//
// Returns:
//  Object - The deleted object.
//
// See: https://fauna.com/documentation/queries#write_functions
func Remove(ref, ts, action interface{}) Expr { return fn3("remove", ref, "ts", ts, "action", action) }

// String

// Concat concatenates a list of strings into a single string.
//
// Parameters:
//  terms []string - A list of strings to concatenate.
//
// Optional parameters:
//  separator string - The separator to use between each string. See Separator() function.
//
// Returns:
//  string - A string with all terms concatenated.
//
// See: https://fauna.com/documentation/queries#string_functions
func Concat(terms interface{}, options ...OptionalParameter) Expr {
	return fn1("concat", terms, options...)
}

// Casefold normalizes strings according to the Unicode Standard section 5.18 "Case Mappings".
//
// Parameters:
//  str string - The string to casefold.
//
// Optional parameters:
//  normalizer string - The algorithm to use. One of: NormalizerNFKCCaseFold, NormalizerNFC, NormalizerNFD, NormalizerNFKC, NormalizerNFKD.
//
// Returns:
//  string - The normalized string.
//
// See: https://fauna.com/documentation/queries#string_functions
func Casefold(str interface{}, options ...OptionalParameter) Expr {
	return fn1("casefold", str, options...)
}

// Time and Date

// Time constructs a time from a ISO 8601 offset date/time string.
//
// Parameters:
//  str string - A string to convert to a time object.
//
// Returns:
//  time - A time object.
//
// See: https://fauna.com/documentation/queries#time_functions
func Time(str interface{}) Expr { return fn1("time", str) }

// Date constructs a date from a ISO 8601 offset date/time string.
//
// Parameters:
//  str string - A string to convert to a date object.
//
// Returns:
//  date - A date object.
//
// See: https://fauna.com/documentation/queries#time_functions
func Date(str interface{}) Expr { return fn1("date", str) }

// Epoch constructs a time relative to the epoch "1970-01-01T00:00:00Z".
//
// Parameters:
//  num int64 - The number of units from Epoch.
//  unit string - The unit of number. One of TimeUnitSecond, TimeUnitMillisecond, TimeUnitMicrosecond, TimeUnitNanosecond.
//
// Returns:
//  time - A time object.
//
// See: https://fauna.com/documentation/queries#time_functions
func Epoch(num, unit interface{}) Expr { return fn2("epoch", num, "unit", unit) }

// Set

// Singleton returns the history of the instance's presence of the provided ref.
//
// Parameters:
//  ref Ref - The reference of the instance for which to retrieve the singleton set.
//
// Returns:
//  SetRef - The singleton SetRef.
//
// See: https://fauna.com/documentation/queries#sets
func Singleton(ref interface{}) Expr { return fn1("singleton", ref) }

// Events returns the history of instance's data of the provided ref.
//
// Parameters:
//  refSet Ref|SetRef - A reference or set reference to retrieve an event set from.
//
// Returns:
//  SetRef - The events SetRef.
//
// See: https://fauna.com/documentation/queries#sets
func Events(refSet interface{}) Expr { return fn1("events", refSet) }

// Match returns the set of instances for the specified index.
//
// Parameters:
//  ref Ref - The reference of the index to match against.
//
// Returns:
//  SetRef
//
// See: https://fauna.com/documentation/queries#sets
func Match(ref interface{}) Expr { return fn1("match", ref) }

// MatchTerm returns th set of instances that match the terms in an index.
//
// Parameters:
//  ref Ref - The reference of the index to match against.
//  terms []Value - A list of terms used in the match.
//
// Returns:
//  SetRef
//
// See: https://fauna.com/documentation/queries#sets
func MatchTerm(ref, terms interface{}) Expr { return fn2("match", ref, "terms", terms) }

// Union returns the set of instances that are present in at least one of the specified sets.
//
// Parameters:
//  sets []SetRef - A list of SetRef to union together.
//
// Returns:
//  SetRef
//
// See: https://fauna.com/documentation/queries#sets
func Union(sets ...interface{}) Expr { return fn1("union", varargs(sets...)) }

// Intersection returns the set of instances that are present in all of the specified sets.
//
// Parameters:
//  sets []SetRef - A list of SetRef to intersect.
//
// Returns:
//  SetRef
//
// See: https://fauna.com/documentation/queries#sets
func Intersection(sets ...interface{}) Expr { return fn1("intersection", varargs(sets...)) }

// Difference returns the set of instances that are present in the first set but not in
// any of the other specified sets.
//
// Parameters:
//  sets []SetRef - A list of SetRef to diff.
//
// Returns:
//  SetRef
//
// See: https://fauna.com/documentation/queries#sets
func Difference(sets ...interface{}) Expr { return fn1("difference", varargs(sets...)) }

// Distinct returns the set of instances with duplicates removed.
//
// Parameters:
//  set []SetRef - A list of SetRef to remove duplicates from.
//
// Returns:
//  SetRef
//
// See: https://fauna.com/documentation/queries#sets
func Distinct(set interface{}) Expr { return fn1("distinct", set) }

// Join derives a set of resources by applying each instance in the source set to the target set.
//
// Parameters:
//  source SetRef - A SetRef of the source set.
//  target Lambda - A Lambda that will accept each element of the source Set and return a Set.
//
// Returns:
//  SetRef
//
// See: https://fauna.com/documentation/queries#sets
func Join(source, target interface{}) Expr { return fn2("join", source, "with", target) }

// Authentication

// Login creates a token for the provided ref.
//
// Parameters:
//  ref Ref - A reference with credentials to authenticate against.
//  params Object - An object of parameters to pass to the login function
//    - password: The password used to login
//
// Returns:
//  Key - a key with the secret to login.
//
// See: https://fauna.com/documentation/queries#auth_functions
func Login(ref, params interface{}) Expr { return fn2("login", ref, "params", params) }

// Logout deletes the current session token. If invalidateAll is true, logout will delete all tokens associated with the current session.
//
// Parameters:
//  invalidateAll bool - If true, log out all tokens associated with the current session.
//
// See: https://fauna.com/documentation/queries#auth_functions
func Logout(invalidateAll interface{}) Expr { return fn1("logout", invalidateAll) }

// Identify checks the given password against the provided ref's credentials.
//
// Parameters:
//  ref Ref - The reference to check the password against.
//  password string - The credentials password to check.
//
// Returns:
//  bool - true if the password is correct, false otherwise.
//
// See: https://fauna.com/documentation/queries#auth_functions
func Identify(ref, password interface{}) Expr { return fn2("identify", ref, "password", password) }

// Identity returns the instance reference associated with the current key.
//
// For example, the current key token created using:
//	Create(Tokens(), Obj{"instance": someRef})
// or via:
//	Login(someRef, Obj{"password":"sekrit"})
// will return "someRef" as the result of this function.
//
// Returns:
//  Ref - The reference associated with the current key.
//
// See: https://fauna.com/documentation/queries#auth_functions
func Identity() Expr { return fn1("identity", NullV{}) }

// HasIdentity checks if the current key has an identity associated to it.
//
// Returns:
//  bool - true if the current key has an identity, false otherwise.
//
// See: https://fauna.com/documentation/queries#auth_functions
func HasIdentity() Expr { return fn1("has_identity", NullV{}) }

// Miscellaneous

// NextID produces a new identifier suitable for use when constructing refs.
//
// Deprecated: Use NewId instead
//
// Returns:
//  string - The new ID.
//
// See: https://fauna.com/documentation/queries#misc_functions
func NextID() Expr { return fn1("new_id", NullV{}) }

// NewId produces a new identifier suitable for use when constructing refs.
//
// Returns:
//  string - The new ID.
//
// See: https://fauna.com/documentation/queries#misc_functions
func NewId() Expr { return fn1("new_id", NullV{}) }

// Database creates a new database ref.
//
// Parameters:
//  name string - The name of the database.
//
// Returns:
//  Ref - The database reference.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Database(name interface{}) Expr { return fn1("database", name) }

// ScopedDatabase creates a new database ref inside a database.
//
// Parameters:
//  name string - The name of the database.
//  scope Ref - The reference of the database's scope.
//
// Returns:
//  Ref - The database reference.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedDatabase(name interface{}, scope interface{}) Expr {
	return fn2("database", name, "scope", scope)
}

// Index creates a new index ref.
//
// Parameters:
//  name string - The name of the index.
//
// Returns:
//  Ref - The index reference.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Index(name interface{}) Expr { return fn1("index", name) }

// ScopedIndex creates a new index ref inside a database.
//
// Parameters:
//  name string - The name of the index.
//  scope Ref - The reference of the index's scope.
//
// Returns:
//  Ref - The index reference.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedIndex(name interface{}, scope interface{}) Expr { return fn2("index", name, "scope", scope) }

// Class creates a new class ref.
//
// Parameters:
//  name string - The name of the class.
//
// Returns:
//  Ref - The class reference.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Class(name interface{}) Expr { return fn1("class", name) }

// ScopedClass creates a new class ref inside a database.
//
// Parameters:
//  name string - The name of the class.
//  scope Ref - The reference of the class's scope.
//
// Returns:
//  Ref - The class reference.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedClass(name interface{}, scope interface{}) Expr { return fn2("class", name, "scope", scope) }

// Function create a new function ref.
//
// Parameters:
//  name string - The name of the functions.
//
// Returns:
//  Ref - The function reference.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Function(name interface{}) Expr { return fn1("function", name) }

// ScopedFunction creates a new function ref inside a database.
//
// Parameters:
//  name string - The name of the function.
//  scope Ref - The reference of the function's scope.
//
// Returns:
//  Ref - The function reference.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedFunction(name interface{}, scope interface{}) Expr {
	return fn2("function", name, "scope", scope)
}

// Classes creates a native ref for classes.
//
// Returns:
//  Ref - The reference of the class set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Classes() Expr { return fn1("classes", NullV{}) }

// ScopedClasses creates a native ref for classes inside a database.
//
// Parameters:
//  scope Ref - The reference of the class set's scope.
//
// Returns:
//  Ref - The reference of the class set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedClasses(scope interface{}) Expr { return fn1("classes", scope) }

// Indexes creates a native ref for indexes.
//
// Returns:
//  Ref - The reference of the index set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Indexes() Expr { return fn1("indexes", NullV{}) }

// ScopedIndexes creates a native ref for indexes inside a database.
//
// Parameters:
//  scope Ref - The reference of the index set's scope.
//
// Returns:
//  Ref - The reference of the index set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedIndexes(scope interface{}) Expr { return fn1("indexes", scope) }

// Databases creates a native ref for databases.
//
// Returns:
//  Ref - The reference of the datbase set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Databases() Expr { return fn1("databases", NullV{}) }

// ScopedDatabases creates a native ref for databases inside a database.
//
// Parameters:
//  scope Ref - The reference of the database set's scope.
//
// Returns:
//  Ref - The reference of the database set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedDatabases(scope interface{}) Expr { return fn1("databases", scope) }

// Functions creates a native ref for functions.
//
// Returns:
//  Ref - The reference of the function set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Functions() Expr { return fn1("functions", NullV{}) }

// ScopedFunctions creates a native ref for functions inside a database.
//
// Parameters:
//  scope Ref - The reference of the function set's scope.
//
// Returns:
//  Ref - The reference of the function set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedFunctions(scope interface{}) Expr { return fn1("functions", scope) }

// Keys creates a native ref for keys.
//
// Returns:
//  Ref - The reference of the key set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Keys() Expr { return fn1("keys", NullV{}) }

// ScopedKeys creates a native ref for keys inside a database.
//
// Parameters:
//  scope Ref - The reference of the key set's scope.
//
// Returns:
//  Ref - The reference of the key set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedKeys(scope interface{}) Expr { return fn1("keys", scope) }

// Tokens creates a native ref for tokens.
//
// Returns:
//  Ref - The reference of the token set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Tokens() Expr { return fn1("tokens", NullV{}) }

// ScopedTokens creates a native ref for tokens inside a database.
//
// Parameters:
//  scope Ref - The reference of the token set's scope.
//
// Returns:
//  Ref - The reference of the token set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedTokens(scope interface{}) Expr { return fn1("tokens", scope) }

// Credentials creates a native ref for credentials.
//
// Returns:
//  Ref - The reference of the credential set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Credentials() Expr { return fn1("credentials", NullV{}) }

// ScopedCredentials creates a native ref for credentials inside a database.
//
// Parameters:
//  scope Ref - The reference of the credential set's scope.
//
// Returns:
//  Ref - The reference of the credential set.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedCredentials(scope interface{}) Expr { return fn1("credentials", scope) }

// Equals checks if all args are equivalents.
//
// Parameters:
//  args []Value - A collection of expressions to check for equivalence.
//
// Returns:
//  bool - true if all elements are equals, false otherwise.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Equals(args ...interface{}) Expr { return fn1("equals", varargs(args...)) }

// Contains checks if the provided value contains the path specified.
//
// Parameters:
//  path Path - An array representing a path to check for the existence of. Path can be either strings or ints.
//  value Object - An object to search against.
//
// Returns:
//  bool - true if the path contains any value, false otherwise.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Contains(path, value interface{}) Expr { return fn2("contains", path, "in", value) }

// Add computes the sum of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to sum together.
//
// Returns:
//  number - The sum of all elements.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Add(args ...interface{}) Expr { return fn1("add", varargs(args...)) }

// Multiply computes the product of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to multiply together.
//
// Returns:
//  number - The multiplication of all elements.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Multiply(args ...interface{}) Expr { return fn1("multiply", varargs(args...)) }

// Subtract computes the difference of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to compute the difference of.
//
// Returns:
//  number - The difference of all elements.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Subtract(args ...interface{}) Expr { return fn1("subtract", varargs(args...)) }

// Divide computes the quotient of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to compute the quotient of.
//
// Returns:
//  number - The quotient of all elements.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Divide(args ...interface{}) Expr { return fn1("divide", varargs(args...)) }

// Modulo computes the reminder after the division of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to compute the quotient of. The remainder will be returned.
//
// Returns:
//  number - The remainder of the quotient of all elements.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Modulo(args ...interface{}) Expr { return fn1("modulo", varargs(args...)) }

// LT returns true if each specified value is less than all the subsequent values. Otherwise LT returns false.
//
// Parameters:
//  args []number - A collection of terms to compare.
//
// Returns:
//  bool - true if all elements are less than each other from left to right.
//
// See: https://fauna.com/documentation/queries#misc_functions
func LT(args ...interface{}) Expr { return fn1("lt", varargs(args...)) }

// LTE returns true if each specified value is less than or equal to all subsequent values. Otherwise LTE returns false.
//
// Parameters:
//  args []number - A collection of terms to compare.
//
// Returns:
//  bool - true if all elements are less than of equals each other from left to right.
//
// See: https://fauna.com/documentation/queries#misc_functions
func LTE(args ...interface{}) Expr { return fn1("lte", varargs(args...)) }

// GT returns true if each specified value is greater than all subsequent values. Otherwise GT returns false.
// and false otherwise.
//
// Parameters:
//  args []number - A collection of terms to compare.
//
// Returns:
//  bool - true if all elements are greather than to each other from left to right.
//
// See: https://fauna.com/documentation/queries#misc_functions
func GT(args ...interface{}) Expr { return fn1("gt", varargs(args...)) }

// GTE returns true if each specified value is greater than or equal to all subsequent values. Otherwise GTE returns false.
//
// Parameters:
//  args []number - A collection of terms to compare.
//
// Returns:
//  bool - true if all elements are greather than or equals to each other from left to right.
//
// See: https://fauna.com/documentation/queries#misc_functions
func GTE(args ...interface{}) Expr { return fn1("gte", varargs(args...)) }

// And returns the conjunction of a list of boolean values.
//
// Parameters:
//  args []bool - A collection to compute the conjunction of.
//
// Returns:
//  bool - true if all elements are true, false otherwise.
//
// See: https://fauna.com/documentation/queries#misc_functions
func And(args ...interface{}) Expr { return fn1("and", varargs(args...)) }

// Or returns the disjunction of a list of boolean values.
//
// Parameters:
//  args []bool - A collection to compute the disjunction of.
//
// Returns:
//  bool - true if at least one element is true, false otherwise.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Or(args ...interface{}) Expr { return fn1("or", varargs(args...)) }

// Not returns the negation of a boolean value.
//
// Parameters:
//  boolean bool - A boolean to produce the negation of.
//
// Returns:
//  bool - The value negated.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Not(boolean interface{}) Expr { return fn1("not", boolean) }

// Select traverses into the provided value, returning the value at the given path.
//
// Parameters:
//  path []Path - An array representing a path to pull from an object. Path can be either strings or numbers.
//  value Object - The object to select from.
//
// Optional parameters:
//  default Value - A default value if the path does not exist. See Default() function.
//
// Returns:
//  Value - The value at the given path location.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Select(path, value interface{}, options ...OptionalParameter) Expr {
	return fn2("select", path, "from", value, options...)
}

// SelectAll traverses into the provided value flattening all values under the desired path.
//
// Parameters:
//  path []Path - An array representing a path to pull from an object. Path can be either strings or numbers.
//  value Object - The object to select from.
//
// Returns:
//  Value - The value at the given path location.
//
// See: https://fauna.com/documentation/queries#misc_functions
func SelectAll(path, value interface{}) Expr {
	return fn2("select_all", path, "from", value)
}
