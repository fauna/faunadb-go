package faunadb

// Time and Date

// Time constructs a time from a ISO 8601 offset date/time string.
//
// Parameters:
//  str string - A string to convert to a time object.
//
// Returns:
//  time - A time object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#time-and-date
func Time(str interface{}) Expr { return timeFn{Time: wrap(str)} }

type timeFn struct {
	fnApply
	Time Expr `json:"time"`
}

// TimeAdd returns a new time or date with the offset in terms of the unit
// added.
//
// Parameters:
// base        -  the base time or data
// offset      -  the number of units
// unit        -  the unit type
//
// Returns:
// Expr
//
//See: https://docs.fauna.com/fauna/current/api/fql/functions/timeadd
func TimeAdd(base interface{}, offset interface{}, unit interface{}) Expr {
	return timeAddFn{TimeAdd: wrap(base), Offset: wrap(offset), Unit: wrap(unit)}
}

type timeAddFn struct {
	fnApply
	TimeAdd Expr `json:"time_add"`
	Offset  Expr `json:"offset"`
	Unit    Expr `json:"unit"`
}

// TimeSubtract returns a new time or date with the offset in terms of the unit
// subtracted.
//
// Parameters:
// base        -  the base time or data
// offset      -  the number of units
// unit        -  the unit type
//
// Returns:
// Expr
//
//See: https://docs.fauna.com/fauna/current/api/fql/functions/timesubtract
func TimeSubtract(base interface{}, offset interface{}, unit interface{}) Expr {
	return timeSubtractFn{TimeSubtract: wrap(base), Offset: wrap(offset), Unit: wrap(unit)}
}

type timeSubtractFn struct {
	fnApply
	TimeSubtract Expr `json:"time_subtract"`
	Offset       Expr `json:"offset"`
	Unit         Expr `json:"unit"`
}

// TimeDiff returns the number of intervals in terms of the unit between
// two times or dates. Both start and finish must be of the same
// type.
//
// Parameters:
//   start the starting time or date, inclusive
//   finish the ending time or date, exclusive
//   unit the unit type//
// Returns:
// Expr
//
//See: https://docs.fauna.com/fauna/current/api/fql/functions/timediff
func TimeDiff(start interface{}, finish interface{}, unit interface{}) Expr {
	return timeDiffFn{TimeDiff: wrap(start), Other: wrap(finish), Unit: wrap(unit)}
}

type timeDiffFn struct {
	fnApply
	TimeDiff Expr `json:"time_diff"`
	Other    Expr `json:"other"`
	Unit     Expr `json:"unit"`
}

// Date constructs a date from a ISO 8601 offset date/time string.
//
// Parameters:
//  str string - A string to convert to a date object.
//
// Returns:
//  date - A date object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#time-and-date
func Date(str interface{}) Expr { return dateFn{Date: wrap(str)} }

type dateFn struct {
	fnApply
	Date Expr `json:"date"`
}

// Epoch constructs a time relative to the epoch "1970-01-01T00:00:00Z".
//
// Parameters:
//  num int64 - The number of units from Epoch.
//  unit string - The unit of number. One of TimeUnitSecond, TimeUnitMillisecond, TimeUnitMicrosecond, TimeUnitNanosecond.
//
// Returns:
//  time - A time object.
//
// See: https://app.fauna.com/documentation/reference/queryapi#time-and-date
func Epoch(num, unit interface{}) Expr { return epochFn{Epoch: wrap(num), Unit: wrap(unit)} }

type epochFn struct {
	fnApply
	Epoch Expr `json:"epoch"`
	Unit  Expr `json:"unit"`
}

// Now returns the current snapshot time.
//
// Returns:
// Expr
//
// See: https://docs.fauna.com/fauna/current/api/fql/functions/now
func Now() Expr {
	return nowFn{Now: NullV{}}
}

type nowFn struct {
	fnApply
	Now Expr `json:"now" faunarepr:"noargs"`
}

// ToSeconds converts a time expression to seconds since the UNIX epoch.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - A time literal.
func ToSeconds(value interface{}) Expr {
	return toSecondsFn{ToSeconds: wrap(value)}
}

type toSecondsFn struct {
	fnApply
	ToSeconds Expr `json:"to_seconds"`
}

// ToMillis converts a time expression to milliseconds since the UNIX epoch.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - A time literal.
func ToMillis(value interface{}) Expr {
	return toMillisFn{ToMillis: wrap(value)}
}

type toMillisFn struct {
	fnApply
	ToMillis Expr `json:"to_millis"`
}

// ToMicros converts a time expression to microseconds since the UNIX epoch.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - A time literal.
func ToMicros(value interface{}) Expr {
	return toMicrosFn{ToMicros: wrap(value)}
}

type toMicrosFn struct {
	fnApply
	ToMicros Expr `json:"to_micros"`
}

// Year returns the time expression's year, following the ISO-8601 standard.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - year.
func Year(value interface{}) Expr {
	return yearFn{Year: wrap(value)}
}

type yearFn struct {
	fnApply
	Year Expr `json:"year"`
}

// Month returns a time expression's month of the year, from 1 to 12.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - Month.
func Month(value interface{}) Expr {
	return monthFn{Month: wrap(value)}
}

type monthFn struct {
	fnApply
	Month Expr `json:"month"`
}

// Hour returns a time expression's hour of the day, from 0 to 23.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - year.
func Hour(value interface{}) Expr {
	return hourFn{Hour: wrap(value)}
}

type hourFn struct {
	fnApply
	Hour Expr `json:"hour"`
}

// Minute returns a time expression's minute of the hour, from 0 to 59.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - year.
func Minute(value interface{}) Expr {
	return minuteFn{Minute: wrap(value)}
}

type minuteFn struct {
	fnApply
	Minute Expr `json:"minute"`
}

// Second returns a time expression's second of the minute, from 0 to 59.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - year.
func Second(value interface{}) Expr {
	return secondFn{Second: wrap(value)}
}

type secondFn struct {
	fnApply
	Second Expr `json:"second"`
}

// DayOfMonth returns a time expression's day of the month, from 1 to 31.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - day of month.
func DayOfMonth(value interface{}) Expr {
	return dayOfMonthFn{DayOfMonth: wrap(value)}
}

type dayOfMonthFn struct {
	fnApply
	DayOfMonth Expr `json:"day_of_month"`
}

// DayOfWeek returns a time expression's day of the week following ISO-8601 convention, from 1 (Monday) to 7 (Sunday).
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - day of week.
func DayOfWeek(value interface{}) Expr {
	return dayOfWeekFn{DayOfWeek: wrap(value)}
}

type dayOfWeekFn struct {
	fnApply
	DayOfWeek Expr `json:"day_of_week"`
}

// DayOfYear eturns a time expression's day of the year, from 1 to 365, or 366 in a leap year.
//
// Parameters:
//    value Object - The expression to convert.
//
// Returns:
//   time - Day of the year.
func DayOfYear(value interface{}) Expr {
	return dayOfYearFn{DayOfYear: wrap(value)}
}

type dayOfYearFn struct {
	fnApply
	DayOfYear Expr `json:"day_of_year"`
}
