package faunadb

type Field struct {
	path path
}

type FieldValue interface {
	GetValue() (Value, error)
	Get(i interface{}) error
}

func ObjKey(keys ...string) Field {
	return Field{pathFromKeys(keys...)}
}

func ArrIndex(indexes ...int) Field {
	return Field{pathFromIndexes(indexes...)}
}

func (f Field) At(other Field) Field {
	return Field{f.path.subPath(other.path)}
}

func (f Field) AtKey(keys ...string) Field {
	return f.At(ObjKey(keys...))
}

func (f Field) AtIndex(indexes ...int) Field {
	return f.At(ArrIndex(indexes...))
}

func (f *Field) get(value Value) FieldValue {
	value, err := f.path.get(value)

	if err != nil {
		return invalidField{err}
	}

	return validField{value}
}

type validField struct {
	value Value
}

func (v validField) GetValue() (Value, error) { return v.value, nil }
func (v validField) Get(i interface{}) error  { return v.value.Get(i) }

type invalidField struct {
	err error
}

func (v invalidField) GetValue() (Value, error) { return nil, v.err }
func (v invalidField) Get(i interface{}) error  { return v.err }
