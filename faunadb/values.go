package faunadb

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

/*
Value represents valid FaunaDB values returned from the server. Values also implement the Expr interface.
They can be sent back and forth to FaunaDB with no extra escaping needed.

The Get method is used to decode a FaunaDB value into a Go type. For example:

	var t time.Time

	faunaTime, _ := client.Query(Time("now"))
	_ := faunaTime.Get(&t)

The At method uses field extractors to traverse the data to specify a field:

	var firstEmail string

	profile, _ := client.Query(RefCollection(Collection("profile), "43"))
	profile.At(ObjKey("emails").AtIndex(0)).Get(&firstEmail)

For more information, check https://app.fauna.com/documentation/reference/queryapi#simple-type.
*/
type Value interface {
	Expr
	Get(interface{}) error // Decode a FaunaDB value into a native Go type
	At(Field) FieldValue   // Traverse the value using the provided field extractor
}

// StringV represents a valid JSON string.
type StringV string

// Get implements the Value interface by decoding the underlying value to either a StringV or a string type.
func (str StringV) Get(i interface{}) error { return newValueDecoder(i).assign(str) }

// At implements the Value interface by returning an invalid field since StringV is not traversable.
func (str StringV) At(field Field) FieldValue { return field.get(str) }

// LongV represents a valid JSON number.
type LongV int64

// Get implements the Value interface by decoding the underlying value to either a LongV or a numeric type.
func (num LongV) Get(i interface{}) error { return newValueDecoder(i).assign(num) }

// At implements the Value interface by returning an invalid field since LongV is not traversable.
func (num LongV) At(field Field) FieldValue { return field.get(num) }

// DoubleV represents a valid JSON double.
type DoubleV float64

// Get implements the Value interface by decoding the underlying value to either a DoubleV or a float type.
func (num DoubleV) Get(i interface{}) error { return newValueDecoder(i).assign(num) }

// At implements the Value interface by returning an invalid field since DoubleV is not traversable.
func (num DoubleV) At(field Field) FieldValue { return field.get(num) }

// BooleanV represents a valid JSON boolean.
type BooleanV bool

// Get implements the Value interface by decoding the underlying value to either a BooleanV or a boolean type.
func (boolean BooleanV) Get(i interface{}) error { return newValueDecoder(i).assign(boolean) }

// At implements the Value interface by returning an invalid field since BooleanV is not traversable.
func (boolean BooleanV) At(field Field) FieldValue { return field.get(boolean) }

// DateV represents a FaunaDB date type.
type DateV time.Time

// Get implements the Value interface by decoding the underlying value to either a DateV or a time.Time type.
func (date DateV) Get(i interface{}) error { return newValueDecoder(i).assign(date) }

