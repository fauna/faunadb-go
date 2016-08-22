package values

import (
	"encoding/json"
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
	return jsonReader.read()
}

type valueReader interface {
	read() (Value, error)
}

type jsonReader struct {
	scan *scan
}

func (reader *jsonReader) read() (Value, error) {
	return reader.scan.readNext()
}

type scan struct {
	decoder *json.Decoder
}

func (s *scan) readNext() (Value, error) {
	if next, err := s.next(); err == nil {
		return next.read()
	} else {
		return Value{}, err
	}
}

func (s *scan) next() (reader valueReader, err error) {
	var token json.Token

	if token, err = s.decoder.Token(); err != nil {
		return
	}

	switch token {
	default:
		reader = literalReader{s, token}
	case startOfObject:
		reader, err = s.readSpecialObject()
	case startOfArray:
		reader = arrayReader{s}
	case endOfArray, endOfObject:
		err = io.EOF
	}

	return
}

func (s *scan) readSpecialObject() (reader valueReader, err error) {
	if !s.hasMore() {
		reader = objectReader{s, ""}
		return
	}

	var firstKey string

	if firstKey, err = s.readString(); err == nil {
		switch firstKey {
		case "@ref":
			reader = refReader{s}
		case "@obj":
			reader = literalObjectReader{s}
		case "@set":
			reader = setRefReader{s}
		case "@date":
			reader = dateTimeReader{dateConverter{}, s, "2006-01-02"}
		case "@ts":
			reader = dateTimeReader{timeConverter{}, s, "2006-01-02T15:04:05.999999999Z"}
		default:
			reader = objectReader{s, firstKey}
		}
	}

	return
}

func (s *scan) readString() (str string, err error) {
	var value Value
	var ok bool

	if value, err = s.readNext(); err == nil {
		if str, ok = value.inner.(string); !ok {
			err = fmt.Errorf("Expected string but got %T", value.inner)
		}
	}

	return
}

func (s *scan) readObject() (obj map[string]Value, err error) {
	var value Value
	var ok bool

	if value, err = s.readNext(); err == nil {
		if err = s.ensureNoMoreTokens(); err == nil {
			if obj, ok = value.inner.(map[string]Value); !ok {
				err = fmt.Errorf("Expected single object but got %T", value.inner)
			}
		}
	}

	return
}

func (s *scan) ensureNoMoreTokens() error {
	if s.hasMore() {
		token, _ := s.decoder.Token()
		return fmt.Errorf("Expected end of array or object but got %s", token)
	}

	return nil
}

func (s *scan) hasMore() bool {
	if !s.decoder.More() {
		s.decoder.Token() // Discarts next } or ] token
		return false
	}

	return true
}

type objectReader struct {
	scan     *scan
	firstKey string
}

func (reader objectReader) read() (Value, error) {
	object := make(map[string]Value)

	if key := reader.firstKey; key != "" {
		for {
			if value, err := reader.scan.readNext(); err == nil {
				object[key] = value

				if !reader.scan.hasMore() {
					break
				}

				if key, err = reader.scan.readString(); err != nil {
					return Value{}, err
				}
			} else {
				return Value{}, err
			}
		}
	}

	return Value{object}, nil
}

type arrayReader struct {
	scan *scan
}

func (reader arrayReader) read() (Value, error) {
	var array []Value

	for {
		if !reader.scan.hasMore() {
			break
		}

		if value, err := reader.scan.readNext(); err == nil {
			array = append(array, value)
		} else {
			return Value{}, err
		}
	}

	return Value{array}, nil
}

type literalReader struct {
	scan  *scan
	token json.Token
}

func (reader literalReader) read() (Value, error) {
	if number, ok := reader.token.(json.Number); ok {
		return reader.parseJsonNumber(number)
	}

	return Value{reader.token}, nil
}

func (reader *literalReader) parseJsonNumber(number json.Number) (Value, error) {
	var err error

	if strings.Contains(number.String(), ".") {
		var n float64
		if n, err = number.Float64(); err == nil {
			return Value{n}, nil
		}
	} else {
		var n int64
		if n, err = number.Int64(); err == nil {
			return Value{n}, nil
		}
	}

	return Value{}, err
}

type refReader struct {
	scan *scan
}

func (reader refReader) read() (value Value, err error) {
	var id string

	if id, err = reader.scan.readString(); err == nil {
		if err = reader.scan.ensureNoMoreTokens(); err == nil {
			value = Value{RefV{id}}
		}
	}

	return
}

type literalObjectReader struct {
	scan *scan
}

func (reader literalObjectReader) read() (Value, error) {
	if obj, err := reader.scan.readObject(); err == nil {
		return Value{obj}, nil
	} else {
		return Value{}, err
	}
}

type setRefReader struct {
	scan *scan
}

func (reader setRefReader) read() (Value, error) {
	if obj, err := reader.scan.readObject(); err == nil {
		return Value{SetRefV{obj}}, nil
	} else {
		return Value{}, err
	}
}

type timeToValue interface {
	toValue(time.Time) Value
}

type dateTimeReader struct {
	timeToValue
	scan   *scan
	format string
}

func (reader dateTimeReader) read() (value Value, err error) {
	var str string

	if str, err = reader.scan.readString(); err == nil {
		if err = reader.scan.ensureNoMoreTokens(); err == nil {
			value, err = reader.parseTime(str)
		}
	}

	return
}

func (reader *dateTimeReader) parseTime(raw string) (Value, error) {
	if t, err := time.Parse(reader.format, raw); err == nil {
		return reader.toValue(t), nil
	} else {
		return Value{}, err
	}
}

type dateConverter struct{}

func (reader dateConverter) toValue(t time.Time) Value {
	return Value{DateV{t}}
}

type timeConverter struct{}

func (reader timeConverter) toValue(t time.Time) Value {
	return Value{TimeV{t}}
}
