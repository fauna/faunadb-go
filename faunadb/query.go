package faunadb

// Event's action types. Usually used as a parameter for Insert or Remove functions.
//
// See: https://fauna.com/documentation/queries#values-events
const (
	ActionCreate = "create"
	ActionDelete = "delete"
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
// See: https://fauna.com/documentation/queries#values-special_types
func Ref(id string) Expr { return fn1("@ref", id) }

// RefClass creates a new Ref based on the provided class and ID.
//
// See: https://fauna.com/documentation/queries#values-special_types
func RefClass(classRef, id interface{}) Expr { return fn2("ref", classRef, "id", id) }

// Null creates a NullV value.
//
// See: https://fauna.com/documentation/queries#values
func Null() Expr { return NullV{} }

// Basic forms

// Abort aborts the execution of the query
//
// See: https://fauna.com/documentation/queries#basic_forms
func Abort(msg interface{}) Expr { return fn1("abort", msg) }

// Do sequentially evaluates its arguments, and returns the last expression.
// If no expressions are provided, do returns an error.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Do(exprs ...interface{}) Expr { return fn1("do", varargs(exprs...)) }

// If evaluates and returns then or elze depending on the value of cond.
// If cond evaluates to anything other than a boolean, if returns an “invalid argument” error
//
// See: https://fauna.com/documentation/queries#basic_forms
func If(cond, then, elze interface{}) Expr { return fn3("if", cond, "then", then, "else", elze) }

// Lambda creates an anonymous function. Mostly used with Collection functions.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Lambda(varName, expr interface{}) Expr { return fn2("lambda", varName, "expr", expr) }

// At execute an expression at a given timestamp.
//
// See: https://fauna.com/documentation/queries#basic_forms
func At(timestamp, expr interface{}) Expr { return fn2("at", timestamp, "expr", expr) }

// Let binds values to one or more variables.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Let(bindings Obj, in interface{}) Expr { return fn2("let", unescapedBindings(bindings), "in", in) }

// Var refers to a value of a variable on the current lexical scope.
//
// See: https://fauna.com/documentation/queries#basic_forms
func Var(name string) Expr { return fn1("var", name) }

// Call invokes the specified function passing in a variable number of arguments
//
// See: https://fauna.com/documentation/queries#basic_forms
func Call(ref interface{}, args ...interface{}) Expr {
	return fn2("call", ref, "arguments", varargs(args...))
}

// Query creates an instance of the `@query` type with the specified lambda
//
// See: https://fauna.com/documentation/queries#basic_forms
func Query(lambda interface{}) Expr { return fn1("query", lambda) }

// Collections

// Map applies the lambda expression on each element of a collection or Page.
// It returns the result of each application on a collection of the same type.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Map(coll, lambda interface{}) Expr { return fn2("map", lambda, "collection", coll) }

// Foreach applies the lambda expression on each element of a collection or Page.
// The original collection is returned.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Foreach(coll, lambda interface{}) Expr { return fn2("foreach", lambda, "collection", coll) }

// Filter applies the lambda expression on each element of a collection or Page.
// It returns a new collection of the same type containing only the elements in which the
// function application returned true.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Filter(coll, lambda interface{}) Expr { return fn2("filter", lambda, "collection", coll) }

// Take returns a new collection containing num elements from the head of the original collection.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Take(num, coll interface{}) Expr { return fn2("take", num, "collection", coll) }

// Drop returns a new collection containing the remaining elements from the original collection
// after num elements have been removed.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Drop(num, coll interface{}) Expr { return fn2("drop", num, "collection", coll) }

// Prepend returns a new collection that is the result of prepending elems to coll.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Prepend(elems, coll interface{}) Expr { return fn2("prepend", elems, "collection", coll) }

// Append returns a new collection that is the result of appending elems to coll.
//
// See: https://fauna.com/documentation/queries#collection_functions
func Append(elems, coll interface{}) Expr { return fn2("append", elems, "collection", coll) }

// Read

// Get retrieves the instance identified by the provided ref. Optional parameters: TS.
//
// See: https://fauna.com/documentation/queries#read_functions
func Get(ref interface{}, options ...OptionalParameter) Expr { return fn1("get", ref, options...) }

// KeyFromSecret retrieves the key object from the given secret.
//
// See: https://fauna.com/documentation/queries#read_functions
func KeyFromSecret(secret interface{}) Expr { return fn1("key_from_secret", secret) }

// Exists returns boolean true if the provided ref exists (in the case of an instance),
// or is non-empty (in the case of a set), and false otherwise. Optional parameters: TS.
//
// See: https://fauna.com/documentation/queries#read_functions
func Exists(ref interface{}, options ...OptionalParameter) Expr { return fn1("exists", ref, options...) }

// Paginate retrieves a page from the provided set.
// Optional parameters: TS, After, Before, Size, EventsOpt, and Sources.
//
// See: https://fauna.com/documentation/queries#read_functions
func Paginate(set interface{}, options ...OptionalParameter) Expr {
	return fn1("paginate", set, options...)
}

// Write

// Create creates an instance of the specified class.
//
// See: https://fauna.com/documentation/queries#write_functions
func Create(ref, params interface{}) Expr { return fn2("create", ref, "params", params) }

// CreateClass creates a new class.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateClass(params interface{}) Expr { return fn1("create_class", params) }

// CreateDatabase creates an new database.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateDatabase(params interface{}) Expr { return fn1("create_database", params) }

// CreateIndex creates a new index.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateIndex(params interface{}) Expr { return fn1("create_index", params) }

// CreateKey creates a new key.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateKey(params interface{}) Expr { return fn1("create_key", params) }

// CreateFunction creates a new function.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateFunction(params interface{}) Expr { return fn1("create_function", params) }

// Update updates the provided instance.
//
// See: https://fauna.com/documentation/queries#write_functions
func Update(ref, params interface{}) Expr { return fn2("update", ref, "params", params) }

// Replace replaces the provided instance.
//
// See: https://fauna.com/documentation/queries#write_functions
func Replace(ref, params interface{}) Expr { return fn2("replace", ref, "params", params) }

// Delete deletes the provided instance.
//
// See: https://fauna.com/documentation/queries#write_functions
func Delete(ref interface{}) Expr { return fn1("delete", ref) }

// Insert adds an event to the provided instance's history.
//
// See: https://fauna.com/documentation/queries#write_functions
func Insert(ref, ts, action, params interface{}) Expr {
	return fn4("insert", ref, "ts", ts, "action", action, "params", params)
}

// Remove deletes an event from the provided instance's history.
//
// See: https://fauna.com/documentation/queries#write_functions
func Remove(ref, ts, action interface{}) Expr { return fn3("remove", ref, "ts", ts, "action", action) }

// String

// Concat concatenates a list of strings into a single string.
// Optional parameters: Separator.
//
// See: https://fauna.com/documentation/queries#string_functions
func Concat(terms interface{}, options ...OptionalParameter) Expr {
	return fn1("concat", terms, options...)
}

// Casefold normalizes strings according to the Unicode Standard section 5.18 "Case Mappings".
//
// See: https://fauna.com/documentation/queries#string_functions
func Casefold(str interface{}, options ...OptionalParameter) Expr {
	return fn1("casefold", str, options...)
}

// Time and Date

// Time constructs a time from a ISO 8601 offset date/time string.
//
// See: https://fauna.com/documentation/queries#time_functions
func Time(str interface{}) Expr { return fn1("time", str) }

// Date constructs a date from a ISO 8601 offset date/time string.
//
// See: https://fauna.com/documentation/queries#time_functions
func Date(str interface{}) Expr { return fn1("date", str) }

// Epoch constructs a time relative to the epoch "1970-01-01T00:00:00Z".
//
// See: https://fauna.com/documentation/queries#time_functions
func Epoch(num, unit interface{}) Expr { return fn2("epoch", num, "unit", unit) }

// Set

// Singleton returns the history of the instance's presence of the provided ref.
//
// See: https://fauna.com/documentation/queries#sets
func Singleton(ref interface{}) Expr { return fn1("singleton", ref) }

// Events returns the history of instance's data of the provided ref.
//
// See: https://fauna.com/documentation/queries#sets
func Events(refSet interface{}) Expr { return fn1("events", refSet) }

// Match returns the set of instances for the specified index.
//
// See: https://fauna.com/documentation/queries#sets
func Match(ref interface{}) Expr { return fn1("match", ref) }

// MatchTerm returns th set of instances that match the terms in an index.
//
// See: https://fauna.com/documentation/queries#sets
func MatchTerm(ref, terms interface{}) Expr { return fn2("match", ref, "terms", terms) }

// Union returns the set of instances that are present in at least one of the specified sets.
//
// See: https://fauna.com/documentation/queries#sets
func Union(sets ...interface{}) Expr { return fn1("union", varargs(sets...)) }

// Intersection returns the set of instances that are present in all of the specified sets.
//
// See: https://fauna.com/documentation/queries#sets
func Intersection(sets ...interface{}) Expr { return fn1("intersection", varargs(sets...)) }

// Difference returns the set of instances that are present in the source set but not in
// any of the other specified sets.
//
// See: https://fauna.com/documentation/queries#sets
func Difference(sets ...interface{}) Expr { return fn1("difference", varargs(sets...)) }

// Distinct returns the set of instances with duplicates removed.
//
// See: https://fauna.com/documentation/queries#sets
func Distinct(set interface{}) Expr { return fn1("distinct", set) }

// Join derives a set of resources by applying each instance in the source set to the target set.
//
// See: https://fauna.com/documentation/queries#sets
func Join(source, target interface{}) Expr { return fn2("join", source, "with", target) }

// Authentication

// Login creates a token for the provided ref.
//
// See: https://fauna.com/documentation/queries#auth_functions
func Login(ref, params interface{}) Expr { return fn2("login", ref, "params", params) }

// Logout deletes the current session token. If invalidateAll is true, logout will delete all tokens associated with the current session.
//
// See: https://fauna.com/documentation/queries#auth_functions
func Logout(invalidateAll interface{}) Expr { return fn1("logout", invalidateAll) }

// Identify checks the given password against the provided ref's credentials.
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
// See: https://fauna.com/documentation/queries#auth_functions
func Identity() Expr { return fn1("identity", NullV{}) }

// HasIdentity checks if the current key has an identity associated to it.
//
// See: https://fauna.com/documentation/queries#auth_functions
func HasIdentity() Expr { return fn1("has_identity", NullV{}) }

// Miscellaneous

// NextID produces a new identifier suitable for use when constructing refs.
//
// Deprecated: Use NewId instead
//
// See: https://fauna.com/documentation/queries#misc_functions
func NextID() Expr { return fn1("new_id", NullV{}) }

// NewId produces a new identifier suitable for use when constructing refs.
//
// See: https://fauna.com/documentation/queries#misc_functions
func NewId() Expr { return fn1("new_id", NullV{}) }

// Database creates a new database ref.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Database(name interface{}) Expr { return fn1("database", name) }

// ScopedDatabase creates a new database ref inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedDatabase(name interface{}, scope interface{}) Expr {
	return fn2("database", name, "scope", scope)
}

