package faunadb

// Abs computes the absolute value of a number.
//
// Parameters:
//  value number - The number to take the absolute value of
//
// Returns:
//  number - The abosulte value of a number
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Abs(value interface{}) Expr { return absFn{Abs: wrap(value)} }

type absFn struct {
	fnApply
	Abs Expr `json:"abs"`
}

// Acos computes the arccosine of a number.
//
// Parameters:
//  value number - The number to take the arccosine of
//
// Returns:
//  number - The arccosine of a number
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Acos(value interface{}) Expr { return acosFn{Acos: wrap(value)} }

type acosFn struct {
	fnApply
	Acos Expr `json:"acos"`
}

// Asin computes the arcsine of a number.
//
// Parameters:
//  value number - The number to take the arcsine of
//
// Returns:
//  number - The arcsine of a number
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Asin(value interface{}) Expr { return asinFn{Asin: wrap(value)} }

type asinFn struct {
	fnApply
	Asin Expr `json:"asin"`
}

// Atan computes the arctan of a number.
//
// Parameters:
//  value number - The number to take the arctan of
//
// Returns:
//  number - The arctan of a number
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Atan(value interface{}) Expr { return atanFn{Atan: wrap(value)} }

type atanFn struct {
	fnApply
	Atan Expr `json:"atan"`
}

// Add computes the sum of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to sum together.
//
// Returns:
//  number - The sum of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Add(args ...interface{}) Expr { return addFn{Add: wrap(varargs(args...))} }

type addFn struct {
	fnApply
	Add Expr `json:"add" faunarepr:"varargs"`
}

// BitAnd computes the and of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to and together.
//
// Returns:
//  number - The and of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func BitAnd(args ...interface{}) Expr { return bitAndFn{BitAnd: wrap(varargs(args...))} }

type bitAndFn struct {
	fnApply
	BitAnd Expr `json:"bitand" faunarepr:"varargs"`
}

// BitNot computes the 2's complement of a number
//
// Parameters:
//  value number - A numbers to not
//
// Returns:
//  number - The not of an element
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func BitNot(value interface{}) Expr { return bitNotFn{BitNot: wrap(value)} }

type bitNotFn struct {
	fnApply
	BitNot Expr `json:"bitnot"`
}

// BitOr computes the OR of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to OR together.
//
// Returns:
//  number - The OR of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func BitOr(args ...interface{}) Expr { return bitOrFn{BitOr: wrap(varargs(args...))} }

type bitOrFn struct {
	fnApply
	BitOr Expr `json:"bitor" faunarepr:"varargs"`
}

// BitXor computes the XOR of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to XOR together.
//
// Returns:
//  number - The XOR of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func BitXor(args ...interface{}) Expr { return bitXorFn{BitXor: wrap(varargs(args...))} }

type bitXorFn struct {
	fnApply
	BitXor Expr `json:"bitxor" faunarepr:"varargs"`
}

// Ceil computes the largest integer greater than or equal to
//
// Parameters:
//  value number - A numbers to compute the ceil of
//
// Returns:
//  number - The ceil of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Ceil(value interface{}) Expr { return ceilFn{Ceil: wrap(value)} }

type ceilFn struct {
	fnApply
	Ceil Expr `json:"ceil"`
}

// Cos computes the Cosine of a number
//
// Parameters:
//  value number - A number to compute the cosine of
//
// Returns:
//  number - The cosine of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Cos(value interface{}) Expr { return cosFn{Cos: wrap(value)} }

type cosFn struct {
	fnApply
	Cos Expr `json:"cos"`
}

// Cosh computes the Hyperbolic Cosine of a number
//
// Parameters:
//  value number - A number to compute the Hyperbolic cosine of
//
// Returns:
//  number - The Hyperbolic cosine of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Cosh(value interface{}) Expr { return coshFn{Cosh: wrap(value)} }

type coshFn struct {
	fnApply
	Cosh Expr `json:"cosh"`
}

