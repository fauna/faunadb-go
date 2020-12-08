package faunadb

import "sync"

func newSubscription(client *FaunaClient, query Expr, config ...StreamConfig) StreamSubscription {
	sub := StreamSubscription{
		query,
		streamDispatcher{},
		streamConfig{
			[]StreamField{},
		},
		client,
		streamConnectionStatus{status: StreamConnIdle},
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
	query      Expr
	dispatcher streamDispatcher
	config     streamConfig
	client     *FaunaClient
	status     streamConnectionStatus
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
	return sub.client.startStream(sub)
}

// On
func (sub *StreamSubscription) On(eventType StreamEventType, callback StreamEventCallback) {
	sub.dispatcher.On(eventType, callback)
}

// Close eventually closes the stream
func (sub *StreamSubscription) Close() {
	sub.status.Set(StreamConnClosed)
}