// Index creates a new index ref.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Index(name interface{}) Expr { return fn1("index", name) }

// ScopedIndex creates a new index ref inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedIndex(name interface{}, scope interface{}) Expr { return fn2("index", name, "scope", scope) }

// Class creates a new class ref.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Class(name interface{}) Expr { return fn1("class", name) }

// ScopedClass creates a new class ref inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedClass(name interface{}, scope interface{}) Expr { return fn2("class", name, "scope", scope) }

// Function create a new function ref.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Function(name interface{}) Expr { return fn1("function", name) }

// ScopedFunction creates a new function ref inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedFunction(name interface{}, scope interface{}) Expr {
	return fn2("function", name, "scope", scope)
}

// Classes creates a native ref for classes.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Classes() Expr { return fn1("classes", NullV{}) }

// ScopedClasses creates a native ref for classes inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedClasses(scope interface{}) Expr { return fn1("classes", scope) }

// Indexes creates a native ref for indexes.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Indexes() Expr { return fn1("indexes", NullV{}) }

// ScopedIndexes creates a native ref for indexes inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedIndexes(scope interface{}) Expr { return fn1("indexes", scope) }

// Databases creates a native ref for databases.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Databases() Expr { return fn1("databases", NullV{}) }

// ScopedDatabases creates a native ref for databases inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedDatabases(scope interface{}) Expr { return fn1("databases", scope) }

