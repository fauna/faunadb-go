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

// Events is an boolean optional parameter that describes if the query should include historical events.
// For more information about events, check https://fauna.com/documentation/queries#values-events.
//
// Functions that accept this optional parameter are: Paginate.
func Events(events interface{}) OptionalParameter {
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

// After is a cursor optional parameter that specifies the next page of a cursor, inclusive.
// For more information about pages, check https://fauna.com/documentation/queries#values-pages.
//
// Functions that accept this optional parameter are: Paginate.
func After(ref interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["after"] = wrap(ref)
	}
}

// Before is a cursor optional parameter that specifies the previous page of a cursor, exclusive.
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

// Default is a optional parameter that specifies the default value for a select operation when
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

// Values

// Ref creates a new RefV value with the ID informed.
//
// See: https://fauna.com/documentation/queries#values-special_types
func Ref(id string) Expr { return RefV{id} }

// RefClass creates a new Ref based on the class and ID informed.
//
// See: https://fauna.com/documentation/queries#values-special_types
func RefClass(classRef, id interface{}) Expr { return fn2("ref", classRef, "id", id) }

// Null creates a NullV value.
//
// See: https://fauna.com/documentation/queries#values
func Null() Expr { return NullV{} }

// Basic forms

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

// Get retrieves the instance identified by the ref informed. Optional parameters: TS.
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

// Paginate retrieves a page from the set informed.
// Optional parameters: TS, After, Before, Size, Events, and Sources.
//
// See: https://fauna.com/documentation/queries#read_functions
func Paginate(set interface{}, options ...OptionalParameter) Expr {
	return fn1("paginate", set, options...)
}

// Write

// Create an instance of the class informed.
//
// See: https://fauna.com/documentation/queries#write_functions
func Create(ref, params interface{}) Expr { return fn2("create", ref, "params", params) }

// CreateClass creates an new class.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateClass(params interface{}) Expr { return fn1("create_class", params) }

// CreateDatabase creates an new database.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateDatabase(params interface{}) Expr { return fn1("create_database", params) }

// CreateIndex creates an new index.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateIndex(params interface{}) Expr { return fn1("create_index", params) }

// CreateKey creates an new key.
//
// See: https://fauna.com/documentation/queries#write_functions
func CreateKey(params interface{}) Expr { return fn1("create_key", params) }

// Update the instance informed.
//
// See: https://fauna.com/documentation/queries#write_functions
func Update(ref, params interface{}) Expr { return fn2("update", ref, "params", params) }

// Replace the instance informed.
//
// See: https://fauna.com/documentation/queries#write_functions
func Replace(ref, params interface{}) Expr { return fn2("replace", ref, "params", params) }

// Delete the instance informed.
//
// See: https://fauna.com/documentation/queries#write_functions
func Delete(ref interface{}) Expr { return fn1("delete", ref) }

// Insert adds an event on an instance's history.
//
// See: https://fauna.com/documentation/queries#write_functions
func Insert(ref, ts, action, params interface{}) Expr {
	return fn4("insert", ref, "ts", ts, "action", action, "params", params)
}

// Remove deletes an event on an instance's history.
//
// See: https://fauna.com/documentation/queries#write_functions
func Remove(ref, ts, action interface{}) Expr { return fn3("remove", ref, "ts", ts, "action", action) }

// String

// Concat joins a list of strings into a single string.
// Optional parameters: Separator.
//
// See: https://fauna.com/documentation/queries#string_functions
func Concat(terms interface{}, options ...OptionalParameter) Expr {
	return fn1("concat", terms, options...)
}

// Casefold normalizes strings according to the Unicode Standard section 5.18 "Case Mappings".
//
// See: https://fauna.com/documentation/queries#string_functions
func Casefold(str interface{}) Expr { return fn1("casefold", str) }

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

// Match returns the set of instances in the ref informed.
//
// See: https://fauna.com/documentation/queries#sets
func Match(ref interface{}) Expr { return fn1("match", ref) }

// MatchTerm returns the set of instances that match the terms informed.
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

// Difference returns the set of instances that are present in the source set and not in
// any of the other specified sets.
//
// See: https://fauna.com/documentation/queries#sets
func Difference(sets ...interface{}) Expr { return fn1("difference", varargs(sets...)) }

// Distinct returns the set after removing duplicates.
//
// See: https://fauna.com/documentation/queries#sets
func Distinct(set interface{}) Expr { return fn1("distinct", set) }

// Join derives a set of resources from target by applying each instance in source to target.
//
// See: https://fauna.com/documentation/queries#sets
func Join(source, target interface{}) Expr { return fn2("join", source, "with", target) }

// Authentication

// Login creates a token for the provided ref.
//
// See: https://fauna.com/documentation/queries#auth_functions
func Login(ref, params interface{}) Expr { return fn2("login", ref, "params", params) }

// Logout deletes all tokens associated with the current session, if invalidateAll is true. Otherwise,
// it deletes only current session token.
//
// See: https://fauna.com/documentation/queries#auth_functions
func Logout(invalidateAll interface{}) Expr { return fn1("logout", invalidateAll) }

// Identify checks the given password against the ref's credentials.
//
// See: https://fauna.com/documentation/queries#auth_functions
func Identify(ref, password interface{}) Expr { return fn2("identify", ref, "password", password) }

// Miscellaneous

// NextID produces a new identifier suitable for use when constructing refs.
//
// See: https://fauna.com/documentation/queries#misc_functions
func NextID() Expr { return fn1("next_id", NullV{}) }

// Database creates a new database ref.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Database(name interface{}) Expr { return fn1("database", name) }

// Index creates a new index ref.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Index(name interface{}) Expr { return fn1("index", name) }

// Class creates a new class ref.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Class(name interface{}) Expr { return fn1("class", name) }

// Equals checks if all args are equivalents.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Equals(args ...interface{}) Expr { return fn1("equals", varargs(args...)) }

// Contains checks if the value informed contains the path specified.
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

// LT returns true if each specified value compares as less than the ones following it,
// and false otherwise.
//
// See: https://fauna.com/documentation/queries#misc_functions
func LT(args ...interface{}) Expr { return fn1("lt", varargs(args...)) }

// LTE returns true if each specified value compares as less than or equal the ones following it,
// and false otherwise.
//
// See: https://fauna.com/documentation/queries#misc_functions
func LTE(args ...interface{}) Expr { return fn1("lte", varargs(args...)) }

// GT returns true if each specified value compares as greater than the ones following it,
// and false otherwise.
//
// See: https://fauna.com/documentation/queries#misc_functions
func GT(args ...interface{}) Expr { return fn1("gt", varargs(args...)) }

// GTE returns true if each specified value compares as greater than or equal the ones following it,
// and false otherwise.
//
// See: https://fauna.com/documentation/queries#misc_functions
func GTE(args ...interface{}) Expr { return fn1("gte", varargs(args...)) }

// And computes the conjunction of a list of boolean values.
//
// See: https://fauna.com/documentation/queries#misc_functions
func And(args ...interface{}) Expr { return fn1("and", varargs(args...)) }

// Or computes the disjunction of a list of boolean values.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Or(args ...interface{}) Expr { return fn1("or", varargs(args...)) }

// Not computes the negation of a boolean value.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Not(boolean interface{}) Expr { return fn1("not", boolean) }

// Select traverses into the value informed returning the value under the desired path.
//
// See: https://fauna.com/documentation/queries#misc_functions
func Select(path, value interface{}, options ...OptionalParameter) Expr {
	return fn2("select", path, "from", value, options...)
}
