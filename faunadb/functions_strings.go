package faunadb

// String

// Format formats values into a string.
//
// Parameters:
//  format string - format a string with format specifiers.
//
// Optional parameters:
//  values []string - list of values to format into string.
//
// Returns:
//  string - A string.
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func Format(format interface{}, values ...interface{}) Expr {
	return formatFn{Format: wrap(format), Values: wrap(varargs(values...))}
}

type formatFn struct {
	fnApply
	Format Expr `json:"format"`
	Values Expr `json:"values"`
}

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
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func Concat(terms interface{}, options ...OptionalParameter) Expr {
	fn := concatFn{Concat: wrap(terms)}
	return applyOptionals(fn, options)
}

type concatFn struct {
	fnApply
	Concat    Expr `json:"concat"`
	Separator Expr `json:"separator,omitempty" faunarepr:"optfn"`
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
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func Casefold(str interface{}, options ...OptionalParameter) Expr {
	fn := casefoldFn{Casefold: wrap(str)}
	return applyOptionals(fn, options)
}

type casefoldFn struct {
	fnApply
	Casefold   Expr `json:"casefold"`
	Normalizer Expr `json:"normalizer,omitempty" faunarepr:"optfn"`
}

// StartsWith returns true if the string starts with the given prefix value, or false if otherwise
//
// Parameters:
//
//  value  string -  the string to evaluate
//  search string -  the prefix to search for
//
// Returns:
//   boolean       - does `value` start with `search
//
// See https://docs.fauna.com/fauna/current/api/fql/functions/startswith
func StartsWith(value interface{}, search interface{}) Expr {
	return startsWithFn{StartsWith: wrap(value), Search: wrap(search)}
}

type startsWithFn struct {
	fnApply
	StartsWith Expr `json:"startswith"`
	Search     Expr `json:"search"`
}

// EndsWith returns true if the string ends with the given suffix value, or false if otherwise
//
// Parameters:
//
// value  string  -  the string to evaluate
// search  string -  the suffix to search for
//
// Returns:
// boolean       - does `value` end with `search`
//
// See https://docs.fauna.com/fauna/current/api/fql/functions/endswith
func EndsWith(value interface{}, search interface{}) Expr {
	return endsWithFn{EndsWith: wrap(value), Search: wrap(search)}
}

type endsWithFn struct {
	fnApply
	EndsWith Expr `json:"endswith"`
	Search   Expr `json:"search"`
}

// ContainsStr returns true if the string contains the given substring, or false if otherwise
//
// Parameters:
//
// value string  -  the string to evaluate
// search string -  the substring to search for
//
// Returns:
// boolean      - was the search result found
//
// See https://docs.fauna.com/fauna/current/api/fql/functions/containsstr
func ContainsStr(value interface{}, search interface{}) Expr {
	return containsStrFn{ContainsStr: wrap(value), Search: wrap(search)}
}

type containsStrFn struct {
	fnApply
	ContainsStr Expr `json:"containsstr"`
	Search      Expr `json:"search"`
}

// ContainsStrRegex returns true if the string contains the given pattern, or false if otherwise
//
// Parameters:
//
// value   string      -  the string to evaluate
// pattern string      -  the pattern to search for
//
// Returns:
// boolean      - was the search result found
//
// See https://docs.fauna.com/fauna/current/api/fql/functions/containsstrregex
func ContainsStrRegex(value interface{}, pattern interface{}) Expr {
	return containsStrRegexFn{ContainsStrRegex: wrap(value), Pattern: wrap(pattern)}
}

type containsStrRegexFn struct {
	fnApply
	ContainsStrRegex Expr `json:"containsstrregex"`
	Pattern          Expr `json:"pattern"`
}

// RegexEscape It takes a string and returns a regex which matches the input string verbatim.
//
// Parameters:
//
// value  string     - the string to analyze
// pattern       -  the pattern to search for
//
// Returns:
// boolean      - was the search result found
//
// See https://docs.fauna.com/fauna/current/api/fql/functions/regexescape
func RegexEscape(value interface{}) Expr {
	return regexEscapeFn{RegexEscape: wrap(value)}
}

type regexEscapeFn struct {
	fnApply
	RegexEscape Expr `json:"regexescape"`
}

// FindStr locates a substring in a source string.  Optional parameters: Start
//
// Parameters:
//  str string  - The source string
//  find string - The string to locate
//
// Optional parameters:
//  start int - a position to start the search. See Start() function.
//
// Returns:
//  string - The offset of where the substring starts or -1 if not found
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func FindStr(str, find interface{}, options ...OptionalParameter) Expr {
	fn := findStrFn{FindStr: wrap(str), Find: wrap(find)}
	return applyOptionals(fn, options)
}

type findStrFn struct {
	fnApply
	FindStr Expr `json:"findstr"`
	Find    Expr `json:"find"`
	Start   Expr `json:"start,omitempty" faunarepr:"optfn"`
}

// FindStrRegex locates a java regex pattern in a source string.  Optional parameters: Start
//
// Parameters:
//  str string      - The sourcestring
//  pattern string  - The pattern to locate.
//
// Optional parameters:
//  start long - a position to start the search.  See Start() function.
//
// Returns:
//  string - The offset of where the substring starts or -1 if not found
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func FindStrRegex(str, pattern interface{}, options ...OptionalParameter) Expr {
	fn := findStrRegexFn{FindStrRegex: wrap(str), Pattern: wrap(pattern)}
	return applyOptionals(fn, options)
}

type findStrRegexFn struct {
	fnApply
	FindStrRegex Expr `json:"findstrregex"`
	Pattern      Expr `json:"pattern"`
	Start        Expr `json:"start,omitempty" faunarepr:"optfn"`
	NumResults   Expr `json:"num_results,omitempty" faunarepr:"optfn"`
}

// Length finds the length of a string in codepoints
//
// Parameters:
//  str string - A string to find the length in codepoints
//
// Returns:
//  int - A length of a string.
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func Length(str interface{}) Expr { return lengthFn{Length: wrap(str)} }

type lengthFn struct {
	fnApply
	Length Expr `json:"length"`
}

// LowerCase changes all characters in the string to lowercase
//
// Parameters:
//  str string - A string to convert to lowercase
//
// Returns:
//  string - A string in lowercase.
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func LowerCase(str interface{}) Expr { return lowercaseFn{Lowercase: wrap(str)} }

type lowercaseFn struct {
	fnApply
	Lowercase Expr `json:"lowercase"`
}

// LTrim returns a string wtih leading white space removed.
//
// Parameters:
//  str string - A string to remove leading white space
//
// Returns:
//  string - A string with all leading white space removed
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func LTrim(str interface{}) Expr { return lTrimFn{LTrim: wrap(str)} }

type lTrimFn struct {
	fnApply
	LTrim Expr `json:"ltrim"`
}

// Repeat returns a string wtih repeated n times
//
// Parameters:
//  str string - A string to repeat
//  number int - The number of times to repeat the string
//
// Optional parameters:
//  Number - Only replace the first found pattern.  See OnlyFirst() function.
//
// Returns:
//  string - A string concatendanted the specified number of times
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func Repeat(str interface{}, options ...OptionalParameter) Expr {
	fn := repeatFn{Repeat: wrap(str)}
	return applyOptionals(fn, options)
}

type repeatFn struct {
	fnApply
	Repeat Expr `json:"repeat"`
	Number Expr `json:"number,omitempty" faunarepr:"fn=optfn,name=Number"`
}

// ReplaceStr returns a string with every occurence of the "find" string changed to "replace" string
//
// Parameters:
//  str string     - A source string
//  find string    - The substring to locate in in the source string
//  replace string - The string to replaice the "find" string when located
//
// Returns:
//  string - returns a string with every occurence of the "find" string changed to "replace"
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func ReplaceStr(str, find, replace interface{}) Expr {
	return replaceStrFn{
		ReplaceStr: wrap(str),
		Find:       wrap(find),
		Replace:    wrap(replace),
	}
}

type replaceStrFn struct {
	fnApply
	ReplaceStr Expr `json:"replacestr"`
	Find       Expr `json:"find"`
	Replace    Expr `json:"replace"`
}

// ReplaceStrRegex returns a string with occurence(s) of the java regular expression "pattern" changed to "replace" string.   Optional parameters: OnlyFirst
//
// Parameters:
//  value string   - The source string
//  pattern string - A java regular expression to locate
//  replace string - The string to replace the pattern when located
//
// Optional parameters:
//  OnlyFirst - Only replace the first found pattern.  See OnlyFirst() function.
//
// Returns:
//  string - A string with occurence(s) of the java regular expression "pattern" changed to "replace" string
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func ReplaceStrRegex(value, pattern, replace interface{}, options ...OptionalParameter) Expr {
	fn := replaceStrRegexFn{
		ReplaceStrRegex: wrap(value),
		Pattern:         wrap(pattern),
		Replace:         wrap(replace),
	}
	return applyOptionals(fn, options)
}

type replaceStrRegexFn struct {
	fnApply
	ReplaceStrRegex Expr `json:"replacestrregex"`
	Pattern         Expr `json:"pattern"`
	Replace         Expr `json:"replace"`
	First           Expr `json:"first,omitempty" faunarepr:"fn=optfn,name=OnlyFirst,noargs=true"`
}

// RTrim returns a string wtih trailing white space removed.
//
// Parameters:
//  str string - A string to remove trailing white space
//
// Returns:
//  string - A string with all trailing white space removed
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func RTrim(str interface{}) Expr { return rTrimFn{RTrim: wrap(str)} }

type rTrimFn struct {
	fnApply
	RTrim Expr `json:"rtrim"`
}

// Space function returns "N" number of spaces
//
// Parameters:
//  value int - the number of spaces
//
// Returns:
//  string - function returns string with n spaces
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func Space(value interface{}) Expr { return spaceFn{Space: wrap(value)} }

type spaceFn struct {
	fnApply
	Space Expr `json:"space"`
}

// SubString returns a subset of the source string.   Optional parameters: StrLength
//
// Parameters:
//  str string - A source string
//  start int  - The position in the source string where SubString starts extracting characters
//
// Optional parameters:
//  StrLength int - A value for the length of the extracted substring. See StrLength() function.
//
// Returns:
//  string - function returns a subset of the source string
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func SubString(str, start interface{}, options ...OptionalParameter) Expr {
	fn := subStringFn{SubString: wrap(str), Start: wrap(start)}
	return applyOptionals(fn, options)
}

type subStringFn struct {
	fnApply
	SubString Expr `json:"substring"`
	Start     Expr `json:"start"`
	Length    Expr `json:"length,omitempty" faunarepr:"fn=optfn,name=StrLength"`
}

// TitleCase changes all characters in the string to TitleCase
//
// Parameters:
//  str string - A string to convert to TitleCase
//
// Returns:
//  string - A string in TitleCase.
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func TitleCase(str interface{}) Expr { return titleCaseFn{Titlecase: wrap(str)} }

type titleCaseFn struct {
	fnApply
	Titlecase Expr `json:"titlecase"`
}

// Trim returns a string wtih trailing white space removed.
//
// Parameters:
//  str string - A string to remove trailing white space
//
// Returns:
//  string - A string with all trailing white space removed
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func Trim(str interface{}) Expr { return trimFn{Trim: wrap(str)} }

type trimFn struct {
	fnApply
	Trim Expr `json:"trim"`
}

// UpperCase changes all characters in the string to uppercase
//
// Parameters:
//  string - A string to convert to uppercase
//
// Returns:
//  string - A string in uppercase.
//
// See: https://app.fauna.com/documentation/reference/queryapi#string-functions
func UpperCase(str interface{}) Expr { return upperCaseFn{UpperCase: wrap(str)} }

type upperCaseFn struct {
	fnApply
	UpperCase Expr `json:"uppercase"`
}
