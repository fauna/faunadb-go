package faunadb

import "time"

type Value interface {
	Expr
	To(interface{}) error
}

type StringV string

func (str StringV) To(i interface{}) error {
	return newValueDecoder(i).assign(str)
}

func (str StringV) toJSON() (interface{}, error) {
	return str, nil
}

type LongV int64

func (num LongV) To(i interface{}) error {
	return newValueDecoder(i).assign(num)
}

func (num LongV) toJSON() (interface{}, error) {
	return num, nil
}

type DoubleV float64

func (num DoubleV) To(i interface{}) error {
	return newValueDecoder(i).assign(num)
}

func (num DoubleV) toJSON() (interface{}, error) {
	return num, nil
}

type BooleanV bool

func (boolean BooleanV) To(i interface{}) error {
	return newValueDecoder(i).assign(boolean)
}

func (boolean BooleanV) toJSON() (interface{}, error) {
	return boolean, nil
}

type DateV time.Time

func (date DateV) To(i interface{}) error {
	return newValueDecoder(i).assign(date)
}

func (date DateV) toJSON() (interface{}, error) {
	return map[string]interface{}{"@date": date}, nil
}

type TimeV time.Time

func (locaTime TimeV) To(i interface{}) error {
	return newValueDecoder(i).assign(locaTime)
}

func (locaTime TimeV) toJSON() (interface{}, error) {
	return map[string]interface{}{"@ts": locaTime}, nil
}

type RefV struct {
	ID string
}

func (ref RefV) To(i interface{}) error {
	return newValueDecoder(i).assign(ref)
}

func (ref RefV) toJSON() (interface{}, error) {
	return map[string]interface{}{"@ref": ref.ID}, nil
}

type SetRefV struct {
	Parameters ObjectV
}

func (set SetRefV) To(i interface{}) error {
	return newValueDecoder(i).assign(set)
}

func (set SetRefV) toJSON() (interface{}, error) {
	return map[string]interface{}{"@set": set.Parameters}, nil
}

type ObjectV map[string]Value

func (obj ObjectV) To(i interface{}) error {
	return newValueDecoder(i).decodeMap(obj)
}

func (obj ObjectV) toJSON() (interface{}, error) {
	return map[string]interface{}{"object": obj}, nil
}

type ArrayV []Value

func (arr ArrayV) To(i interface{}) error {
	return newValueDecoder(i).decodeArray(arr)
}

func (arr ArrayV) toJSON() (interface{}, error) {
	return arr, nil
}
