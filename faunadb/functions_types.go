package faunadb

// ToString attempts to convert an expression to a string literal.
//
// Parameters:
//   value Object - The expression to convert.
//
// Returns:
//   string - A string literal.
func ToString(value interface{}) Expr {
	return toStringFn{ToString: wrap(value)}
}

type toStringFn struct {
	fnApply
	ToString Expr `json:"to_string"`
}

func (fn toStringFn) String() string { return printFn(fn) }

// ToNumber attempts to convert an expression to a numeric literal -
// either an int64 or float64.
//
// Parameters:
//   value Object - The expression to convert.
//
// Returns:
//   number - A numeric literal.
func ToNumber(value interface{}) Expr {
	return toNumberFn{ToNumber: wrap(value)}
}

type toNumberFn struct {
	fnApply
	ToNumber Expr `json:"to_number"`
}

func (fn toNumberFn) String() string { return printFn(fn) }

// ToDouble attempts to convert an expression to a double
//
// Parameters:
//   value Object - The expression to convert.
//
// Returns:
//   number - A double literal.
func ToDouble(value interface{}) Expr {
	return toDoubleFn{ToDouble: wrap(value)}
}

type toDoubleFn struct {
	fnApply
	ToDouble Expr `json:"to_double"`
}

func (fn toDoubleFn) String() string { return printFn(fn) }

// ToInteger attempts to convert an expression to an integer literal
//
// Parameters:
//   value Object - The expression to convert.
//
// Returns:
//   number - An integer literal.
func ToInteger(value interface{}) Expr {
	return toIntegerFn{ToInteger: wrap(value)}
}

type toIntegerFn struct {
	fnApply
	ToInteger Expr `json:"to_integer"`
}

func (fn toIntegerFn) String() string { return printFn(fn) }

// ToObject attempts to convert an expression to an object
//
// Parameters:
//   value Object - The expression to convert.
//
// Returns:
//   object - An object.
func ToObject(value interface{}) Expr {
	return toObjectFn{ToObject: wrap(value)}
}

type toObjectFn struct {
	fnApply
	ToObject Expr `json:"to_object"`
}

func (fn toObjectFn) String() string { return printFn(fn) }

// ToArray attempts to convert an expression to an array
//
// Parameters:
//   value Object - The expression to convert.
//
// Returns:
//   array - An array.
func ToArray(value interface{}) Expr {
	return toArrayFn{ToArray: wrap(value)}
}

type toArrayFn struct {
	fnApply
	ToArray Expr `json:"to_array"`
}

func (fn toArrayFn) String() string { return printFn(fn) }

// ToTime attempts to convert an expression to a time literal.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - A time literal.
func ToTime(value interface{}) Expr {
	return toTimeFn{ToTime: wrap(value)}
}

type toTimeFn struct {
	fnApply
	ToTime Expr `json:"to_time"`
}

func (fn toTimeFn) String() string { return printFn(fn) }

// ToDate attempts to convert an expression to a date literal.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   date - A date literal.
func ToDate(value interface{}) Expr {
	return toDateFn{ToDate: wrap(value)}
}

type toDateFn struct {
	fnApply
	ToDate Expr `json:"to_date"`
}

func (fn toDateFn) String() string { return printFn(fn) }

// IsNumber checks if the expression is a number
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool      -  returns true if the expression is a number
func IsNumber(expr interface{}) Expr {
	return isNumberFn{IsNumber: wrap(expr)}
}

type isNumberFn struct {
	fnApply
	IsNumber Expr `json:"is_number"`
}

func (fn isNumberFn) String() string { return printFn(fn) }

// IsDouble checks if the expression is a double
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a double
func IsDouble(expr interface{}) Expr {
	return isDoubleFn{IsDouble: wrap(expr)}
}

type isDoubleFn struct {
	fnApply
	IsDouble Expr `json:"is_double"`
}

func (fn isDoubleFn) String() string { return printFn(fn) }

// IsInteger checks if the expression is an integer
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is an integer
func IsInteger(expr interface{}) Expr {
	return isIntegerFn{IsInteger: wrap(expr)}
}

type isIntegerFn struct {
	fnApply
	IsInteger Expr `json:"is_integer"`
}

func (fn isIntegerFn) String() string { return printFn(fn) }

// IsBoolean checks if the expression is a boolean
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a boolean
func IsBoolean(expr interface{}) Expr {
	return isBooleanFn{IsBoolean: wrap(expr)}
}

type isBooleanFn struct {
	fnApply
	IsBoolean Expr `json:"is_boolean"`
}

func (fn isBooleanFn) String() string { return printFn(fn) }

// IsNull checks if the expression is null
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is null
func IsNull(expr interface{}) Expr {
	return isNullFn{IsNull: wrap(expr)}
}

type isNullFn struct {
	fnApply
	IsNull Expr `json:"is_null"`
}

func (fn isNullFn) String() string { return printFn(fn) }

// IsBytes checks if the expression are bytes
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression are bytes
func IsBytes(expr interface{}) Expr {
	return isBytesFn{IsBytes: wrap(expr)}
}

type isBytesFn struct {
	fnApply
	IsBytes Expr `json:"is_bytes"`
}

func (fn isBytesFn) String() string { return printFn(fn) }

// IsTimestamp checks if the expression is a timestamp
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a timestamp
func IsTimestamp(expr interface{}) Expr {
	return isTimestampFn{IsTimestamp: wrap(expr)}
}

