package faunadb

type fnApply interface {
	Expr
}

// OptionalParameter describes optional parameters for query language functions
type OptionalParameter func(Expr) Expr

func applyOptionals(expr Expr, options []OptionalParameter) Expr {
	for _, option := range options {
		switch expr.(type) {
		case invalidExpr:
			return expr
		default:
			expr = option(expr)
		}
	}
	return expr
}

// Event's action types. Usually used as a parameter for Insert or Remove functions.
//
// See: https://app.fauna.com/documentation/reference/queryapi#simple-type-events
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionAdd    = "add"
	ActionRemove = "remove"
)

// Time unit. Usually used as a parameter for Time functions.
//
// See: https://app.fauna.com/documentation/reference/queryapi#epochnum-unit
const (
	TimeUnitDay         = "day"
	TimeUnitHalfDay     = "half day"
	TimeUnitHour        = "hour"
	TimeUnitMinute      = "minute"
	TimeUnitSecond      = "second"
	TimeUnitMillisecond = "millisecond"
	TimeUnitMicrosecond = "microsecond"
	TimeUnitNanosecond  = "nanosecond"
)

// Normalizers for Casefold
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
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

// Optional parameters

// EventsOpt is an boolean optional parameter that describes if the query should include historical events.
// For more information about events, check https://app.fauna.com/documentation/reference/queryapi#simple-type-events.
//
// Functions that accept this optional parameter are: Paginate.
//
// Deprecated: The Events function was renamed to EventsOpt to support the new history API.
// EventsOpt is provided here for backwards compatibility. Instead of using Paginate with the EventsOpt parameter,
// you should use the new Events function.
func EventsOpt(events interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case eventsParam:
			return e.setEvents(wrap(events))
		default:
			return e
		}
	}
}

type eventsParam interface {
	setEvents(length Expr) Expr
}

func (fn paginateFn) setEvents(e Expr) Expr {
	fn.Events = e
	return fn
}

// TS is a timestamp optional parameter that specifies in which timestamp a query should be executed.
//
// Functions that accept this optional parameter are: Get, Exists, and Paginate.
func TS(timestamp interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case tsParam:
			return e.setTS(wrap(timestamp))
		default:
			return e
		}
	}
}

type tsParam interface {
	setTS(length Expr) Expr
}

func (fn paginateFn) setTS(e Expr) Expr {
	fn.TS = e
	return fn
}

func (fn existsFn) setTS(e Expr) Expr {
	fn.TS = e
	return fn
}

func (fn getFn) setTS(e Expr) Expr {
	fn.TS = e
	return fn
}

func Cursor(ref interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case cursorParam:
			return e.setCursor(wrap(ref))
		default:
			return e
		}
	}
}

type cursorParam interface {
	setCursor(cursor Expr) Expr
}

func (fn paginateFn) setCursor(e Expr) Expr {
	fn.Cursor = e
	return fn
}

// After is an optional parameter used when cursoring that refers to the specified cursor's the next page, inclusive.
// For more information about pages, check https://app.fauna.com/documentation/reference/queryapi#simple-type-pages.
//
// Functions that accept this optional parameter are: Paginate.
func After(ref interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case afterParam:
			return e.setAfter(wrap(ref))
		default:
			return e
		}
	}
}

type afterParam interface {
	setAfter(after Expr) Expr
}

func (fn paginateFn) setAfter(e Expr) Expr {
	fn.After = e
	return fn
}

// Before is an optional parameter used when cursoring that refers to the specified cursor's previous page, exclusive.
// For more information about pages, check https://app.fauna.com/documentation/reference/queryapi#simple-type-pages.
//
// Functions that accept this optional parameter are: Paginate.
func Before(ref interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case beforeParam:
			return e.setBefore(wrap(ref))
		default:
			return e
		}
	}
}

type beforeParam interface {
	setBefore(before Expr) Expr
}

func (fn paginateFn) setBefore(e Expr) Expr {
	fn.Before = e
	return fn
}

// Number is a numeric optional parameter that specifies an optional number.
//
// Functions that accept this optional parameter are: Repeat.
func Number(num interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case numberParam:
			return e.setNumber(wrap(num))
		default:
			return e
		}
	}
}

type numberParam interface {
	setNumber(num Expr) Expr
}

func (fn repeatFn) setNumber(e Expr) Expr {
	fn.Number = e
	return fn
}

// Size is a numeric optional parameter that specifies the size of a pagination cursor.
//
// Functions that accept this optional parameter are: Paginate.
func Size(size interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case sizeParam:
			return e.setSize(wrap(size))
		default:
			return e
		}
	}
}

type sizeParam interface {
	setSize(size Expr) Expr
}

func (fn paginateFn) setSize(e Expr) Expr {
	fn.Size = e
	return fn
}