// Degrees computes the degress of a number
//
// Parameters:
//  value number - A number to compute the degress of
//
// Returns:
//  number - The degrees of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Degrees(value interface{}) Expr { return degreesFn{Degrees: wrap(value)} }

type degreesFn struct {
	fnApply
	Degrees Expr `json:"degrees"`
}

// Divide computes the quotient of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to compute the quotient of.
//
// Returns:
//  number - The quotient of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Divide(args ...interface{}) Expr { return divideFn{Divide: wrap(varargs(args...))} }

type divideFn struct {
	fnApply
	Divide Expr `json:"divide" faunarepr:"varargs"`
}

// Exp computes the Exp of a number
//
// Parameters:
//  value number - A number to compute the exp of
//
// Returns:
//  number - The exp of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Exp(value interface{}) Expr { return expFn{Exp: wrap(value)} }

type expFn struct {
	fnApply
	Exp Expr `json:"exp"`
}

// Floor computes the Floor of a number
//
// Parameters:
//  value number - A number to compute the Floor of
//
// Returns:
//  number - The Floor of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Floor(value interface{}) Expr { return floorFn{Floor: wrap(value)} }

type floorFn struct {
	fnApply
	Floor Expr `json:"floor"`
}

// Hypot computes the Hypotenuse of two numbers
//
// Parameters:
//  a number - A side of a right triangle
//  b number - A side of a right triangle
//
// Returns:
//  number - The hypotenuse of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Hypot(a, b interface{}) Expr { return hypotFn{Hypot: wrap(a), B: wrap(b)} }

type hypotFn struct {
	fnApply
	Hypot Expr `json:"hypot"`
	B     Expr `json:"b"`
}

// Ln computes the natural log of a number
//
// Parameters:
//  value number - A number to compute the natural log of
//
// Returns:
//  number - The ln of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Ln(value interface{}) Expr { return lnFn{Ln: wrap(value)} }

type lnFn struct {
	fnApply
	Ln Expr `json:"ln"`
}

// Log computes the Log of a number
//
// Parameters:
//  value number - A number to compute the Log of
//
// Returns:
//  number - The Log of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Log(value interface{}) Expr { return logFn{Log: wrap(value)} }

type logFn struct {
	fnApply
	Log Expr `json:"log"`
}

// Max computes the max of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to find the max of.
//
// Returns:
//  number - The max of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Max(args ...interface{}) Expr { return maxFn{Max: wrap(varargs(args...))} }

type maxFn struct {
	fnApply
	Max Expr `json:"max" faunarepr:"varargs"`
}

// Min computes the Min of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to find the min of.
//
// Returns:
//  number - The min of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Min(args ...interface{}) Expr { return minFn{Min: wrap(varargs(args...))} }

type minFn struct {
	fnApply
	Min Expr `json:"min" faunarepr:"varargs"`
}

// Modulo computes the reminder after the division of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to compute the quotient of. The remainder will be returned.
//
// Returns:
//  number - The remainder of the quotient of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Modulo(args ...interface{}) Expr { return moduloFn{Modulo: wrap(varargs(args...))} }

type moduloFn struct {
	fnApply
	Modulo Expr `json:"modulo" faunarepr:"varargs"`
}

// Multiply computes the product of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to multiply together.
//
// Returns:
//  number - The multiplication of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Multiply(args ...interface{}) Expr { return multiplyFn{Multiply: wrap(varargs(args...))} }

type multiplyFn struct {
	fnApply
	Multiply Expr `json:"multiply" faunarepr:"varargs"`
}

// Pow computes the Power of a number
//
// Parameters:
//  base number - A number which is the base
//  exp number  - A number which is the exponent
//
// Returns:
//  number - The Pow of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Pow(base, exp interface{}) Expr { return powFn{Pow: wrap(base), Exp: wrap(exp)} }

