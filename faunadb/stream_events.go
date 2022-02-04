package faunadb

import (
	"errors"
	"fmt"
)

// StreamEventType is a stream eveny type
type StreamEventType = string

const (
	// ErrorEventT is the stream error event type
	ErrorEventT StreamEventType = "error"

	// HistoryRewriteEventT is the stream history rewrite event type
	HistoryRewriteEventT StreamEventType = "history_rewrite"

	// StartEventT is the stream start event type
	StartEventT StreamEventType = "start"

	// VersionEventT is the stream version event type
	VersionEventT StreamEventType = "version"

	// SetT is the stream set event type
	SetEventT StreamEventType = "set"
)

// StreamEvent represents a stream event with a `type` and `txn`
type StreamEvent interface {
	Type() StreamEventType
	Txn() int64
	String() string
}

// StartEvent emitted when a valid stream subscription begins.
// Upcoming events are guaranteed to have transaction timestamps equal to or greater than
// the stream's start timestamp.
type StartEvent struct {
	StreamEvent
	txn   int64
	event Value
}

// Type returns the stream event type
func (event StartEvent) Type() StreamEventType {
	return StartEventT
}

// Txn returns the stream event timestamp
func (event StartEvent) Txn() int64 {
	return event.txn
}

// Event returns the stream event as a `f.Value`
func (event StartEvent) Event() Value {
	return event.event
}

func (event StartEvent) String() string {
	return fmt.Sprintf("StartEvent{event=%d, txn=%d} ", event.Event(), event.Txn())
}

// VersionEvent represents a version event that occurs upon any
// modifications to the current state of the subscribed document.
type VersionEvent struct {
	StreamEvent
	txn   int64
	event Value
}

// Txn returns the stream event timestamp
func (event VersionEvent) Txn() int64 {
	return event.txn
}

// Event returns the stream event as a `Value`
func (event VersionEvent) Event() Value {
	return event.event
}

func (event VersionEvent) String() string {
	return fmt.Sprintf("VersionEvent{txn=%d, event=%s}", event.Txn(), event.Event())
}

// Type returns the stream event type
func (event VersionEvent) Type() StreamEventType {
	return VersionEventT
}

// VersionEvent represents a version event that occurs upon any
// modifications to the current state of the subscribed document.
type SetEvent struct {
	StreamEvent
	txn   int64
	event Value
}

// Txn returns the stream event timestamp
func (event SetEvent) Txn() int64 {
	return event.txn
}

// Event returns the stream event as a `Value`
func (event SetEvent) Event() Value {
	return event.event
}

func (event SetEvent) String() string {
	return fmt.Sprintf("SetEvent{txn=%d, event=%s}", event.Txn(), event.Event())
}

// Type returns the stream event type
func (event SetEvent) Type() StreamEventType {
	return SetEventT
}

// HistoryRewriteEvent represents a history rewrite event which occurs upon any modifications
// to the history of the subscribed document.
type HistoryRewriteEvent struct {
	StreamEvent
	txn   int64
	event Value
}

// Txn returns the stream event timestamp
func (event HistoryRewriteEvent) Txn() int64 {
	return event.txn
}

// Event returns the stream event as a `Value`
func (event HistoryRewriteEvent) Event() Value {
	return event.event
}

func (event HistoryRewriteEvent) String() string {
	return fmt.Sprintf("HistoryRewriteEvent{txn=%d, event=%s}", event.Txn(), event.Event())
}

// Type returns the stream event type
func (event HistoryRewriteEvent) Type() StreamEventType {
	return HistoryRewriteEventT
}

// ErrorEvent represents an error event fired both for client and server errors
// that may occur as a result of a subscription.
type ErrorEvent struct {
	StreamEvent
	txn int64
	err error
}

// Type returns the stream event type
func (event ErrorEvent) Type() StreamEventType {
	return ErrorEventT
}

// Txn returns the stream event timestamp
func (event ErrorEvent) Txn() int64 {
	return event.txn
}

// Error returns the event error
func (event ErrorEvent) Error() error {
	return event.err
}

func (event ErrorEvent) String() string {
	return fmt.Sprintf("ErrorEvent{error=%s}", event.err)
}

func unMarshalStreamEvent(data Obj) (evt StreamEvent, err error) {
	if tpe, ok := data["type"]; ok {
		switch StreamEventType(tpe.(StringV)) {
		case StartEventT:
			evt = StartEvent{
				txn:   int64(data["txn"].(LongV)),
				event: data["event"].(LongV),
			}
		case VersionEventT:
			evt = VersionEvent{
				txn:   int64(data["txn"].(LongV)),
				event: data["event"].(ObjectV),
			}
		case SetEventT:
			evt = SetEvent{
				txn:   int64(data["txn"].(LongV)),
				event: data["event"].(ObjectV),
			}
		case ErrorEventT:
			evt = ErrorEvent{
				txn: int64(data["txn"].(LongV)),
				err: errorFromStreamError(data["event"].(ObjectV)),
			}
		case HistoryRewriteEventT:
			evt = HistoryRewriteEvent{
				txn:   int64(data["txn"].(LongV)),
				event: data["event"].(ObjectV),
			}
		}
	} else {
		err = errors.New("unparseable event type")
	}
	return
}