// Functions creates a native ref for functions.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Functions() Expr { return fn1("functions", NullV{}) }

// ScopedFunctions creates a native ref for functions inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedFunctions(scope interface{}) Expr { return fn1("functions", scope) }

// Keys creates a native ref for keys.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Keys() Expr { return fn1("keys", NullV{}) }

// ScopedKeys creates a native ref for keys inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedKeys(scope interface{}) Expr { return fn1("keys", scope) }

// Tokens creates a native ref for tokens.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Tokens() Expr { return fn1("tokens", NullV{}) }

// ScopedTokens creates a native ref for tokens inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedTokens(scope interface{}) Expr { return fn1("tokens", scope) }

// Credentials creates a native ref for credentials.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Credentials() Expr { return fn1("credentials", NullV{}) }

// ScopedCredentials creates a native ref for credentials inside a database.
//
// See: https://fauna.com/documentation/queries#misc_functions
func ScopedCredentials(scope interface{}) Expr { return fn1("credentials", scope) }

// Equals checks if all args are equivalents.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Equals(args ...interface{}) Expr { return fn1("equals", varargs(args...)) }

// Contains checks if the provided value contains the path specified.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Contains(path, value interface{}) Expr { return fn2("contains", path, "in", value) }

// Add computes the sum of a list of numbers.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Add(args ...interface{}) Expr { return fn1("add", varargs(args...)) }

