package faunadb

// Basic forms

// Abort aborts the execution of the query
//
// Parameters:
//  msg string - An error message.
//
// Returns:
//  Error
//
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func Abort(msg interface{}) Expr { return abortFn{Abort: wrap(msg)} }

type abortFn struct {
	fnApply
	Abort Expr `json:"abort"`
}

// Do sequentially evaluates its arguments, and returns the last expression.
// If no expressions are provided, do returns an error.
//
// Parameters:
//  exprs []Expr - A variable number of expressions to be evaluated.
//
// Returns:
//  Value - The result of the last expression in the list.
//
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func Do(exprs ...interface{}) Expr { return doFn{Do: wrap(varargs(exprs))} }

type doFn struct {
	fnApply
	Do Expr `json:"do" faunarepr:"varargs"`
}

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
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func If(cond, then, elze interface{}) Expr {
	return ifFn{
		If:   wrap(cond),
		Then: wrap(then),
		Else: wrap(elze),
	}
}

type ifFn struct {
	fnApply
	If   Expr `json:"if"`
	Then Expr `json:"then"`
	Else Expr `json:"else"`
}

// Lambda creates an anonymous function. Mostly used with Collection functions.
//
// Parameters:
//  varName string|[]string - A string or an array of strings of arguments name to be bound in the body of the lambda.
//  expr Expr - An expression used as the body of the lambda.
//
// Returns:
//  Value - The result of the body expression.
//
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func Lambda(varName, expr interface{}) Expr {
	return lambdaFn{
		Lambda:     wrap(varName),
		Expression: wrap(expr),
	}
}

type lambdaFn struct {
	fnApply
	Lambda     Expr `json:"lambda"`
	Expression Expr `json:"expr"`
}

// At execute an expression at a given timestamp.
//
// Parameters:
//  timestamp time - The timestamp in which the expression will be evaluated.
//  expr Expr - An expression to be evaluated.
//
// Returns:
//  Value - The result of the given expression.
//
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func At(timestamp, expr interface{}) Expr {
	return atFn{
		At:         wrap(timestamp),
		Expression: wrap(expr),
	}
}

type atFn struct {
	fnApply
	At         Expr `json:"at"`
	Expression Expr `json:"expr"`
}

// LetBuilder builds Let expressions
type LetBuilder struct {
	bindings unescapedArr
}

type letFn struct {
	fnApply
	Let Expr `json:"let"`
	In  Expr `json:"in"`
}

// Bind binds a variable name to a value and returns a LetBuilder
func (lb *LetBuilder) Bind(key string, in interface{}) *LetBuilder {
	binding := make(unescapedObj, 1)
	binding[key] = wrap(in)
	lb.bindings = append(lb.bindings, binding)
	return lb
}

// In sets the expression to be evaluated and returns the prepared Let.
func (lb *LetBuilder) In(in Expr) Expr {
	return letFn{
		Let: wrap(lb.bindings),
		In:  in,
	}
}

// Let binds values to one or more variables.
//
// Returns:
//  *LetBuilder - Returns a LetBuilder.
//
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func Let() *LetBuilder { return &LetBuilder{nil} }

// Var refers to a value of a variable on the current lexical scope.
//
// Parameters:
//  name string - The variable name.
//
// Returns:
//  Value - The variable value.
//
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func Var(name string) Expr { return varFn{Var: wrap(name)} }

type varFn struct {
	fnApply
	Var Expr `json:"var"`
}

// Call invokes the specified function passing in a variable number of arguments
//
// Parameters:
//  ref Ref - The reference to the user defined functions to call.
//  args []Value - A series of values to pass as arguments to the user defined function.
//
// Returns:
//  Value - The return value of the user defined function.
//
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func Call(ref interface{}, args ...interface{}) Expr {
	return callFn{Call: wrap(ref), Params: wrap(varargs(args...))}
}

type callFn struct {
	fnApply
	Call   Expr `json:"call"`
	Params Expr `json:"arguments"`
}

// Query creates an instance of the @query type with the specified lambda
//
// Parameters:
//  lambda Lambda - A lambda representation. See Lambda() function.
//
// Returns:
//  Query - The lambda wrapped in a @query type.
//
// See: https://app.fauna.com/documentation/reference/queryapi#basic-forms
func Query(lambda interface{}) Expr { return queryFn{Query: wrap(lambda)} }

type queryFn struct {
	fnApply
	Query Expr `json:"query"`
}

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
// See: https://app.fauna.com/documentation/reference/queryapi#read-functions
func Select(path, value interface{}, options ...OptionalParameter) Expr {
	fn := selectFn{Select: wrap(path), From: wrap(value)}
	return applyOptionals(fn, options)
}

type selectFn struct {
	fnApply
	Select  Expr `json:"select"`
	From    Expr `json:"from"`
	Default Expr `json:"default,omitempty" faunarepr:"optfn"`
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
// See: https://app.fauna.com/documentation/reference/queryapi#read-functions
func SelectAll(path, value interface{}) Expr {
	return selectAllFn{SelectAll: wrap(path), From: wrap(value)}
}

type selectAllFn struct {
	fnApply
	SelectAll Expr `json:"select_all"`
	From      Expr `json:"from"`
	Default   Expr `json:"default,omitempty" faunarepr:"optfn"`
}
