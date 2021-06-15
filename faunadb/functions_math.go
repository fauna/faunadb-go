package faunadb

// Abs computes the absolute value of a number.
//
// Parameters:
//  value number - The number to take the absolute value of.
//
// Returns:
//  number - The absolute value of a number.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/abs?lang=go
func Abs(value interface{}) Expr { return absFn{Abs: wrap(value)} }

type absFn struct {
	fnApply
	Abs Expr `json:"abs"`
}

// Acos computes the arccosine of a number.
//
// Parameters:
//  value number - The number to take the arccosine of.
//
// Returns:
//  number - The arccosine of a number.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/acos?lang=go
func Acos(value interface{}) Expr { return acosFn{Acos: wrap(value)} }

type acosFn struct {
	fnApply
	Acos Expr `json:"acos"`
}

// Asin computes the arcsine of a number.
//
// Parameters:
//  value number - The number to take the arcsine of.
//
// Returns:
//  number - The arcsine of a number.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/asin?lang=go
func Asin(value interface{}) Expr { return asinFn{Asin: wrap(value)} }

type asinFn struct {
	fnApply
	Asin Expr `json:"asin"`
}

// Atan computes the arctan of a number.
//
// Parameters:
//  value number - The number to take the arctan of.
//
// Returns:
//  number - The arctan of a number.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/atan?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/add?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/bitand?lang=go
func BitAnd(args ...interface{}) Expr { return bitAndFn{BitAnd: wrap(varargs(args...))} }

type bitAndFn struct {
	fnApply
	BitAnd Expr `json:"bitand" faunarepr:"varargs"`
}

// BitNot computes the two's complement of a number.
//
// Parameters:
//  value number - A numbers to not.
//
// Returns:
//  number - The not of an element.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/bitnot?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/bitor?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/bitxor?lang=go
func BitXor(args ...interface{}) Expr { return bitXorFn{BitXor: wrap(varargs(args...))} }

type bitXorFn struct {
	fnApply
	BitXor Expr `json:"bitxor" faunarepr:"varargs"`
}

// Ceil computes the smallest integer greater than or equal to the
// provided value.
//
// Parameters:
//  value number - A number to compute the ceil of.
//
// Returns:
//  number - The ceil of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/ceil?lang=go
func Ceil(value interface{}) Expr { return ceilFn{Ceil: wrap(value)} }

type ceilFn struct {
	fnApply
	Ceil Expr `json:"ceil"`
}

// Cos computes the cosine of a number.
//
// Parameters:
//  value number - A number to compute the cosine of.
//
// Returns:
//  number - The cosine of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/cos?lang=go
func Cos(value interface{}) Expr { return cosFn{Cos: wrap(value)} }

type cosFn struct {
	fnApply
	Cos Expr `json:"cos"`
}

// Cosh computes the hyperbolic cosine of a number.
//
// Parameters:
//  value number - A number to compute the hyperbolic cosine of.
//
// Returns:
//  number - The hyperbolic cosine of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/cosh?lang=go
func Cosh(value interface{}) Expr { return coshFn{Cosh: wrap(value)} }

type coshFn struct {
	fnApply
	Cosh Expr `json:"cosh"`
}

// Degrees converts the provided radians to degrees.
//
// Parameters:
//  value number - A number in radians to compute the degrees of.
//
// Returns:
//  number - The degrees of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/degrees?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/divide?lang=go
func Divide(args ...interface{}) Expr { return divideFn{Divide: wrap(varargs(args...))} }

type divideFn struct {
	fnApply
	Divide Expr `json:"divide" faunarepr:"varargs"`
}

// Exp computes the value of e to the given exponent.
//
// Parameters:
//  value number - A number to compute the exponent of.
//
// Returns:
//  number - The exponent of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/exp?lang=go
func Exp(value interface{}) Expr { return expFn{Exp: wrap(value)} }

type expFn struct {
	fnApply
	Exp Expr `json:"exp"`
}

// Floor computes the largest integer that is smaller then, or equal to,
// the provided value.
//
// Parameters:
//  value number - A number to compute the floor of.
//
// Returns:
//  number - The floor of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/floor?lang=go
func Floor(value interface{}) Expr { return floorFn{Floor: wrap(value)} }

type floorFn struct {
	fnApply
	Floor Expr `json:"floor"`
}

// Hypot computes the hypotenuse of a right triangle whose other two
// sides are of length a and b.
//
// Parameters:
//  a number - A side of a right triangle.
//  b number - A side of a right triangle.
//
// Returns:
//  number - The hypotenuse of a right triangle whose other two sides
//           are of length a and b.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/hypot?lang=go
func Hypot(a, b interface{}) Expr { return hypotFn{Hypot: wrap(a), B: wrap(b)} }

type hypotFn struct {
	fnApply
	Hypot Expr `json:"hypot"`
	B     Expr `json:"b"`
}

// Ln computes the natural log of a number.
//
// Parameters:
//  value number - A number to compute the natural log of.
//
// Returns:
//  number - The ln of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/ln?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/log?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/max?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/min?lang=go
func Min(args ...interface{}) Expr { return minFn{Min: wrap(varargs(args...))} }

