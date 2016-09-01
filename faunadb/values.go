package faunadb

import "time"

type Value interface {
	Expr
	Get(interface{}) error
	At(Field) FieldValue
}

type StringV string

func (str StringV) Get(i interface{}) error      { return newValueDecoder(i).assign(str) }
func (str StringV) At(field Field) FieldValue    { return field.get(str) }
func (str StringV) toJSON() (interface{}, error) { return str, nil }

type LongV int64

func (num LongV) Get(i interface{}) error      { return newValueDecoder(i).assign(num) }
func (num LongV) At(field Field) FieldValue    { return field.get(num) }
func (num LongV) toJSON() (interface{}, error) { return num, nil }

type DoubleV float64

func (num DoubleV) Get(i interface{}) error      { return newValueDecoder(i).assign(num) }
func (num DoubleV) At(field Field) FieldValue    { return field.get(num) }
func (num DoubleV) toJSON() (interface{}, error) { return num, nil }

type BooleanV bool

func (boolean BooleanV) Get(i interface{}) error      { return newValueDecoder(i).assign(boolean) }
func (boolean BooleanV) At(field Field) FieldValue    { return field.get(boolean) }
func (boolean BooleanV) toJSON() (interface{}, error) { return boolean, nil }

type DateV time.Time

func (date DateV) Get(i interface{}) error      { return newValueDecoder(i).assign(date) }
func (date DateV) At(field Field) FieldValue    { return field.get(date) }
func (date DateV) toJSON() (interface{}, error) { return map[string]interface{}{"@date": date}, nil }

type TimeV time.Time

func (locaTime TimeV) Get(i interface{}) error   { return newValueDecoder(i).assign(locaTime) }
func (locaTime TimeV) At(field Field) FieldValue { return field.get(locaTime) }

func (locaTime TimeV) toJSON() (interface{}, error) {
	return map[string]interface{}{"@ts": locaTime}, nil
}

type RefV struct {
	ID string
}

func (ref RefV) Get(i interface{}) error      { return newValueDecoder(i).assign(ref) }
func (ref RefV) At(field Field) FieldValue    { return field.get(ref) }
func (ref RefV) toJSON() (interface{}, error) { return map[string]interface{}{"@ref": ref.ID}, nil }

type SetRefV struct {
	Parameters ObjectV
}

func (set SetRefV) Get(i interface{}) error   { return newValueDecoder(i).assign(set) }
func (set SetRefV) At(field Field) FieldValue { return field.get(set) }

func (set SetRefV) toJSON() (interface{}, error) {
	return map[string]interface{}{"@set": set.Parameters}, nil
}

type ObjectV map[string]Value

func (obj ObjectV) Get(i interface{}) error      { return newValueDecoder(i).decodeMap(obj) }
func (obj ObjectV) At(field Field) FieldValue    { return field.get(obj) }
func (obj ObjectV) toJSON() (interface{}, error) { return map[string]interface{}{"object": obj}, nil }

type ArrayV []Value

func (arr ArrayV) Get(i interface{}) error      { return newValueDecoder(i).decodeArray(arr) }
func (arr ArrayV) At(field Field) FieldValue    { return field.get(arr) }
func (arr ArrayV) toJSON() (interface{}, error) { return arr, nil }
