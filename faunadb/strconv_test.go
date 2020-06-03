package faunadb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func assertString(t *testing.T, expr Expr, expected string) {
	str := expr.String()
	require.Equal(t, expected, str)
}

func TestStringifyFaunaValues(t *testing.T) {
	assertString(t, StringV("a string"), `"a string"`)
	assertString(t, LongV(90), `90`)

	arr := Arr{
		LongV(1), LongV(2), LongV(3),
		NullV{}, Obj{"a": "bcde"},
		BooleanV(true), BooleanV(false),
	}
	assertString(t,
		arr,
		`Arr{1, 2, 3, nil, Obj{"a": "bcde"}, true, false}`,
	)

	assertString(t,
		NullV{},
		`nil`,
	)

	assertString(t,
		SetRefV{Parameters: map[string]Value{"x": StringV("y")}},
		`SetRefV{ Parameters: map[string]Value{"x": "y"} }`,
	)

	assertString(t,
		nativeCollections,
		`RefV{ID: "collections"}`,
	)

	assertString(t,
		RefV{Database: &RefV{ID: "db1"}},
		`RefV{Database: &RefV{ID: "db1"}}`,
	)

	now := time.Now()
	assertString(t,
		TimeV(now),
		`TimeV("`+now.Format("2006-01-02T15:04:05.999999999Z")+`")`,
	)

}