// NumResults is a numeric optional parameter that specifies the number of results returned.
//
// Functions that accept this optional parameter are: FindStrRegex.
func NumResults(num interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case numParam:
			return e.setNum(wrap(num))
		default:
			return e
		}
	}
}

type numParam interface {
	setNum(num Expr) Expr
}

func (fn findStrRegexFn) setNum(e Expr) Expr {
	fn.NumResults = e
	return fn
}

// Start is a numeric optional parameter that specifies the start of where to search.
//
// Functions that accept this optional parameter are: FindStr .
func Start(start interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case startParam:
			return e.setStart(wrap(start))
		default:
			return e
		}
	}
}

type startParam interface {
	setStart(start Expr) Expr
}

func (fn findStrFn) setStart(e Expr) Expr {
	fn.Start = e
	return fn
}

func (fn findStrRegexFn) setStart(e Expr) Expr {
	fn.Start = e
	return fn
}

// StrLength is a numeric optional parameter that specifies the amount to copy.
//
// Functions that accept this optional parameter are: FindStr and FindStrRegex.
func StrLength(length interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case strLengthParam:
			return e.setLength(wrap(length))
		default:
			return e
		}
	}
}

type strLengthParam interface {
	setLength(length Expr) Expr
}

func (fn subStringFn) setLength(e Expr) Expr {
	fn.Length = e
	return fn
}

// OnlyFirst is a boolean optional parameter that only replace the first string
//
// Functions that accept this optional parameter are: ReplaceStrRegex
func OnlyFirst() OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case firstParam:
			return e.setOnlyFirst()
		default:
			return e
		}
	}
}

type firstParam interface {
	setOnlyFirst() Expr
}

func (fn replaceStrRegexFn) setOnlyFirst() Expr {
	fn.First = BooleanV(true)
	return fn
}

// Sources is a boolean optional parameter that specifies if a pagination cursor should include
// the source sets along with each element.
//
// Functions that accept this optional parameter are: Paginate.
func Sources(sources interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case sourcesParam:
			return e.setSources(wrap(sources))
		default:
			return e
		}
	}
}

type sourcesParam interface {
	setSources(sources Expr) Expr
}

func (fn paginateFn) setSources(e Expr) Expr {
	fn.Sources = e
	return fn
}

// Default is an optional parameter that specifies the default value for a select operation when
// the desired value path is absent.
//
// Functions that accept this optional parameter are: Select.
func Default(value interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case defaultParam:
			return e.setDefault(wrap(value))
		default:
			return e
		}
	}
}

type defaultParam interface {
	setDefault(value Expr) Expr
}

func (fn selectFn) setDefault(e Expr) Expr {
	fn.Default = e
	return fn
}

func (fn selectAllFn) setDefault(e Expr) Expr {
	fn.Default = e
	return fn
}

// Separator is a string optional parameter that specifies the separator for a concat operation.
//
// Functions that accept this optional parameter are: Concat.
func Separator(sep interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case separatorParam:
			return e.setSeparator(wrap(sep))
		default:
			return e
		}
	}
}

type separatorParam interface {
	setSeparator(sep Expr) Expr
}

func (fn concatFn) setSeparator(e Expr) Expr {
	fn.Separator = e
	return fn
}

// Precision is an optional parameter that specifies the precision for a Trunc and Round operations.
//
// Functions that accept this optional parameter are: Round and Trunc.
func Precision(precision interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case precisionParam:
			return e.setPrecision(wrap(precision))
		default:
			return e
		}
	}
}

type precisionParam interface {
	setPrecision(precision Expr) Expr
}

func (fn roundFn) setPrecision(e Expr) Expr {
	fn.Precision = e
	return fn
}

func (fn truncFn) setPrecision(e Expr) Expr {
	fn.Precision = e
	return fn
}

// ConflictResolver is an optional parameter that specifies the lambda for resolving Merge conflicts
//
// Functions that accept this optional parameter are: Merge
func ConflictResolver(lambda interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case lambdaParam:
			return e.setLambda(wrap(lambda))
		default:
			return e
		}
	}
}

type lambdaParam interface {
	setLambda(lambda Expr) Expr
}

func (fn mergeFn) setLambda(e Expr) Expr {
	fn.Lambda = e
	return fn
}

// Normalizer is a string optional parameter that specifies the normalization function for casefold operation.
//
// Functions that accept this optional parameter are: Casefold.
func Normalizer(norm interface{}) OptionalParameter {
	return func(expr Expr) Expr {
		switch e := expr.(type) {
		case normalizerParam:
			return e.setNormalizer(wrap(norm))
		default:
			return e
		}
	}
}

type normalizerParam interface {
	setNormalizer(normalizer Expr) Expr
}

func (fn casefoldFn) setNormalizer(e Expr) Expr {
	fn.Normalizer = e
	return fn
}