type powFn struct {
	fnApply
	Exp Expr `json:"exp"`
	Pow Expr `json:"pow"`
}

// Radians computes the Radians of a number
//
// Parameters:
//  value number - A number which is convert to radians
//
// Returns:
//  number - The Radians of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Radians(value interface{}) Expr { return radiansFn{Radians: wrap(value)} }

type radiansFn struct {
	fnApply
	Radians Expr `json:"radians"`
}

// Round a number at the given percission
//
// Parameters:
//  value number - The number to truncate
//  precision number - precision where to truncate, defaults is 2
//
// Returns:
//  number - The Rounded value.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Round(value interface{}, options ...OptionalParameter) Expr {
	fn := roundFn{Round: wrap(value)}
	return applyOptionals(fn, options)
}

type roundFn struct {
	fnApply
	Round     Expr `json:"round"`
	Precision Expr `json:"precision,omitempty" faunarepr:"optfn"`
}

// Sign computes the Sign of a number
//
// Parameters:
//  value number - A number to compute the Sign of
//
// Returns:
//  number - The Sign of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Sign(value interface{}) Expr { return signFn{Sign: wrap(value)} }

type signFn struct {
	fnApply
	Sign Expr `json:"sign"`
}

// Sin computes the Sine of a number
//
// Parameters:
//  value number - A number to compute the Sine of
//
// Returns:
//  number - The Sine of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Sin(value interface{}) Expr { return sinFn{Sin: wrap(value)} }

type sinFn struct {
	fnApply
	Sin Expr `json:"sin"`
}

// Sinh computes the Hyperbolic Sine of a number
//
// Parameters:
//  value number - A number to compute the Hyperbolic Sine of
//
// Returns:
//  number - The Hyperbolic Sine of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Sinh(value interface{}) Expr { return sinhFn{Sinh: wrap(value)} }

type sinhFn struct {
	fnApply
	Sinh Expr `json:"sinh"`
}

// Sqrt computes the square root of a number
//
// Parameters:
//  value number - A number to compute the square root of
//
// Returns:
//  number - The square root of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Sqrt(value interface{}) Expr { return sqrtFn{Sqrt: wrap(value)} }

type sqrtFn struct {
	fnApply
	Sqrt Expr `json:"sqrt"`
}

// Subtract computes the difference of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to compute the difference of.
//
// Returns:
//  number - The difference of all elements.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Subtract(args ...interface{}) Expr { return subtractFn{Subtract: wrap(varargs(args...))} }

type subtractFn struct {
	fnApply
	Subtract Expr `json:"subtract" faunarepr:"varargs"`
}

// Tan computes the Tangent of a number
//
// Parameters:
//  value number - A number to compute the Tangent of
//
// Returns:
//  number - The Tangent of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Tan(value interface{}) Expr { return tanFn{Tan: wrap(value)} }

type tanFn struct {
	fnApply
	Tan Expr `json:"tan"`
}

// Tanh computes the Hyperbolic Tangent of a number
//
// Parameters:
//  value number - A number to compute the Hyperbolic Tangent of
//
// Returns:
//  number - The Hyperbolic Tangent of value
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Tanh(value interface{}) Expr { return tanhFn{Tanh: wrap(value)} }

type tanhFn struct {
	fnApply
	Tanh Expr `json:"tanh"`
}

// Trunc truncates a number at the given percission
//
// Parameters:
//  value number - The number to truncate
//  precision number - precision where to truncate, defaults is 2
//
// Returns:
//  number - The truncated value.
//
// See: https://app.fauna.com/documentation/reference/queryapi#mathematical-functions
func Trunc(value interface{}, options ...OptionalParameter) Expr {
	fn := truncFn{Trunc: wrap(value)}
	return applyOptionals(fn, options)
}

type truncFn struct {
	fnApply
	Trunc     Expr `json:"trunc"`
	Precision Expr `json:"precision,omitempty" faunarepr:"optfn"`
}
