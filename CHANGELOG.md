
# v2.12.0 (May, 2020) [current]

- Add client specified query timeout
- Update Ref() to take two arguments
- Handle nil as NullV{}
- Add type check functions:
   IsEmpty(), IsNonEmpty(), IsNumber(), IsDouble(), IsInteger()
   IsBoolean(), IsNull(), IsBytes(), IsTimestamp(), IsDate()
   IsString(), IsArray(), IsObject(), IsRef(), IsSet(), IsDoc()
   IsLambda(), IsCollection(), IsDatabase(), IsIndex(), IsFunction()
   IsKey(), IsToken(), IsCredentials(), IsRole()

# v2.11.0 (February, 2020) [current]

- Add StartsWith(), EndsWith(), ContainsStr(), ContainsStrRegex(), RegexEscape()
- Add TimeAdd(), TimeSubtract(), TimeDiff()
- Add Count(), Sum(), Mean()
- Add Documents(), Now()
- Add Any(), All()

# v2.10.0 (November, 2019)

# v2.9.0 (October, 2019)

- Add CHANGELOG.md
- Add Range(), Reduce(), Merge() functions
- Add MoveDatabase()
- Add Format() function
- Send "X-Fauna-Driver: Go" header with http requests