func TestStringifyFixedArityFunctions(t *testing.T) {

	assertString(t, Add(1, 2, 3), `Add(1, 2, 3)`)

	assertString(t, Append(Arr{3, 4}, Arr{1, 2}), `Append(Arr{3, 4}, Arr{1, 2})`)
	assertString(t, Exists(Function("a_function")), `Exists(Function("a_function"))`)
	assertString(t, DayOfMonth(Epoch(0, "second")), `DayOfMonth(Epoch(0, "second"))`)
	assertString(t, DayOfMonth(Epoch(2147483648, "second")), `DayOfMonth(Epoch(2147483648, "second"))`)
	assertString(t, DayOfMonth(0), `DayOfMonth(0)`)
	assertString(t, DayOfMonth(2147483648000000), `DayOfMonth(2147483648000000)`)
	assertString(t, DayOfWeek(Epoch(0, "second")), `DayOfWeek(Epoch(0, "second"))`)
	assertString(t, DayOfWeek(Epoch(2147483648, "second")), `DayOfWeek(Epoch(2147483648, "second"))`)
	assertString(t, DayOfWeek(0), `DayOfWeek(0)`)
	assertString(t, DayOfWeek(2147483648000000), `DayOfWeek(2147483648000000)`)
	assertString(t, DayOfYear(Epoch(0, "second")), `DayOfYear(Epoch(0, "second"))`)
	assertString(t, DayOfYear(Epoch(2147483648, "second")), `DayOfYear(Epoch(2147483648, "second"))`)
	assertString(t, DayOfYear(0), `DayOfYear(0)`)
	assertString(t, DayOfYear(2147483648000000), `DayOfYear(2147483648000000)`)

	assertString(t, Drop(2, Arr{1, 2, 3}), `Drop(2, Arr{1, 2, 3})`)
	assertString(t, Query(Lambda("x", Var("x"))), `Query(Lambda("x", Var("x")))`)
	assertString(t, Abs(-2), `Abs(-2)`)
	assertString(t, Acos(1), `Acos(1)`)
	assertString(t, Add(Arr{2, 3}), `Add(2, 3)`)
	assertString(t, And(Arr{true, true}), `And(true, true)`)
	assertString(t, Arr{Any(Arr{true, true, false}), All(Arr{true, true, true}), Any(Arr{false, false, false}), All(Arr{true, true, false})}, `Arr{Any(Arr{true, true, false}), All(Arr{true, true, true}), Any(Arr{false, false, false}), All(Arr{true, true, false})}`)
	assertString(t, Trunc(Asin(0.5)), `Trunc(Asin(0.5))`)
	assertString(t, Trunc(Atan(0.5)), `Trunc(Atan(0.5))`)
	assertString(t, BitAnd(Arr{2, 3}), `BitAnd(2, 3)`)
	assertString(t, BitNot(2), `BitNot(2)`)
	assertString(t, BitOr(Arr{2, 1}), `BitOr(2, 1)`)
	assertString(t, BitXor(Arr{2, 3}), `BitXor(2, 3)`)
	assertString(t, Casefold("GET DOWN"), `Casefold("GET DOWN")`)

	assertString(t, Ceil(1.8e+00), `Ceil(1.8)`)
	assertString(t, Concat(Arr{"Hello", "World"}), `Concat(Arr{"Hello", "World"})`)
	assertString(t, ContainsStr("faunadb", "fauna"), `ContainsStr("faunadb", "fauna")`)
	assertString(t, ContainsStrRegex("faunadb", "f(\\w+)a"), `ContainsStrRegex("faunadb", "f(\\w+)a")`)
	assertString(t, ContainsStrRegex("faunadb", "/^\\d*\\.\\d+$/"), `ContainsStrRegex("faunadb", "/^\\d*\\.\\d+$/")`)
	assertString(t, ContainsStrRegex("test data", "\\s"), `ContainsStrRegex("test data", "\\s")`)
	assertString(t, Trunc(Cosh(0.5)), `Trunc(Cosh(0.5))`)
	assertString(t, Arr{Count(Arr{1, 2, 3, 4, 5, 6, 7, 8, 9}), Mean(Arr{1, 2, 3, 4, 5, 6, 7, 8, 9}), Sum(Arr{1, 2, 3, 4, 5, 6, 7, 8, 9})}, `Arr{Count(Arr{1, 2, 3, 4, 5, 6, 7, 8, 9}), Mean(Arr{1, 2, 3, 4, 5, 6, 7, 8, 9}), Sum(Arr{1, 2, 3, 4, 5, 6, 7, 8, 9})}`)
	assertString(t, Arr{Count(Match(Index("countmeansum_idx"))), Trunc(Mean(Match(Index("countmeansum_idx")))), Sum(Match(Index("countmeansum_idx")))}, `Arr{Count(Match(Index("countmeansum_idx"))), Trunc(Mean(Match(Index("countmeansum_idx")))), Sum(Match(Index("countmeansum_idx")))}`)
	assertString(t, Date("1970-01-02"), `Date("1970-01-02")`)
	assertString(t, Trunc(Degrees(0.5)), `Trunc(Degrees(0.5))`)
	assertString(t, Divide(Arr{10, 2}), `Divide(10, 2)`)
	assertString(t, Select(Arr{0}, Count(Paginate(Documents(Collection("collection_2276634055"))))), `Select(Arr{0}, Count(Paginate(Documents(Collection("collection_2276634055")))))`)
	assertString(t, Count(Documents(Collection("collection_2276634055"))), `Count(Documents(Collection("collection_2276634055")))`)
	assertString(t, EndsWith("faunadb", "fauna"), `EndsWith("faunadb", "fauna")`)
	assertString(t, EndsWith("faunadb", "db"), `EndsWith("faunadb", "db")`)
	assertString(t, EndsWith("faunadb", ""), `EndsWith("faunadb", "")`)
	assertString(t, Arr{Epoch(30, "second"), Epoch(30, "millisecond"), Epoch(30, "microsecond"), Epoch(30, "nanosecond")}, `Arr{Epoch(30, "second"), Epoch(30, "millisecond"), Epoch(30, "microsecond"), Epoch(30, "nanosecond")}`)
	assertString(t, Equals(Arr{"fire", "fire"}), `Equals("fire", "fire")`)
	assertString(t, Trunc(Exp(2)), `Trunc(Exp(2))`)
	assertString(t, FindStr("GET DOWN", "DOWN"), `FindStr("GET DOWN", "DOWN")`)
	assertString(t, Floor(2.99), `Floor(2.99)`)
	assertString(t, Format("%2$s%1$s %3$s", Arr{"DB", "Fauna", "rocks"}), `Format("%2$s%1$s %3$s", Arr{"DB", "Fauna", "rocks"})`)
	assertString(t, Format("%d %s %.2f %%", Arr{34, "tEsT ", 3.14159}), `Format("%d %s %.2f %%", Arr{34, "tEsT ", 3.14159})`)
	assertString(t, GTE(Arr{2, 2}), `GTE(2, 2)`)
	assertString(t, GT(Arr{3, 2}), `GT(3, 2)`)

	assertString(t, If(true, "true", "false"), `If(true, "true", "false")`)
	assertString(t, LTE(Arr{2, 2}), `LTE(2, 2)`)
	assertString(t, LT(Arr{2, 3}), `LT(2, 3)`)
	assertString(t, LTrim("    One Fish Two Fish"), `LTrim("    One Fish Two Fish")`)
	assertString(t, Length("One Fish Two Fish"), `Length("One Fish Two Fish")`)

	assertString(t, Trunc(Ln(2)), `Trunc(Ln(2))`)
	assertString(t, Log(100), `Log(100)`)
	assertString(t, LowerCase("One Fish Two Fish"), `Lowercase("One Fish Two Fish")`)
	assertString(t, Max(Arr{2, 3}), `Max(2, 3)`)
	assertString(t, Min(Arr{4, 3}), `Min(4, 3)`)
	assertString(t, Modulo(Arr{10, 2}), `Modulo(10, 2)`)
	assertString(t, Multiply(Arr{2, 3}), `Multiply(2, 3)`)
	assertString(t, Not(false), `Not(false)`)
	assertString(t, Or(Arr{false, true}), `Or(false, true)`)
	assertString(t, RTrim("One Fish Two Fish   "), `RTrim("One Fish Two Fish   ")`)
	assertString(t, Trunc(Radians(2)), `Trunc(Radians(2))`)
	assertString(t, RegexEscape("f(\\w+)a"), `RegexEscape("f(\\w+)a")`)
	assertString(t, ReplaceStr("One Fish Two Fish", "Fish", "Dog"), `ReplaceStr("One Fish Two Fish", "Fish", "Dog")`)
	assertString(t, ReplaceStrRegex("One FIsh Two fish", "[Ff][Ii]sh", "Dog"), `ReplaceStrRegex("One FIsh Two fish", "[Ff][Ii]sh", "Dog")`)
	assertString(t, Sign(-1), `Sign(-1)`)
	assertString(t, Trunc(Sin(20)), `Trunc(Sin(20))`)
	assertString(t, Trunc(Sinh(0.5)), `Trunc(Sinh(0.5))`)
	assertString(t, Space(5), `Space(5)`)
	assertString(t, Sqrt(16), `Sqrt(16)`)
	assertString(t, StartsWith("faunadb", "fauna"), `StartsWith("faunadb", "fauna")`)
	assertString(t, StartsWith("faunadb", "F"), `StartsWith("faunadb", "F")`)
	assertString(t, SubString("ABCDEF", 3), `SubString("ABCDEF", 3)`)
	assertString(t, SubString("ABCDEF", -2), `SubString("ABCDEF", -2)`)
	assertString(t, Subtract(Arr{2, 3}), `Subtract(2, 3)`)
	assertString(t, Trunc(Tan(20)), `Trunc(Tan(20))`)
	assertString(t, Trunc(Tanh(0.5)), `Trunc(Tanh(0.5))`)
	assertString(t, TimeAdd(Epoch(0, "second"), 1, "hour"), `TimeAdd(Epoch(0, "second"), 1, "hour")`)
	assertString(t, TimeAdd(Epoch(0, "second"), 16, "minute"), `TimeAdd(Epoch(0, "second"), 16, "minute")`)
	assertString(t, TimeDiff(Epoch(0, "second"), Epoch(1, "second"), "second"), `TimeDiff(Epoch(0, "second"), Epoch(1, "second"), "second")`)
	assertString(t, TimeDiff(Epoch(24, "hour"), Epoch(1, "day"), "hour"), `TimeDiff(Epoch(24, "hour"), Epoch(1, "day"), "hour")`)
	assertString(t, Time("1970-01-01T00:00:00-04:00"), `Time("1970-01-01T00:00:00-04:00")`)
	assertString(t, TimeSubtract(Epoch(190, "hour"), 78, "minute"), `TimeSubtract(Epoch(190, "hour"), 78, "minute")`)
	assertString(t, TimeSubtract(Epoch(16, "second"), 16, "second"), `TimeSubtract(Epoch(16, "second"), 16, "second")`)
	assertString(t, TitleCase("onE Fish tWO FiSh"), `Titlecase("onE Fish tWO FiSh")`)
	assertString(t, ToDate("1970-01-02"), `ToDate("1970-01-02")`)
	assertString(t, ToNumber("42"), `ToNumber("42")`)
	assertString(t, ToString(42), `ToString(42)`)
	assertString(t, ToInteger(3.14), `ToInteger(3.14)`)
	assertString(t, ToDouble(90), `ToDouble(90)`)
	assertString(t, ToObject(Arr{Arr{"x", 1}}), `ToObject(Arr{Arr{"x", 1}})`)
	assertString(t, ToArray(Obj{"x": 1}), `ToArray(Obj{"x": 1}))`)
	assertString(t, ToTime("1970-01-01T00:00:00-04:00"), `ToTime("1970-01-01T00:00:00-04:00")`)
	assertString(t, Trim("   One Fish Two Fish   "), `Trim("   One Fish Two Fish   ")`)
	assertString(t, Trunc(1.234567), `Trunc(1.234567)`)
	assertString(t, UpperCase("One Fish Two Fish"), `UpperCase("One Fish Two Fish")`)
	assertString(t, Hour(Epoch(0, "second")), `Hour(Epoch(0, "second"))`)
	assertString(t, Hour(Epoch(2147483648, "second")), `Hour(Epoch(2147483648, "second"))`)
	assertString(t, Hour(0), `Hour(0)`)
	assertString(t, Hour(2147483648000000), `Hour(2147483648000000)`)
	assertString(t, IsEmpty(Arr{}), `IsEmpty(Arr{})`)
	assertString(t, IsEmpty(Arr{1}), `IsEmpty(Arr{1})`)
	assertString(t, IsNonEmpty(Arr{}), `IsNonEmpty(Arr{})`)
	assertString(t, IsNonEmpty(Arr{1, 2}), `IsNonEmpty(Arr{1, 2})`)

	assertString(t, Minute(Epoch(0, "second")), `Minute(Epoch(0, "second"))`)
	assertString(t, Minute(Epoch(2147483648, "second")), `Minute(Epoch(2147483648, "second"))`)
	assertString(t, Minute(0), `Minute(0)`)
	assertString(t, Minute(2147483648000000), `Minute(2147483648000000)`)
	assertString(t, Month(Epoch(0, "second")), `Month(Epoch(0, "second"))`)
	assertString(t, Month(Epoch(2147483648, "second")), `Month(Epoch(2147483648, "second"))`)
	assertString(t, Month(0), `Month(0)`)
	assertString(t, Month(2147483648000000), `Month(2147483648000000)`)
	assertString(t, Exists(ScopedCollection("collection_1675362432", ScopedDatabase("child_3306182225", Database("parent_3146154786")))), `Exists(ScopedCollection("collection_1675362432", ScopedDatabase("child_3306182225", Database("parent_3146154786"))))`)
	assertString(t, Prepend(Arr{1, 2}, Arr{3, 4}), `Prepend(Arr{1, 2}, Arr{3, 4})`)
	assertString(t, Select("data", Paginate(Range(Match(Index("range_idx")), 3, 8))), `Select("data", Paginate(Range(Match(Index("range_idx")), 3, 8)))`)
	assertString(t, Select("data", Paginate(Range(Match(Index("range_idx")), 17, 18))), `Select("data", Paginate(Range(Match(Index("range_idx")), 17, 18)))`)
	assertString(t, Select("data", Paginate(Range(Match(Index("range_idx")), 19, 0))), `Select("data", Paginate(Range(Match(Index("range_idx")), 19, 0)))`)
	assertString(t, Reduce(Lambda(Arr{"accum", "value"}, Add(Arr{Var("accum"), Var("value")})), 0, Arr{1, 2, 3, 4, 5, 6, 7, 8, 9}), `Reduce(Lambda(Arr{"accum", "value"}, Add(Var("accum"), Var("value"))), 0, Arr{1, 2, 3, 4, 5, 6, 7, 8, 9})`)
	assertString(t, Reduce(Lambda(Arr{"accum", "value"}, Concat(Arr{Var("accum"), Var("value")})), "", Arr{"Fauna", "DB", " ", "rocks"}), `Reduce(Lambda(Arr{"accum", "value"}, Concat(Arr{Var("accum"), Var("value")})), "", Arr{"Fauna", "DB", " ", "rocks"})`)
	assertString(t, Second(Epoch(0, "second")), `Second(Epoch(0, "second"))`)
	assertString(t, Second(Epoch(2147483648, "second")), `Second(Epoch(2147483648, "second"))`)
	assertString(t, Second(0), `Second(0)`)
	assertString(t, Second(2147483648000000), `Second(2147483648000000)`)
	assertString(t, Take(2, Arr{1, 2, 3}), `Take(2, Arr{1, 2, 3})`)
	assertString(t, ToMillis(Epoch(0, "second")), `ToMillis(Epoch(0, "second"))`)
	assertString(t, ToMillis(Epoch(2147483648000000, "microsecond")), `ToMillis(Epoch(2147483648000000, "microsecond"))`)
	assertString(t, ToMillis(0), `ToMillis(0)`)
	assertString(t, ToMillis(2147483648000000), `ToMillis(2147483648000000)`)
	assertString(t, ToMillis(Epoch(0, "second")), `ToMillis(Epoch(0, "second"))`)
	assertString(t, ToMillis(Epoch(2147483648000, "millisecond")), `ToMillis(Epoch(2147483648000, "millisecond"))`)
	assertString(t, ToMillis(0), `ToMillis(0)`)
	assertString(t, ToMillis(2147483648000000), `ToMillis(2147483648000000)`)
	assertString(t, ToSeconds(0), `ToSeconds(0)`)
	assertString(t, ToSeconds(Epoch(2147483648, "second")), `ToSeconds(Epoch(2147483648, "second"))`)
	assertString(t, ToSeconds(0), `ToSeconds(0)`)
	assertString(t, ToSeconds(2147483648000000), `ToSeconds(2147483648000000)`)
	assertString(t, Year(Epoch(0, "second")), `Year(Epoch(0, "second"))`)
	assertString(t, Year(Epoch(2147483648, "second")), `Year(Epoch(2147483648, "second"))`)
	assertString(t, Year(0), `Year(0)`)
	assertString(t, Year(2147483648000000), `Year(2147483648000000)`)
}

