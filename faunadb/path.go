package faunadb

import (
	"fmt"
	"strings"
)

type InvalidFieldType struct {
	path    path
	segment invalidSegmentType
}

func (i InvalidFieldType) Error() string {
	return fmt.Sprintf("Error while extrating path: %s. %s", i.path, i.segment)
}

type invalidSegmentType struct {
	desired string
	actual  interface{}
}

func (i invalidSegmentType) Error() string {
	return fmt.Sprintf("Expected value to be %s but was a %T", i.desired, i.actual)
}

type ValueNotFound struct {
	path    path
	segment segmentNotFound
}

func (v ValueNotFound) Error() string {
	return fmt.Sprintf("Error while extrating path: %s. %s", v.path, v.segment)
}

type segmentNotFound struct {
	desired string
	segment segment
}

func (s segmentNotFound) Error() string {
	return fmt.Sprintf("%s %v not found", s.desired, s.segment)
}

type segment interface {
	get(Value) (Value, error)
}

type path []segment

func pathFromKeys(keys ...string) path {
	var p path

	for _, key := range keys {
		p = append(p, objectSegment(key))
	}

	return p
}

func pathFromIndexes(indexes ...int) path {
	var p path

	for _, index := range indexes {
		p = append(p, arraySegment(index))
	}

	return p
}

func (p path) subPath(other path) path {
	return append(p, other...)
}

func (p path) get(value Value) (Value, error) {
	var err error

	next := value

	for _, seg := range p {
		if next, err = seg.get(next); err != nil {
			switch segErr := err.(type) {
			case segmentNotFound:
				return nil, ValueNotFound{p, segErr}
			case invalidSegmentType:
				return nil, InvalidFieldType{p, segErr}
			default:
				return nil, err
			}
		}
	}

	return next, nil
}

func (p path) String() (str string) {
	var segments []string

	for _, seg := range p {
		segments = append(segments, fmt.Sprintf("%v", seg))
	}

	str = strings.Join(segments, " / ")

	if str == "" {
		str = "<root>"
	}

	return
}

type objectSegment string

func (seg objectSegment) get(value Value) (res Value, err error) {
	key := string(seg)

	switch obj := value.(type) {
	case ObjectV:
		if value, ok := obj[key]; ok {
			res = value
		} else {
			err = segmentNotFound{"Object key", seg}
		}
	default:
		err = invalidSegmentType{"an object", value}
	}

	return
}

type arraySegment int

func (seg arraySegment) get(value Value) (res Value, err error) {
	index := int(seg)

	switch arr := value.(type) {
	case ArrayV:
		if index >= 0 && index < len(arr) {
			res = arr[index]
		} else {
			err = segmentNotFound{"Array index", seg}
		}
	default:
		err = invalidSegmentType{"an array", value}
	}

	return
}
