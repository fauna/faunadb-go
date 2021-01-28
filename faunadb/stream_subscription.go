package faunadb

import "sync"

func newSubscription(client *FaunaClient, query Expr, config ...StreamConfig) StreamSubscription {
	sub := StreamSubscription{
		query,
		streamConfig{
			[]StreamField{},
		},
		client,
		streamConnectionStatus{status: StreamConnIdle},
		make(chan StreamEvent),
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
	query    Expr
	config   streamConfig
	client   *FaunaClient
	status   streamConnectionStatus
	messages chan StreamEvent
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

func isClosed(ch <-chan StreamEvent) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

// Close eventually closes the stream
func (sub *StreamSubscription) Close() {
	if !isClosed(sub.messages) {
		sub.status.Set(StreamConnClosed)
		close(sub.messages)
	}
}

func (sub *StreamSubscription) Messages() <-chan StreamEvent {
	return sub.messages
}
