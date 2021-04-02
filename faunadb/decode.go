package faunadb

import (
	"fmt"
	"reflect"
)

// A DecodeError describes an error when decoding a Fauna Value to a native Golang type
type DecodeError struct {
	path path
	err  error
}

func (d DecodeError) Error() string {
	path := d.path
	err := d.err

	for {
		if decodeErr, ok := err.(DecodeError); ok {
			path = path.subPath(decodeErr.path)
			err = decodeErr.err
		} else {
			break
		}
	}

	return fmt.Sprintf("Error while decoding fauna value at: %s. %s", path, err)
}

type valueDecoder struct {
	target     reflect.Value
	targetType reflect.Type
}

func newValueDecoder(source interface{}) *valueDecoder {
	target, targetType := indirectValue(source)

	return &valueDecoder{
		target:     target,
		targetType: targetType,
	}
}

func (c *valueDecoder) assign(value interface{}) error {
	source, sourceType := indirectValue(value)

	if sourceType.AssignableTo(c.targetType) {
		c.target.Set(source)
		return nil
	}

	if sourceType.ConvertibleTo(c.targetType) {
		c.target.Set(source.Convert(c.targetType))
		return nil
	}

	return DecodeError{
		err: fmt.Errorf("Can not assign value of type \"%s\" to a value of type \"%s\"", sourceType, c.targetType),
	}
}

func (c *valueDecoder) decodeArray(arr ArrayV) error {
	if err := c.assign(arr); err == nil {
		return nil
	}

	if c.target.Kind() != reflect.Slice {
		return DecodeError{err: fmt.Errorf("Can not decode array into a value of type \"%s\"", c.targetType)}
	}

	return c.makeNewSlice(arr)
}

func (c *valueDecoder) makeNewSlice(arr []Value) error {
	newArray := reflect.MakeSlice(c.targetType, len(arr), len(arr))

	for index, value := range arr {
		if err := value.Get(newArray.Index(index)); err != nil {
			return DecodeError{path: pathFromIndexes(index), err: err}
		}
	}

	return c.assign(newArray)
}

func (c *valueDecoder) decodeMap(obj ObjectV) error {
	if err := c.assign(obj); err == nil {
		return nil
	}

	switch c.target.Kind() {
	case reflect.Map:
		return c.makeNewMap(obj)
	case reflect.Struct:
		return c.fillStructFields(obj)
	default:
		return DecodeError{err: fmt.Errorf("Can not decode map into a value of type \"%s\"", c.targetType)}
	}
}

func (c *valueDecoder) makeNewMap(obj map[string]Value) error {
	newMap := reflect.MakeMap(c.targetType)
	elemType := c.targetType.Elem()

	for key, value := range obj {
		newElem := reflect.New(elemType).Elem()

		if err := value.Get(newElem); err != nil {
			return DecodeError{path: pathFromKeys(key), err: err}
		}

		newMap.SetMapIndex(reflect.ValueOf(key), newElem)
	}

	return c.assign(newMap)
}

func (c *valueDecoder) fillStructFields(obj map[string]Value) (err error) {
	newStruct := reflect.New(c.targetType).Elem()
	aStructType := newStruct.Type()

	//need to get all StructFields and unpack json values into the fields to check if they are empty - use allStructFields rather than
	for key, field := range allStructFields(newStruct) {
		f, _ := aStructType.FieldByName(key)
		if !field.CanInterface() {
			continue
		}
		fieldName, ignored, _, _ := parseTag(f) //need to parse tags

		value, found := obj[fieldName]

		if !found {
			continue
		}

		if ignored {
			continue
		}

		if err = value.Get(field); err != nil {
			return DecodeError{path: pathFromKeys(key), err: err}
		}
	}
	return c.assign(newStruct)
}
