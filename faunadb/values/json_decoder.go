package values

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	startOfObject = json.Delim('{')
	startOfArray  = json.Delim('[')
	endOfObject   = json.Delim('}')
	endOfArray    = json.Delim(']')
)

func ReadValue(reader io.Reader) (Value, error) {
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()

	jsonReader := jsonReader{&scan{decoder: decoder}}
	value, _ := jsonReader.read()

	if jsonReader.scan.err != nil {
		return Value{}, jsonReader.scan.err
	}

	return value, nil
}

type valueReader interface {
	read() (Value, bool)
}

type jsonReader struct {
	scan *scan
}

func (reader *jsonReader) read() (Value, bool) {
	return reader.scan.readNext()
}

type scan struct {
	decoder *json.Decoder
	err     error
}

func (s *scan) readNext() (Value, bool) {
	next, ok := s.next()

	if !ok {
		return Value{}, false
	}

	return next.read()
}

func (s *scan) next() (reader valueReader, _ bool) {
	token, err := s.decoder.Token()

	if err != nil {
		s.err = err
		return nil, false
	}

	switch token {
	default:
		reader = literalReader{s, token}
	case startOfArray:
		reader = arrayReader{s}
	case startOfObject:
		reader = s.readSpecialObject()
	case endOfArray, endOfObject:
		return nil, false
	}

	return reader, true
}

func (s *scan) readSpecialObject() valueReader {
	firstKey, ok := s.readString()

	if ok {
		switch firstKey {
		case "@ref":
			return refReader{s}
		case "@obj":
			return literalObjectReader{s}
		case "@set":
			return setRefReader{s}
		case "@date":
			return dateTimeReader{dateConverter{}, s, "2006-01-02"}
		case "@ts":
			return dateTimeReader{timeConverter{}, s, "2006-01-02T15:04:05.999999999Z"}
		}
	}

	return objectReader{s, firstKey}
}

func (s *scan) readString() (str string, ok bool) {
	var value Value

	if value, ok = s.readNext(); ok {
		if str, ok = value.inner.(string); !ok {
			s.err = fmt.Errorf("Expected string but got %T", value.inner)
		}
	}

	return
}

func (s *scan) ensureNoMoreTokens() bool {
	_, hasMore := s.next()

	if hasMore {
		s.err = errors.New("JSON type is bigger than expected")
	}

	return !hasMore
}

type objectReader struct {
	scan     *scan
	firstKey string
}

func (reader objectReader) read() (Value, bool) {
	object := make(map[string]Value)

	if key := reader.firstKey; key != "" {
		for {
			value, ok := reader.scan.readNext()
			if !ok {
				break
			}

			object[key] = value

			key, ok = reader.scan.readString()
			if !ok {
				break
			}
		}
	}

	return Value{object}, true
}

type arrayReader struct {
	scan *scan
}

func (reader arrayReader) read() (Value, bool) {
	var array []Value

	for {
		value, ok := reader.scan.readNext()
		if !ok {
			break
		}
		array = append(array, value)
	}

	return Value{array}, true
}

type literalReader struct {
	scan  *scan
	token json.Token
}

func (reader literalReader) read() (Value, bool) {
	if number, ok := reader.token.(json.Number); ok {
		return reader.parseJsonNumber(number)
	} else {
		return Value{reader.token}, true
	}
}

func (reader *literalReader) parseJsonNumber(number json.Number) (Value, bool) {
	var err error

	if strings.Contains(number.String(), ".") {
		var n float64
		if n, err = number.Float64(); err == nil {
			return Value{n}, true
		}
	} else {
		var n int64
		if n, err = number.Int64(); err == nil {
			return Value{n}, true
		}
	}

	reader.scan.err = err
	return Value{}, false
}

type refReader struct {
	scan *scan
}

func (reader refReader) read() (Value, bool) {
	if id, ok := reader.scan.readString(); ok && reader.scan.ensureNoMoreTokens() {
		return Value{RefV{id}}, true
	}

	return Value{}, false
}

type literalObjectReader struct {
	scan *scan
}

func (reader literalObjectReader) read() (Value, bool) {
	if obj, ok := reader.scan.readNext(); ok && reader.scan.ensureNoMoreTokens() {
		return obj, true
	}
	return Value{}, false
}

type setRefReader struct {
	scan *scan
}

func (reader setRefReader) read() (Value, bool) {
	if v, ok := reader.scan.readNext(); ok && reader.scan.ensureNoMoreTokens() {
		return Value{SetRefV{v.inner.(map[string]Value)}}, true
	}

	return Value{}, false
}

type timeToValue interface {
	toValue(time.Time) Value
}

type dateTimeReader struct {
	timeToValue
	scan   *scan
	format string
}

func (reader dateTimeReader) read() (Value, bool) {
	if str, ok := reader.scan.readString(); ok && reader.scan.ensureNoMoreTokens() {
		return reader.parseTime(str)
	}

	return Value{}, false
}

func (reader *dateTimeReader) parseTime(raw string) (value Value, ok bool) {
	if t, err := time.Parse(reader.format, raw); err == nil {
		value, ok = reader.toValue(t), true
	} else {
		reader.scan.err = err
	}

	return
}

type dateConverter struct{}

func (reader dateConverter) toValue(t time.Time) Value {
	return Value{DateV{t}}
}

type timeConverter struct{}

func (reader timeConverter) toValue(t time.Time) Value {
	return Value{TimeV{t}}
}
