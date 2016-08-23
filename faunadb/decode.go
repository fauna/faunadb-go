package faunadb

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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

	if !c.target.CanSet() {
		return nil // Don't attempt to set values on unexported variables/pointers
	}

	if sourceType.AssignableTo(c.targetType) {
		c.target.Set(source)
		return nil
	}

	if sourceType.ConvertibleTo(c.targetType) {
		c.target.Set(source.Convert(c.targetType))
		return nil
	}

	return decodeError{
		err: fmt.Errorf("Can not assign value of type \"%s\" to a value of type \"%s\"", sourceType, c.targetType),
	}
}

func (c *valueDecoder) decodeArray(arr ArrayV) error {
	if err := c.assign(arr); err == nil {
		return nil
	}

	if c.target.Kind() != reflect.Slice {
		return decodeError{err: fmt.Errorf("Can not decode array into a value of type \"%s\"", c.targetType)}
	}

	return c.makeNewSlice(arr)
}

func (c *valueDecoder) makeNewSlice(arr []Value) error {
	newArray := reflect.MakeSlice(c.targetType, len(arr), len(arr))

	for index, value := range arr {
		if err := value.To(newArray.Index(index)); err != nil {
			return decodeError{path: strconv.Itoa(index), err: err}
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
		return decodeError{err: fmt.Errorf("Can not decode map into a value of type \"%s\"", c.targetType)}
	}
}

func (c *valueDecoder) makeNewMap(obj map[string]Value) error {
	newMap := reflect.MakeMap(c.targetType)
	elemType := c.targetType.Elem()

	for key, value := range obj {
		newElem := reflect.New(elemType).Elem()

		if err := value.To(newElem); err != nil {
			return decodeError{path: key, err: err}
		}

		newMap.SetMapIndex(reflect.ValueOf(key), newElem)
	}

	return c.assign(newMap)
}

func (c *valueDecoder) fillStructFields(obj map[string]Value) error {
	newStruct := reflect.New(c.targetType).Elem()

	for key, field := range exportedStructFields(newStruct) {
		value, found := obj[key]
		if !found {
			continue
		}

		if err := value.To(field); err != nil {
			return decodeError{path: key, err: err}
		}
	}

	return c.assign(newStruct)
}

type decodeError struct {
	path string
	err  error
}

func (d decodeError) Error() string {
	var segments []string

	path := d.path
	err := d.err

	for {
		if path != "" {
			segments = append(segments, path)
		}

		if next, ok := err.(decodeError); ok {
			path = next.path
			err = next.err
		} else {
			break
		}
	}

	path = strings.Join(segments, " / ")

	if path == "" {
		path = "root"
	}

	return fmt.Sprintf("Error while decoding fauna value at: %s. %s", path, err)
}