// Multiply computes the product of a list of numbers.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Multiply(args ...interface{}) Expr { return fn1("multiply", varargs(args...)) }

// Subtract computes the difference of a list of numbers.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Subtract(args ...interface{}) Expr { return fn1("subtract", varargs(args...)) }

// Divide computes the quotient of a list of numbers.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Divide(args ...interface{}) Expr { return fn1("divide", varargs(args...)) }

// Modulo computes the reminder after the division of a list of numbers.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Modulo(args ...interface{}) Expr { return fn1("modulo", varargs(args...)) }

// LT returns true if each specified value is less than all the subsequent values. Otherwise LT returns false.
//
// See: https://fauna.com/documentation/queries#misc_functions
func LT(args ...interface{}) Expr { return fn1("lt", varargs(args...)) }

// LTE returns true if each specified value is less than or equal to all subsequent values. Otherwise LTE returns false.
//
// See: https://fauna.com/documentation/queries#misc_functions
func LTE(args ...interface{}) Expr { return fn1("lte", varargs(args...)) }

// GT returns true if each specified value is greater than all subsequent values. Otherwise GT returns false.
// and false otherwise.
//
// See: https://fauna.com/documentation/queries#misc_functions
func GT(args ...interface{}) Expr { return fn1("gt", varargs(args...)) }

// GTE returns true if each specified value is greater than or equal to all subsequent values. Otherwise GTE returns false.
//
// See: https://fauna.com/documentation/queries#misc_functions
func GTE(args ...interface{}) Expr { return fn1("gte", varargs(args...)) }

// And returns the conjunction of a list of boolean values.
//
// See: https://fauna.com/documentation/queries#misc_functions
func And(args ...interface{}) Expr { return fn1("and", varargs(args...)) }

// Or returnsj the disjunction of a list of boolean values.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Or(args ...interface{}) Expr { return fn1("or", varargs(args...)) }

// Not returns the negation of a boolean value.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Not(boolean interface{}) Expr { return fn1("not", boolean) }

// Select traverses into the provided value, returning the value at the given path.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Select(path, value interface{}, options ...OptionalParameter) Expr {
	return fn2("select", path, "from", value, options...)
}

// SelectAll traverses into the provided value flattening all values under the desired path.
//
// See: https://fauna.com/documentation/queries#misc_functions
func SelectAll(path, value interface{}) Expr {
	return fn2("select_all", path, "from", value)
}
