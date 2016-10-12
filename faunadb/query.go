package faunadb

// Event's action types. https://fauna.com/documentation/queries#values-events
const (
	ActionCreate = "create"
	ActionDelete = "delete"
)

// Time unit. Usually used as a parameter for Epoch function. https://fauna.com/documentation/queries#time_functions-epoch_num_unit_unit
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

func Events(events interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["events"] = wrap(events)
	}
}

func TS(timestamp interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["ts"] = wrap(timestamp)
	}
}

func After(ref interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["after"] = wrap(ref)
	}
}

func Before(ref interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["before"] = wrap(ref)
	}
}

func Size(size interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["size"] = wrap(size)
	}
}

func Separator(sep interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["separator"] = wrap(sep)
	}
}

func Default(value interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["default"] = wrap(value)
	}
}

func Sources(sources interface{}) OptionalParameter {
	return func(fn unescapedObj) {
		fn["sources"] = wrap(sources)
	}
}

// Values

func Ref(id string) Expr                     { return RefV{id} }
func RefClass(classRef, id interface{}) Expr { return fn2("ref", classRef, "id", id) }
func Null() Expr                             { return NullV{} }

// Basic forms

func Do(exprs ...interface{}) Expr          { return fn1("do", varargs(exprs...)) }
func If(cond, then, elze interface{}) Expr  { return fn3("if", cond, "then", then, "else", elze) }
func Lambda(varName, expr interface{}) Expr { return fn2("lambda", varName, "expr", expr) }
func Let(bindings Obj, in interface{}) Expr { return fn2("let", unescapedBindings(bindings), "in", in) }
func Var(name string) Expr                  { return fn1("var", name) }

// Collections

func Map(coll, lambda interface{}) Expr     { return fn2("map", lambda, "collection", coll) }
func Foreach(coll, lambda interface{}) Expr { return fn2("foreach", lambda, "collection", coll) }
func Filter(coll, lambda interface{}) Expr  { return fn2("filter", lambda, "collection", coll) }
func Take(num, coll interface{}) Expr       { return fn2("take", num, "collection", coll) }
func Drop(num, coll interface{}) Expr       { return fn2("drop", num, "collection", coll) }
func Prepend(elems, coll interface{}) Expr  { return fn2("prepend", elems, "collection", coll) }
func Append(elems, coll interface{}) Expr   { return fn2("append", elems, "collection", coll) }

// Read

func Get(ref interface{}, options ...OptionalParameter) Expr {
	return fn1("get", ref, options...)
}

func Exists(ref interface{}, options ...OptionalParameter) Expr {
	return fn1("exists", ref, options...)
}

func Paginate(set interface{}, options ...OptionalParameter) Expr {
	return fn1("paginate", set, options...)
}

// Write

func Create(ref, params interface{}) Expr    { return fn2("create", ref, "params", params) }
func CreateClass(params interface{}) Expr    { return fn1("create_class", params) }
func CreateDatabase(params interface{}) Expr { return fn1("create_database", params) }
func CreateIndex(params interface{}) Expr    { return fn1("create_index", params) }
func CreateKey(params interface{}) Expr      { return fn1("create_key", params) }
func Update(ref, params interface{}) Expr    { return fn2("update", ref, "params", params) }
func Replace(ref, params interface{}) Expr   { return fn2("replace", ref, "params", params) }
func Delete(ref interface{}) Expr            { return fn1("delete", ref) }

func Insert(ref, ts, action, params interface{}) Expr {
	return fn4("insert", ref, "ts", ts, "action", action, "params", params)
}

func Remove(ref, ts, action interface{}) Expr {
	return fn3("remove", ref, "ts", ts, "action", action)
}

// String

func Concat(terms interface{}, options ...OptionalParameter) Expr {
	return fn1("concat", terms, options...)
}

func Casefold(str interface{}) Expr {
	return fn1("casefold", str)
}

// Time and Date

func Time(str interface{}) Expr        { return fn1("time", str) }
func Date(str interface{}) Expr        { return fn1("date", str) }
func Epoch(num, unit interface{}) Expr { return fn2("epoch", num, "unit", unit) }

// Set

func Match(ref interface{}) Expr            { return fn1("match", ref) }
func MatchTerm(ref, terms interface{}) Expr { return fn2("match", ref, "terms", terms) }
func Union(sets ...interface{}) Expr        { return fn1("union", varargs(sets...)) }
func Intersection(sets ...interface{}) Expr { return fn1("intersection", varargs(sets...)) }
func Difference(sets ...interface{}) Expr   { return fn1("difference", varargs(sets...)) }
func Distinct(set interface{}) Expr         { return fn1("distinct", set) }
func Join(source, target interface{}) Expr  { return fn2("join", source, "with", target) }

// Authentication

func Login(ref, params interface{}) Expr      { return fn2("login", ref, "params", params) }
func Logout(invalidateAll interface{}) Expr   { return fn1("logout", invalidateAll) }
func Identify(ref, password interface{}) Expr { return fn2("identify", ref, "password", password) }

// Miscellaneous

func NextID() Expr                          { return fn1("next_id", NullV{}) }
func Database(name interface{}) Expr        { return fn1("database", name) }
func Index(name interface{}) Expr           { return fn1("index", name) }
func Class(name interface{}) Expr           { return fn1("class", name) }
func Equals(args ...interface{}) Expr       { return fn1("equals", varargs(args...)) }
func Contains(path, value interface{}) Expr { return fn2("contains", path, "in", value) }
func Add(args ...interface{}) Expr          { return fn1("add", varargs(args...)) }
func Multiply(args ...interface{}) Expr     { return fn1("multiply", varargs(args...)) }
func Subtract(args ...interface{}) Expr     { return fn1("subtract", varargs(args...)) }
func Divide(args ...interface{}) Expr       { return fn1("divide", varargs(args...)) }
func Modulo(args ...interface{}) Expr       { return fn1("modulo", varargs(args...)) }
func LT(args ...interface{}) Expr           { return fn1("lt", varargs(args...)) }
func LTE(args ...interface{}) Expr          { return fn1("lte", varargs(args...)) }
func GT(args ...interface{}) Expr           { return fn1("gt", varargs(args...)) }
func GTE(args ...interface{}) Expr          { return fn1("gte", varargs(args...)) }
func And(args ...interface{}) Expr          { return fn1("and", varargs(args...)) }
func Or(args ...interface{}) Expr           { return fn1("or", varargs(args...)) }
func Not(boolean interface{}) Expr          { return fn1("not", boolean) }

func Select(path, value interface{}, options ...OptionalParameter) Expr {
	return fn2("select", path, "from", value, options...)
}
