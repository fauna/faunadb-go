package values

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const fieldNameTag = "fauna"

func decodeValue(value *Value, i interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovered from unexpected error while decoding a fauna value: %s", r)
		}
	}()

	decoder := valueTreeDecoder{tree: value}
	return decoder.decode(reflect.ValueOf(i))
}

type reflectionDecoder interface {
	decode(reflect.Value) error
}

type valueTreeDecoder struct {
	indirectReference
	path path
	tree *Value
}

func (d *valueTreeDecoder) decode(v reflect.Value) error {
	if d.tree.inner == nil {
		return nil
	}

	var decoder reflectionDecoder

	pointer := d.indirect(v)

	switch pointer.Kind() {
	case reflect.Struct:
		decoder = d.decodeSpecialStruct(pointer.Type())
	case reflect.Map:
		decoder = objectDecoder{objectContainer{d.path, d.tree}}
	case reflect.Slice:
		decoder = sliceDecoder{arrayContainer{d.path, d.tree}}
	default:
		decoder = assignDecoder{path: d.path, source: d.indirectValue(d.tree.inner)}
	}

	return decoder.decode(pointer)
}

func (d *valueTreeDecoder) decodeSpecialStruct(pointerType reflect.Type) reflectionDecoder {
	if source := d.indirectValue(d.tree); source.Type().AssignableTo(pointerType) {
		return assignDecoder{path: d.path, source: source}
	}

	if source := d.indirectValue(d.tree.inner); source.Type().AssignableTo(pointerType) {
		return assignDecoder{path: d.path, source: source}
	}

	return structDecoder{objectContainer{path: d.path, value: d.tree}}
}

type objectContainer struct {
	path  path
	value *Value
}

func (c *objectContainer) toMap() (map[string]Value, error) {
	if obj, ok := c.value.inner.(map[string]Value); ok {
		return obj, nil
	}

	return nil, newDecodeError(c.path, "Expected internal value to be a json object but was %T", c.value.inner)
}

type structDecoder struct {
	objectContainer
}

func (d structDecoder) decode(aStruct reflect.Value) error {
	object, err := d.toMap()
	if err != nil {
		return err
	}

	structType := aStruct.Type()

	for i, numFields := 0, aStruct.NumField(); i < numFields; i++ {
		key, value, found := d.getKeyPair(object, structType.Field(i))

		if !found {
			continue
		}

		fieldValue := aStruct.Field(i)
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			continue
		}

		decoder := valueTreeDecoder{
			path: d.path.subpathKey(key),
			tree: &value,
		}

		if err := decoder.decode(fieldValue); err != nil {
			return err
		}
	}

	return nil
}

func (d *structDecoder) getKeyPair(object map[string]Value, field reflect.StructField) (key string, value Value, ok bool) {
	if key = field.Tag.Get(fieldNameTag); key == "" {
		key = field.Name
	}

	value, ok = object[key]
	return
}

type objectDecoder struct {
	objectContainer
}

func (d objectDecoder) decode(mapValue reflect.Value) error {
	object, err := d.toMap()
	if err != nil {
		return err
	}

	elemType := mapValue.Type().Elem()
	newMap := reflect.MakeMap(mapValue.Type())

	for key, value := range object {
		newElem := reflect.Indirect(reflect.New(elemType))

		decoder := valueTreeDecoder{
			path: d.path.subpathKey(key),
			tree: &value,
		}

		if err := decoder.decode(newElem); err != nil {
			return err
		}

		newMap.SetMapIndex(reflect.ValueOf(key), newElem)
	}

	mapValue.Set(newMap)
	return nil
}

type arrayContainer struct {
	path  path
	value *Value
}

func (c *arrayContainer) toArray() ([]Value, error) {
	if arr, ok := c.value.inner.([]Value); ok {
		return arr, nil
	}

	return nil, newDecodeError(c.path, "Expected internal value to be a json array but was %T", c.value.inner)
}

type sliceDecoder struct {
	arrayContainer
}

func (d sliceDecoder) decode(sliceValue reflect.Value) error {
	array, err := d.toArray()
	if err != nil {
		return err
	}

	newSlice := reflect.MakeSlice(sliceValue.Type(), len(array), len(array))

	for index, value := range array {
		decoder := valueTreeDecoder{
			path: d.path.subpathIndex(index),
			tree: &value,
		}

		if err := decoder.decode(newSlice.Index(index)); err != nil {
			return err
		}
	}

	sliceValue.Set(newSlice)
	return nil
}

type assignDecoder struct {
	indirectReference
	path   path
	source reflect.Value
}

func (d assignDecoder) decode(pointer reflect.Value) error {
	if !pointer.CanSet() {
		return newDecodeError(d.path, "Target reference is not assignable")
	}

	sourceType := d.source.Type()
	pointerType := pointer.Type()

	if sourceType.ConvertibleTo(pointerType) {
		pointer.Set(d.source.Convert(pointerType))
		return nil
	}

	if sourceType.AssignableTo(pointerType) {
		pointer.Set(d.source)
		return nil
	}

	return newDecodeError(d.path, "Can not assign value of type %s to a value of type %s", sourceType, pointerType)
}

type path struct {
	segments []string
}

func (p *path) String() string {
	str := strings.Join(p.segments, " / ")

	if str == "" {
		str = "root"
	}

	return str
}

func (p *path) subpathKey(other string) path {
	return path{segments: append(p.segments, other)}
}

func (p *path) subpathIndex(index int) path {
	return path{segments: append(p.segments, strconv.Itoa(index))}
}

type decodeError struct {
	path path
	err  error
}

func newDecodeError(path path, message string, args ...interface{}) decodeError {
	return decodeError{path, fmt.Errorf(message, args...)}
}

func (e decodeError) Error() string {
	return fmt.Sprintf("Error while decoding fauna value at: %s. %s", e.path.String(), e.err)
}

// Reflect on the references trying to reach to an addressable value
type indirectReference struct{}

func (d *indirectReference) indirectValue(i interface{}) reflect.Value {
	return d.indirect(reflect.ValueOf(i))
}

// FIXME: I think we can replace this method by reflect.Indirect
func (d *indirectReference) indirect(source reflect.Value) (value reflect.Value) {
	value = source

	for {
		if value.Kind() == reflect.Interface {
			value = value.Elem()
			continue
		}

		if value.Kind() != reflect.Ptr {
			break
		}

		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}

		value = value.Elem()
	}

	return
}
