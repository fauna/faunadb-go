package faunadb

// Field is a field extractor for FaunaDB values.
type Field struct{ path path }

// FieldValue describes an extracted field value.
type FieldValue interface {
	GetValue() (Value, error) // GetValue returns the extracted FaunaDB value.
	Get(i interface{}) error  // Get decodes a FaunaDB value to a native Go type.
}

// ObjKey creates a field extractor for a JSON object based on the provided keys.
func ObjKey(keys ...string) Field { return Field{pathFromKeys(keys...)} }

// ArrIndex creates a field extractor for a JSON array based on the provided indexes.
func ArrIndex(indexes ...int) Field { return Field{pathFromIndexes(indexes...)} }

// At creates a new field extractor based on the provided path.
func (f Field) At(other Field) Field { return Field{f.path.subPath(other.path)} }

// AtKey creates a new field extractor based on the provided key.
func (f Field) AtKey(keys ...string) Field { return f.At(ObjKey(keys...)) }

// AtIndex creates a new field extractor based on the provided index.
func (f Field) AtIndex(indexes ...int) Field { return f.At(ArrIndex(indexes...)) }

func (f *Field) get(value Value) FieldValue {
	value, err := f.path.get(value)

	if err != nil {
		return invalidField{err}
	}

	return validField{value}
}

type validField struct{ value Value }

func (v validField) GetValue() (Value, error) { return v.value, nil }
func (v validField) Get(i interface{}) error  { return v.value.Get(i) }

type invalidField struct{ err error }

func (v invalidField) GetValue() (Value, error) { return nil, v.err }
func (v invalidField) Get(i interface{}) error  { return v.err }
