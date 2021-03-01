package faunadb

import (
	"errors"
	"sync"
)

// StreamSubscription dispatches events received to the registered listener functions.
// New subscriptions must be constructed via the FaunaClient stream method.
type StreamSubscription struct {
	mu     sync.Mutex
	query  Expr
	config streamConfig
	client *FaunaClient
	status StreamConnectionStatus
	events chan StreamEvent
	closed chan bool
}

func newSubscription(client *FaunaClient, query Expr, config ...StreamConfig) StreamSubscription {
	sub := StreamSubscription{
		query: query,
		config: streamConfig{
			[]StreamField{},
		},
		client: client,
		status: StreamConnIdle,
		events: make(chan StreamEvent),
		closed: make(chan bool),
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

// Query returns the query used to initiate the stream
func (sub *StreamSubscription) Query() Expr {
	return sub.query
}

// Status returns the current stream status
func (sub *StreamSubscription) Status() StreamConnectionStatus {
	sub.mu.Lock()
	defer sub.mu.Unlock()
	return sub.status
}

func (sub *StreamSubscription) Start() (err error) {
	sub.mu.Lock()
	defer sub.mu.Unlock()

	if sub.status != StreamConnIdle {
		err = errors.New("stream subscription already started")
	} else {
		sub.status = StreamConnActive
		if err = sub.client.startStream(sub); err != nil {
			sub.status = StreamConnError
		}
	}
	return
}


func (sub *StreamSubscription) Close() {
	sub.mu.Lock()
	defer sub.mu.Unlock()
	if sub.status == StreamConnActive {
		sub.status = StreamConnClosed
		close(sub.closed)
		close(sub.events)
	}
}

func (sub *StreamSubscription) StreamEvents() <-chan StreamEvent {
	return sub.events
}