func TestStringify0ArityFunctions(t *testing.T) {
	assertString(t, Now(), `Now()`)

	assertString(t, NewId(), `NewId()`)
	assertString(t, NextID(), `NextID()`)
	assertString(t, Indexes(), `Indexes()`)
	assertString(t, Functions(), `Functions()`)
	assertString(t, Roles(), `Roles()`)
	assertString(t, Credentials(), `Credentials()`)
	assertString(t, Tokens(), `Tokens()`)
	assertString(t, Collections(), `Collections()`)
	assertString(t, Databases(), `Databases()`)
	assertString(t, Keys(), `Keys()`)
	assertString(t, Classes(), `Classes()`)
}

func TestStringifyVariableFunctions(t *testing.T) {
	assertString(t,
		Paginate(Match(Index("idx")), After(Now()), Size(1), Before("yesterday"), EventsOpt(9), TS(Now())),
		`Paginate(Match(Index("idx")), After(Now()), Before("yesterday"), EventsOpt(9), Size(1), TS(Now()))`,
	)

	assertString(t, Concat(Arr{"hey", "there"}), `Concat(Arr{"hey", "there"})`)
	assertString(t, Concat(Arr{"hey", "there"}, Separator("/")), `Concat(Arr{"hey", "there"}, Separator("/"))`)

	assertString(t,
		ReplaceStrRegex("One FIsh Two fish", "[Ff][Ii]sh", "Dog", OnlyFirst()),
		`ReplaceStrRegex("One FIsh Two fish", "[Ff][Ii]sh", "Dog", OnlyFirst())`,
	)

	assertString(t, SubString("ABCDEFZZZ", 3, StrLength(3)), `SubString("ABCDEFZZZ", 3, StrLength(3))`)

	merge1 := Merge(Obj{}, Obj{"a": 1})
	merge2 := Merge(Obj{"b": 2},
		Obj{},
		ConflictResolver(Lambda(Arr{"key", "left", "right"}, Var("right"))),
	)
	assertString(t, merge1, `Merge(Obj{}, Obj{"a": 1})`)
	assertString(t, merge2, `Merge(Obj{"b": 2}, Obj{}, ConflictResolver(Lambda(Arr{"key", "left", "right"}, Var("right"))))`)

	assertString(t, FindStr("One Fish Two Fish", "Fish", Start(8)), `FindStr("One Fish Two Fish", "Fish", Start(8))`)
	assertString(t, Concat(Arr{"Hello", "World"}, Separator(" ")), `Concat(Arr{"Hello", "World"}, Separator(" "))`)

	assertString(t, Trunc(0.5, Precision(2)), `Trunc(0.5, Precision(2))`)
	assertString(t, Round(1.666666, Precision(3)), `Round(1.666666, Precision(3))`)

	assertString(t, Casefold("Å", Normalizer("NFD")), `Casefold("Å", Normalizer("NFD"))`)
	assertString(t, Casefold("Å", Normalizer("NFC")), `Casefold("Å", Normalizer("NFC"))`)
	assertString(t, Casefold("ẛ̣", Normalizer("NFKD")), `Casefold("ẛ̣", Normalizer("NFKD"))`)
	assertString(t, Casefold("ẛ̣", Normalizer("NFKC")), `Casefold("ẛ̣", Normalizer("NFKC"))`)
	assertString(t, Casefold("Å", Normalizer("NFKCCaseFold")), `Casefold("Å", Normalizer("NFKCCaseFold"))`)

	assertString(t, Exists(Function("a_function"), TS(0)), `Exists(Function("a_function"), TS(0))`)
	assertString(t, Get(Function("a_function"), TS(0)), `Get(Function("a_function"), TS(0))`)
}

