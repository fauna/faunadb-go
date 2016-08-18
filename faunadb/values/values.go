package values

import (
	"bytes"
	"time"
)

type Value struct {
	inner interface{}
}

// TODO: This function should be removed after we refactor the library packages
// Pure values should be created anywhere
func NewValue(i interface{}) Value {
	return Value{i}
}

func (value *Value) UnmarshalJSON(b []byte) (err error) {
	var decoded Value

	if decoded, err = ReadValue(bytes.NewReader(b)); err == nil {
		value.inner = decoded.inner
	}

	return
}

func (value *Value) Get(i interface{}) error {
	return decodeValue(value, i)
}

type RefV struct {
	Id string `json:"@ref"`
}

type DateV struct {
	date time.Time
}

type TimeV struct {
	time time.Time
}

type SetRefV struct {
	parameters map[string]Value
}
