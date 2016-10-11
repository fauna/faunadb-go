package faunadb

import (
	"encoding/json"
	"time"
)

type Value interface {
	Expr
	Get(interface{}) error
	At(Field) FieldValue
}

func escape(key string, value interface{}) ([]byte, error) {
	return json.Marshal(map[string]interface{}{key: value})
}

type StringV string

func (str StringV) expr()                     {}
func (str StringV) Get(i interface{}) error   { return newValueDecoder(i).assign(str) }
func (str StringV) At(field Field) FieldValue { return field.get(str) }

type LongV int64

func (num LongV) expr()                     {}
func (num LongV) Get(i interface{}) error   { return newValueDecoder(i).assign(num) }
func (num LongV) At(field Field) FieldValue { return field.get(num) }

type DoubleV float64

func (num DoubleV) expr()                     {}
func (num DoubleV) Get(i interface{}) error   { return newValueDecoder(i).assign(num) }
func (num DoubleV) At(field Field) FieldValue { return field.get(num) }

type BooleanV bool

func (boolean BooleanV) expr()                     {}
func (boolean BooleanV) Get(i interface{}) error   { return newValueDecoder(i).assign(boolean) }
func (boolean BooleanV) At(field Field) FieldValue { return field.get(boolean) }

type DateV time.Time

func (date DateV) expr()                     {}
func (date DateV) Get(i interface{}) error   { return newValueDecoder(i).assign(date) }
func (date DateV) At(field Field) FieldValue { return field.get(date) }

func (date DateV) MarshalJSON() ([]byte, error) {
	return escape("@date", time.Time(date).Format("2006-01-02"))
}

type TimeV time.Time

func (localTime TimeV) expr()                     {}
func (localTime TimeV) Get(i interface{}) error   { return newValueDecoder(i).assign(localTime) }
func (localTime TimeV) At(field Field) FieldValue { return field.get(localTime) }

func (localTime TimeV) MarshalJSON() ([]byte, error) {
	return escape("@ts", time.Time(localTime).Format("2006-01-02T15:04:05.999999999Z"))
}

type RefV struct {
	ID string
}

func (ref RefV) expr()                        {}
func (ref RefV) Get(i interface{}) error      { return newValueDecoder(i).assign(ref) }
func (ref RefV) At(field Field) FieldValue    { return field.get(ref) }
func (ref RefV) MarshalJSON() ([]byte, error) { return escape("@ref", ref.ID) }

type SetRefV struct {
	Parameters map[string]Value
}

func (set SetRefV) expr()                        {}
func (set SetRefV) Get(i interface{}) error      { return newValueDecoder(i).assign(set) }
func (set SetRefV) At(field Field) FieldValue    { return field.get(set) }
func (set SetRefV) MarshalJSON() ([]byte, error) { return escape("@set", set.Parameters) }

type ObjectV map[string]Value

func (obj ObjectV) expr()                        {}
func (obj ObjectV) Get(i interface{}) error      { return newValueDecoder(i).decodeMap(obj) }
func (obj ObjectV) At(field Field) FieldValue    { return field.get(obj) }
func (obj ObjectV) MarshalJSON() ([]byte, error) { return escape("object", map[string]Value(obj)) }

type ArrayV []Value

func (arr ArrayV) expr()                     {}
func (arr ArrayV) Get(i interface{}) error   { return newValueDecoder(i).decodeArray(arr) }
func (arr ArrayV) At(field Field) FieldValue { return field.get(arr) }

type NullV struct{}

func (null NullV) expr()                        {}
func (null NullV) Get(i interface{}) error      { return nil }
func (null NullV) At(field Field) FieldValue    { return field.get(null) }
func (null NullV) MarshalJSON() ([]byte, error) { return []byte("null"), nil }
