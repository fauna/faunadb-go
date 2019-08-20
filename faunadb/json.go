package faunadb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// Unmarshal json string into a FaunaDB value.
func UnmarshalJSON(buffer []byte, outValue *Value) error {
	reader := bytes.NewReader(buffer)

	if value, err := parseJSON(reader); err != nil {
		return err
	} else {
		*outValue = value
	}

	return nil
}

// Marshal a FaunaDB value into a json string.
func MarshalJSON(value Value) ([]byte, error) {
	return json.Marshal(unwrap(value))
}

func unwrap(value Value) interface{} {
	switch v := value.(type) {
	case ArrayV:
		array := make([]interface{}, len(v))

		for i, elem := range v {
			array[i] = unwrap(elem)
		}

		return array

	case ObjectV:
		object := make(map[string]interface{}, len(v))

		for key, elem := range v {
			object[key] = unwrap(elem)
		}

		return object

	default:
		return value
	}
}

func parseJSON(reader io.Reader) (Value, error) {
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()

	parser := jsonParser{decoder}
	return parser.parseNext()
}

type wrongToken struct {
	expected string
	got      json.Token
}

func (w wrongToken) Error() string {
	return fmt.Sprintf("Expected %s but got %#v", w.expected, w.got)
}

type jsonParser struct {
	decoder *json.Decoder
}

func (p *jsonParser) parseNext() (Value, error) {
	token, err := p.decoder.Token()
	if err != nil {
		return nil, err
	}

	switch token {
	case json.Delim('{'):
		return p.parseSpecialObject()
	case json.Delim('['):
		return p.parseArray()
	default:
		return p.parseLiteral(token)
	}
}

func (p *jsonParser) parseLiteral(token json.Token) (value Value, err error) {
	switch v := token.(type) {
	case string:
		value = StringV(v)
	case bool:
		value = BooleanV(v)
	case json.Number:
		value, err = p.parseJSONNumber(v)
	case nil:
		value = NullV{}
	default:
		err = wrongToken{"a literal", v}
	}

	return
}

func (p *jsonParser) parseSpecialObject() (value Value, err error) {
	if !p.hasMore() {
		value = ObjectV{}
		return
	}

	var firstKey string

	if firstKey, err = p.readString(); err == nil {
		switch firstKey {
		case "@ref":
			value, err = p.parseRef()
		case "@set":
			value, err = p.parseSet()
		case "@date":
			value, err = p.parseDate("2006-01-02", func(t time.Time) Value { return DateV(t) })
		case "@ts":
			value, err = p.parseDate("2006-01-02T15:04:05.999999999Z", func(t time.Time) Value { return TimeV(t) })
		case "@obj":
			value, err = p.readSingleObject()
		case "@bytes":
			value, err = p.parseBytes()
		case "@query":
			value, err = p.parseQuery()
		default:
			value, err = p.parseObject(firstKey)
		}
	}

	return
}

func (p *jsonParser) parseRef() (value Value, err error) {
	var obj ObjectV

	if obj, err = p.readSingleObject(); err == nil {
		var id string
		var col, db *RefV

		if v, ok := obj["id"]; ok {
			if err = v.Get(&id); err != nil {
				return
			}
		}

		if v, ok := obj["collection"]; ok {
			if err = v.Get(&col); err != nil {
				return
			}
		}

		if v, ok := obj["database"]; ok {
			if err = v.Get(&db); err != nil {
				return
			}
		}

		if col == nil && db == nil {
			value = nativeFromName(id)
		} else {
			value = RefV{id, col, col, db}
		}
	}

	return
}

func (p *jsonParser) parseSet() (value Value, err error) {
	var obj ObjectV

	if obj, err = p.readSingleObject(); err == nil {
		value = SetRefV{obj}
	}

	return
}

func (p *jsonParser) parseBytes() (value Value, err error) {
	var encoded string

	if encoded, err = p.readSingleString(); err == nil {
		bytes, err := base64.StdEncoding.DecodeString(encoded)
		if err == nil {
			value = BytesV(bytes)
		}
	}

	return
}

func (p *jsonParser) parseQuery() (value Value, err error) {
	var lambda json.RawMessage

	if err = p.decoder.Decode(&lambda); err == nil {
		value = QueryV{lambda}
	}

	var token json.Token
	if token, err = p.decoder.Token(); err == nil {
		if token != json.Delim('}') {
			err = wrongToken{"end of object", token}
		}
	}

	return
}

func (p *jsonParser) parseObject(firstKey string) (Value, error) {
	object := make(map[string]Value)

	if key := firstKey; key != "" {
		for {
			if value, err := p.parseNext(); err == nil {
				object[key] = value

				if !p.hasMore() {
					break
				}

				if key, err = p.readString(); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	return ObjectV(object), nil
}

func (p *jsonParser) parseArray() (Value, error) {
	var array []Value

	for {
		if !p.hasMore() {
			break
		}

		if value, err := p.parseNext(); err == nil {
			array = append(array, value)
		} else {
			return nil, err
		}
	}

	return ArrayV(array), nil
}

func (p *jsonParser) parseDate(format string, fn func(t time.Time) Value) (value Value, err error) {
	var str string

	if str, err = p.readSingleString(); err == nil {
		value, err = p.parseStrTime(str, format, fn)
	}

	return
}

func (p *jsonParser) parseStrTime(raw string, format string, fn func(time.Time) Value) (value Value, err error) {
	var t time.Time

	if t, err = time.Parse(format, raw); err == nil {
		value = fn(t)
	}

	return
}

func (p *jsonParser) parseJSONNumber(number json.Number) (Value, error) {
	var err error

	if strings.Contains(number.String(), ".") {
		var n float64
		if n, err = number.Float64(); err == nil {
			return DoubleV(n), nil
		}
	} else {
		var n int64
		if n, err = number.Int64(); err == nil {
			return LongV(n), nil
		}
	}

	return nil, err
}

func (p *jsonParser) readSingleString() (str string, err error) {
	if str, err = p.readString(); err == nil {
		err = p.ensureNoMoreTokens()
	}

	return
}

func (p *jsonParser) readSingleObject() (obj ObjectV, err error) {
	var value Value
	var ok bool

	if value, err = p.parseNext(); err == nil {
		if err = p.ensureNoMoreTokens(); err == nil {
			if obj, ok = value.(ObjectV); !ok {
				err = wrongToken{"a single object", value}
			}
		}
	}

	return
}

func (p *jsonParser) readString() (str string, err error) {
	var token json.Token
	var ok bool

	if token, err = p.decoder.Token(); err == nil {
		if str, ok = token.(string); !ok {
			err = wrongToken{"a string", token}
		}
	}

	return
}

func (p *jsonParser) ensureNoMoreTokens() error {
	if p.hasMore() {
		token, _ := p.decoder.Token()
		return wrongToken{"end of array or object", token}
	}

	return nil
}

func (p *jsonParser) hasMore() bool {
	if !p.decoder.More() {
		_, _ = p.decoder.Token() // Discarts next } or ] token
		return false
	}

	return true
}