func TestStringifyVarArgFunctions(t *testing.T) {
	assertString(t, Add(1, 2, 3), `Add(1, 2, 3)`)
	assertString(t, Divide(1, 2, 3), `Divide(1, 2, 3)`)
	assertString(t, Multiply(1, 2, 3), `Multiply(1, 2, 3)`)

	assertString(t, Add(Var("x")), `Add(Var("x"))`)

	assertString(t, BitOr(1, 2, 3), `BitOr(1, 2, 3)`)
	assertString(t, Do(Arr{1, 2, 3}), `Do(Arr{1, 2, 3})`)
	assertString(t, Equals(0, Modulo(Var("i"), 2)), `Equals(0, Modulo(Var("i"), 2))`)

	assertString(t, Or(1, 2, 3), `Or(1, 2, 3)`)
	assertString(t, And(1, 2, 3), `And(1, 2, 3)`)
	assertString(t, BitOr(1, 2, 3), `BitOr(1, 2, 3)`)
	assertString(t, BitXor(1, 2, 3), `BitXor(1, 2, 3)`)

	assertString(t, Max(1, Var("x"), 3), `Max(1, Var("x"), 3)`)
	assertString(t, Min(1, Var("x"), 3), `Min(1, Var("x"), 3)`)
	assertString(t, Modulo(1, Var("x"), 3), `Modulo(1, Var("x"), 3)`)

	assertString(t, LT(1, Var("x"), 3), `LT(1, Var("x"), 3)`)
	assertString(t, GT(1, Var("x"), 3), `GT(1, Var("x"), 3)`)
	assertString(t, LTE(1, Var("x"), 3), `LTE(1, Var("x"), 3)`)
	assertString(t, GTE(1, Var("x"), 3), `GTE(1, Var("x"), 3)`)

	assertString(t, Paginate(Union(
		MatchTerm(Null(), "arcane"),
		MatchTerm(Null(), "fire"),
	)),
		`Paginate(Union(MatchTerm(nil, "arcane"), MatchTerm(nil, "fire")))`)

	assertString(t, Paginate(Difference(
		MatchTerm(Null(), "arcane"),
		MatchTerm(Null(), "fire"),
	)),
		`Paginate(Difference(MatchTerm(nil, "arcane"), MatchTerm(nil, "fire")))`)

	assertString(t, Paginate(Intersection(
		MatchTerm(Null(), "arcane"),
		MatchTerm(Null(), "fire"),
	)),
		`Paginate(Intersection(MatchTerm(nil, "arcane"), MatchTerm(nil, "fire")))`)
}