type isTimestampFn struct {
	fnApply
	IsTimestamp Expr `json:"is_timestamp"`
}

func (fn isTimestampFn) String() string { return printFn(fn) }

// IsDate checks if the expression is a date
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a date
func IsDate(expr interface{}) Expr {
	return isDateFn{IsDate: wrap(expr)}
}

type isDateFn struct {
	fnApply
	IsDate Expr `json:"is_date"`
}

func (fn isDateFn) String() string { return printFn(fn) }

// IsString checks if the expression is a string
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a string
func IsString(expr interface{}) Expr {
	return isStringFn{IsString: wrap(expr)}
}

type isStringFn struct {
	fnApply
	IsString Expr `json:"is_string"`
}

func (fn isStringFn) String() string { return printFn(fn) }

// IsArray checks if the expression is an array
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is an array
func IsArray(expr interface{}) Expr {
	return isArrayFn{IsArray: wrap(expr)}
}

type isArrayFn struct {
	fnApply
	IsArray Expr `json:"is_array"`
}

func (fn isArrayFn) String() string { return printFn(fn) }

// IsObject checks if the expression is an object
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is an object
func IsObject(expr interface{}) Expr {
	return isObjectFn{IsObject: wrap(expr)}
}

type isObjectFn struct {
	fnApply
	IsObject Expr `json:"is_object"`
}

func (fn isObjectFn) String() string { return printFn(fn) }

// IsRef checks if the expression is a ref
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a ref
func IsRef(expr interface{}) Expr {
	return isRefFn{IsRef: wrap(expr)}
}

type isRefFn struct {
	fnApply
	IsRef Expr `json:"is_ref"`
}

func (fn isRefFn) String() string { return printFn(fn) }

// IsSet checks if the expression is a set
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a set
func IsSet(expr interface{}) Expr {
	return isSetFn{IsSet: wrap(expr)}
}

type isSetFn struct {
	fnApply
	IsSet Expr `json:"is_set"`
}

func (fn isSetFn) String() string { return printFn(fn) }

// IsDoc checks if the expression is a document
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a document
func IsDoc(expr interface{}) Expr {
	return isDocFn{IsDoc: wrap(expr)}
}

type isDocFn struct {
	fnApply
	IsDoc Expr `json:"is_doc"`
}

func (fn isDocFn) String() string { return printFn(fn) }

// IsLambda checks if the expression is a Lambda
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a Lambda
func IsLambda(expr interface{}) Expr {
	return isLambdaFn{IsLambda: wrap(expr)}
}

type isLambdaFn struct {
	fnApply
	IsLambda Expr `json:"is_lambda"`
}

func (fn isLambdaFn) String() string { return printFn(fn) }

// IsCollection checks if the expression is a collection
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a collection
func IsCollection(expr interface{}) Expr {
	return isCollectionFn{IsCollection: wrap(expr)}
}

type isCollectionFn struct {
	fnApply
	IsCollection Expr `json:"is_collection"`
}

func (fn isCollectionFn) String() string { return printFn(fn) }

// IsDatabase checks if the expression is a database
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a database
func IsDatabase(expr interface{}) Expr {
	return isDatabaseFn{IsDatabase: wrap(expr)}
}

type isDatabaseFn struct {
	fnApply
	IsDatabase Expr `json:"is_database"`
}

func (fn isDatabaseFn) String() string { return printFn(fn) }

// IsIndex checks if the expression is an index
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is an index
func IsIndex(expr interface{}) Expr {
	return isIndexFn{IsIndex: wrap(expr)}
}

type isIndexFn struct {
	fnApply
	IsIndex Expr `json:"is_index"`
}

func (fn isIndexFn) String() string { return printFn(fn) }

// IsFunction checks if the expression is a function
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a function
func IsFunction(expr interface{}) Expr {
	return isFunctionFn{IsFunction: wrap(expr)}
}

type isFunctionFn struct {
	fnApply
	IsFunction Expr `json:"is_function"`
}

func (fn isFunctionFn) String() string { return printFn(fn) }

// IsKey checks if the expression is a key
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a key
func IsKey(expr interface{}) Expr {
	return isKeyFn{IsKey: wrap(expr)}
}

type isKeyFn struct {
	fnApply
	IsKey Expr `json:"is_key"`
}

func (fn isKeyFn) String() string { return printFn(fn) }

// IsToken checks if the expression is a token
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a token
func IsToken(expr interface{}) Expr {
	return isTokenFn{IsToken: wrap(expr)}
}

type isTokenFn struct {
	fnApply
	IsToken Expr `json:"is_token"`
}

func (fn isTokenFn) String() string { return printFn(fn) }

// IsCredentials checks if the expression is a credentials
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a credential
func IsCredentials(expr interface{}) Expr {
	return isCredentialsFn{IsCredentials: wrap(expr)}
}

type isCredentialsFn struct {
	fnApply
	IsCredentials Expr `json:"is_credentials"`
}

func (fn isCredentialsFn) String() string { return printFn(fn) }

// IsRole checks if the expression is a role
//
// Parameters:
//  expr Expr - The expression to check.
//
// Returns:
//  bool         -  returns true if the expression is a role
func IsRole(expr interface{}) Expr {
	return isRoleFn{IsRole: wrap(expr)}
}

type isRoleFn struct {
	fnApply
	IsRole Expr `json:"is_role"`
}

func (fn isRoleFn) String() string { return printFn(fn) }