// At implements the Value interface by returning an invalid field since DateV is not traversable.
func (date DateV) At(field Field) FieldValue { return field.get(date) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB date representation.
func (date DateV) MarshalJSON() ([]byte, error) {
	return escape("@date", time.Time(date).Format("2006-01-02"))
}

// TimeV represents a FaunaDB time type.
type TimeV time.Time

// Get implements the Value interface by decoding the underlying value to either a TimeV or a time.Time type.
func (localTime TimeV) Get(i interface{}) error { return newValueDecoder(i).assign(localTime) }

// At implements the Value interface by returning an invalid field since TimeV is not traversable.
func (localTime TimeV) At(field Field) FieldValue { return field.get(localTime) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB time representation.
func (localTime TimeV) MarshalJSON() ([]byte, error) {
	return escape("@ts", time.Time(localTime).Format("2006-01-02T15:04:05.999999999Z"))
}

// RefV represents a FaunaDB ref type.
type RefV struct {
	ID         string
	Collection *RefV
	Class      *RefV //Deprecated: As of 2.7 Class is deprecated, use Collection instead
	Database   *RefV
}

// Get implements the Value interface by decoding the underlying ref to a RefV.
func (ref RefV) Get(i interface{}) error { return newValueDecoder(i).assign(ref) }

// At implements the Value interface by returning an invalid field since RefV is not traversable.
func (ref RefV) At(field Field) FieldValue { return field.get(ref) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB ref representation.
func (ref RefV) MarshalJSON() ([]byte, error) {
	values := map[string]interface{}{"id": ref.ID}

	if ref.Collection != nil {
		values["collection"] = ref.Collection
	}

	if ref.Database != nil {
		values["database"] = ref.Database
	}

	return escape("@ref", values)
}

var (
	nativeClasses     = RefV{"classes", nil, nil, nil}
	nativeCollections = RefV{"collections", nil, nil, nil}
	nativeIndexes     = RefV{"indexes", nil, nil, nil}
	nativeDatabases   = RefV{"databases", nil, nil, nil}
	nativeFunctions   = RefV{"functions", nil, nil, nil}
	nativeRoles       = RefV{"roles", nil, nil, nil}
	nativeKeys        = RefV{"keys", nil, nil, nil}
	nativeTokens      = RefV{"tokens", nil, nil, nil}
	nativeCredentials = RefV{"credentials", nil, nil, nil}
)

func NativeClasses() *RefV     { return &nativeClasses }
func NativeCollections() *RefV { return &nativeCollections }
func NativeIndexes() *RefV     { return &nativeIndexes }
func NativeDatabases() *RefV   { return &nativeDatabases }
func NativeFunctions() *RefV   { return &nativeFunctions }
func NativeRoles() *RefV       { return &nativeRoles }
func NativeKeys() *RefV        { return &nativeKeys }
func NativeTokens() *RefV      { return &nativeTokens }
func NativeCredentials() *RefV { return &nativeCredentials }

func nativeFromName(id string) *RefV {
	switch id {
	case "collections":
		return &nativeCollections
	case "classes":
		return &nativeClasses
	case "indexes":
		return &nativeIndexes
	case "databases":
		return &nativeDatabases
	case "functions":
		return &nativeFunctions
	case "roles":
		return &nativeRoles
	case "keys":
		return &nativeKeys
	case "tokens":
		return &nativeTokens
	case "credentials":
		return &nativeCredentials
	}

	return &RefV{id, nil, nil, nil}
}

// SetRefV represents a FaunaDB setref type.
type SetRefV struct {
	Parameters map[string]Value
}

// Get implements the Value interface by decoding the underlying value to a SetRefV.
func (set SetRefV) Get(i interface{}) error { return newValueDecoder(i).assign(set) }

// At implements the Value interface by returning an invalid field since SetRefV is not traversable.
func (set SetRefV) At(field Field) FieldValue { return field.get(set) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB setref representation.
func (set SetRefV) MarshalJSON() ([]byte, error) { return escape("@set", set.Parameters) }

// ObjectV represents a FaunaDB object type.
type ObjectV map[string]Value

// Get implements the Value interface by decoding the underlying value to either a ObjectV or a native map type.
func (obj ObjectV) Get(i interface{}) error { return newValueDecoder(i).decodeMap(obj) }

// At implements the Value interface by traversing the object and extracting the provided field.
func (obj ObjectV) At(field Field) FieldValue { return field.get(obj) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB object representation.
func (obj ObjectV) MarshalJSON() ([]byte, error) { return escape("object", map[string]Value(obj)) }

// ArrayV represents a FaunaDB array type.
type ArrayV []Value

// Get implements the Value interface by decoding the underlying value to either an ArrayV or a native slice type.
func (arr ArrayV) Get(i interface{}) error { return newValueDecoder(i).decodeArray(arr) }

// At implements the Value interface by traversing the array and extracting the provided field.
func (arr ArrayV) At(field Field) FieldValue { return field.get(arr) }

// NullV represents a valid JSON null.
type NullV struct{}

// Get implements the Value interface by decoding the underlying value to a either a NullV or a nil pointer.
func (null NullV) Get(i interface{}) error { return nil }

// At implements the Value interface by returning an invalid field since NullV is not traversable.
func (null NullV) At(field Field) FieldValue { return field.get(null) }

// MarshalJSON implements json.Marshaler by escaping its value according to JSON null representation.
func (null NullV) MarshalJSON() ([]byte, error) { return []byte("null"), nil }

// BytesV represents a FaunaDB binary blob type.
type BytesV []byte

// Get implements the Value interface by decoding the underlying value to either a ByteV or a []byte type.
func (bytes BytesV) Get(i interface{}) error { return newValueDecoder(i).assign(bytes) }

// At implements the Value interface by returning an invalid field since BytesV is not traversable.
func (bytes BytesV) At(field Field) FieldValue { return field.get(bytes) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB bytes representation.
func (bytes BytesV) MarshalJSON() ([]byte, error) {
	encoded := base64.StdEncoding.EncodeToString(bytes)
	return escape("@bytes", encoded)
}

// QueryV represents a `@query` value in FaunaDB.
type QueryV struct {
	lambda json.RawMessage
}

// Get implements the Value interface by decoding the underlying value to a QueryV.
func (query QueryV) Get(i interface{}) error { return newValueDecoder(i).assign(query) }

// At implements the Value interface by returning an invalid field since QueryV is not traversable.
func (query QueryV) At(field Field) FieldValue { return field.get(query) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB query representation.
func (query QueryV) MarshalJSON() ([]byte, error) { return escape("@query", &query.lambda) }

// Implement Expr for all values

func (str StringV) expr()      {}
func (num LongV) expr()        {}
func (num DoubleV) expr()      {}
func (boolean BooleanV) expr() {}
func (date DateV) expr()       {}
func (localTime TimeV) expr()  {}
func (ref RefV) expr()         {}
func (set SetRefV) expr()      {}
func (obj ObjectV) expr()      {}
func (arr ArrayV) expr()       {}
func (null NullV) expr()       {}
func (bytes BytesV) expr()     {}
func (query QueryV) expr()     {}

func escape(key string, value interface{}) ([]byte, error) {
	return json.Marshal(map[string]interface{}{key: value})
}
