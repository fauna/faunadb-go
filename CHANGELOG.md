# v4.0.0 (April, 2021) [current]

- Add document streaming.
- Add third-party authentication functions: AccessProvider, AccessProviders, CreateAccessProvider,
  CurrentIdentity, CurrentToken, HasCurrentIdentity, HasCurrentToken.
- Add `omitempty` tag support
- Add support for partners info via request headers.

# v3.0.0 (August, 2020)

- Added Reverse()
- Add ContainsPath(), ContainsField(), ContainsValue()
- Add ToArray(), ToObject(), ToInteger(), ToDouble() functions
- Deprecate Contains()
- Add tests for versioned queries
- Bump apiVersion to 3
- Fix DoubleV string formatting

# v2.12.1 (May, 2020)

- Make base64 encoding of secret match other drivers.

# v2.12.0 (May, 2020)

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
