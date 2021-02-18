package faunadb

import (
	"sync"
)

func newSubscription(client *FaunaClient, query Expr, config ...StreamConfig) StreamSubscription {
	sub := StreamSubscription{
		query,
		streamConfig{
			[]StreamField{},
		},
		client,
		streamConnectionStatus{status: StreamConnIdle},
		make(chan StreamEvent),
		make(chan bool),
		make(chan bool),
	}
	for _, fn := range config {
		fn(&sub)
	}
	return sub
}

type streamConnectionStatus struct {
	mu     sync.Mutex
	status StreamConnectionStatus
}

func (s *streamConnectionStatus) Set(status StreamConnectionStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = status
}

func (s *streamConnectionStatus) Get() StreamConnectionStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

// StreamSubscription dispatches events received to the registered listener functions.
// New subscriptions must be constructed via the FaunaClient stream method.
type StreamSubscription struct {
	query          Expr
	config         streamConfig
	client         *FaunaClient
	status         streamConnectionStatus
	eventsMessages chan StreamEvent
	closingMessage chan bool
	getNext        chan bool
}

// Query returns the query used to initiate the stream
func (sub *StreamSubscription) Query() Expr {
	return sub.query
}

// Status returns the current stream status
func (sub *StreamSubscription) Status() StreamConnectionStatus {
	return sub.status.Get()
}

// Start initiates the stream subscription.
func (sub *StreamSubscription) Start() error {
	if sub.status.Get() != StreamConnActive {
		go func() {
			sub.getNext <- true
		}()
	}

	return sub.client.startStream(sub)
}

func (sub *StreamSubscription) Close() {
	go func() {
		sub.closingMessage <- true
	}()
}

func (sub *StreamSubscription) EventsMessages() chan StreamEvent {
	return sub.eventsMessages
}

func (sub *StreamSubscription) Request() {
	go func() {
		sub.getNext <- true
	}()
}