func TestStringifyScopedFunctions(t *testing.T) {
	assertString(t, NewId(), `NewId()`)
	assertString(t, NextID(), `NextID()`)
	assertString(t, ScopedIndexes(Database("db1")), `ScopedIndexes(Database("db1"))`)
	assertString(t, ScopedFunctions(Database("db1")), `ScopedFunctions(Database("db1"))`)
	assertString(t, ScopedRoles(Database("db1")), `ScopedRoles(Database("db1"))`)
	assertString(t, ScopedCredentials(Database("db1")), `ScopedCredentials(Database("db1"))`)
	assertString(t, ScopedTokens(Database("db1")), `ScopedTokens(Database("db1"))`)
	assertString(t, ScopedCollections(Database("db1")), `ScopedCollections(Database("db1"))`)
	assertString(t, ScopedDatabases(Database("db1")), `ScopedDatabases(Database("db1"))`)
	assertString(t, ScopedKeys(Database("db1")), `ScopedKeys(Database("db1"))`)
	assertString(t, ScopedClasses(Database("db1")), `ScopedClasses(Database("db1"))`)
}

func TestStringifyCustomFunctions(t *testing.T) {
	let := Let().Bind("x", 90).Bind("y", 65).In(Add(Var("x"), Var("y")))
	assertString(t, let, `Let().Bind("x", 90).Bind("y", 65).In(Add(Var("x"), Var("y")))`)

	assertString(t, Match(Index("idx")), `Match(Index("idx"))`)
	assertString(t, MatchTerm(Index("idx"), "term"), `MatchTerm(Index("idx"), "term")`)
}