type minFn struct {
	fnApply
	Min Expr `json:"min" faunarepr:"varargs"`
}

// Modulo computes the reminder after the division of a list of numbers.
//
// Parameters:
//  args []number - A collection of numbers to compute the quotient of.
//                  The remainder is returned.
//
// Returns:
//  number - The remainder of the quotient of all elements.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/modulo?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/multiply?lang=go
func Multiply(args ...interface{}) Expr { return multiplyFn{Multiply: wrap(varargs(args...))} }

type multiplyFn struct {
	fnApply
	Multiply Expr `json:"multiply" faunarepr:"varargs"`
}

// Pow computes the value of base raised to the given exp.
//
// Parameters:
//  base number - A number which is the base.
//  exp number  - A number which is the exponent.
//
// Returns:
//  number - The power of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/pow?lang=go
func Pow(base, exp interface{}) Expr { return powFn{Pow: wrap(base), Exp: wrap(exp)} }

type powFn struct {
	fnApply
	Exp Expr `json:"exp"`
	Pow Expr `json:"pow"`
}

// Radians converts the provided degrees to radians.
//
// Parameters:
//  value number - A number which is converted to radians.
//
// Returns:
//  number - The radians of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/radians?lang=go
func Radians(value interface{}) Expr { return radiansFn{Radians: wrap(value)} }

type radiansFn struct {
	fnApply
	Radians Expr `json:"radians"`
}

// Round a number at the given precision.
//
// Parameters:
//  value number     - The number to truncate.
//  precision number - The decimal precision to round to, default is 2.
//
// Returns:
//  number - The rounded value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/round?lang=go
func Round(value interface{}, options ...OptionalParameter) Expr {
	fn := roundFn{Round: wrap(value)}
	return applyOptionals(fn, options)
}

type roundFn struct {
	fnApply
	Round     Expr `json:"round"`
	Precision Expr `json:"precision,omitempty" faunarepr:"optfn"`
}

// Sign computes the sign of a number, returning 1 when the value is
// positive, 0 when the value is zero, and -1 when the value is
// negative.
//
// Parameters:
//  value number - A number to compute the sign of.
//
// Returns:
//  number - The sign of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/sign?lang=go
func Sign(value interface{}) Expr { return signFn{Sign: wrap(value)} }

type signFn struct {
	fnApply
	Sign Expr `json:"sign"`
}

// Sin computes the sine of a number.
//
// Parameters:
//  value number - A number to compute the sine of.
//
// Returns:
//  number - The sine of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/sin?lang=go
func Sin(value interface{}) Expr { return sinFn{Sin: wrap(value)} }

type sinFn struct {
	fnApply
	Sin Expr `json:"sin"`
}

// Sinh computes the hyperbolic sine of a number.
//
// Parameters:
//  value number - A number to compute the hyperbolic sine of.
//
// Returns:
//  number - The hyperbolic sine of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/sinh?lang=go
func Sinh(value interface{}) Expr { return sinhFn{Sinh: wrap(value)} }

type sinhFn struct {
	fnApply
	Sinh Expr `json:"sinh"`
}

// Sqrt computes the square root of a number.
//
// Parameters:
//  value number - A number to compute the square root of.
//
// Returns:
//  number - The square root of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/sqrt?lang=go
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
// See: https://docs.fauna.com/fauna/current/api/fql/functions/subtract?lang=go
func Subtract(args ...interface{}) Expr { return subtractFn{Subtract: wrap(varargs(args...))} }

type subtractFn struct {
	fnApply
	Subtract Expr `json:"subtract" faunarepr:"varargs"`
}

// Tan computes the tangent of a number.
//
// Parameters:
//  value number - A number to compute the tangent of.
//
// Returns:
//  number - The tangent of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/tan?lang=go
func Tan(value interface{}) Expr { return tanFn{Tan: wrap(value)} }

type tanFn struct {
	fnApply
	Tan Expr `json:"tan"`
}

// Tanh computes the hyperbolic tangent of a number.
//
// Parameters:
//  value number - A number to compute the hyperbolic tangent of.
//
// Returns:
//  number - The hyperbolic tangent of value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/tanh?lang=go
func Tanh(value interface{}) Expr { return tanhFn{Tanh: wrap(value)} }

type tanhFn struct {
	fnApply
	Tanh Expr `json:"tanh"`
}

// Trunc truncates a number at the given precision.
//
// Parameters:
//  value number     - The number to truncate.
//  precision number - The decimal precision to truncate to, defaults is 2
//
// Returns:
//  number - The truncated value.
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/trunc?lang=go
func Trunc(value interface{}, options ...OptionalParameter) Expr {
	fn := truncFn{Trunc: wrap(value)}
	return applyOptionals(fn, options)
}

type truncFn struct {
	fnApply
	Trunc     Expr `json:"trunc"`
	Precision Expr `json:"precision,omitempty" faunarepr:"optfn"`
}
