package faunadb

// StreamField represents a stream field
type StreamField string

const (
	DiffField     StreamField = "diff"
	PrevField     StreamField = "prev"
	DocumentField StreamField = "document"
	ActionField   StreamField = "action"
)

type streamConfig struct {
	Fields []StreamField
}

// StreamEventCallback describes the listener function for a
// stream event
type StreamEventCallback func(StreamEvent)

// StreamConfig describes optional parameters for a stream subscription
type StreamConfig func(*StreamSubscription)

// StreamConnectionStatus is a expression shortcut to represent
// the stream connection status
type StreamConnectionStatus int

const (
	// StreamConnIdle represents an idle stream subscription
	StreamConnIdle StreamConnectionStatus = iota
	// StreamConnOpening describes an opening stream subscription
	StreamConnOpening
	// StreamConnActive describes an active/established stream subscription
	StreamConnActive
	// StreamConnClosed describes a closed stream subscription
	StreamConnClosed
	// StreamConnError describes a stream subscription error
	StreamConnError
)

// Fields is optional stream parameter that describes the fields received on version and history_rewrite stream events.
func Fields(fields ...StreamField) StreamConfig {
	return func(sub *StreamSubscription) {
		sub.config.Fields = []StreamField(fields)
	}
}

// OnStart is an optional stream parameter that specifies the listener function for the start event
// alternatively `subscription.On(f.StartEventT, fn)` can be used
func OnStart(fn StreamEventCallback) StreamConfig {
	return func(sub *StreamSubscription) {
		sub.dispatcher.On(StartEventT, fn)
	}
}

// OnVersion is an optional stream parameter that specifies the listener function for a version event
// alternatively `subscription.On(f.VersionEventT, fn)` can be used
func OnVersion(fn StreamEventCallback) StreamConfig {
	return func(sub *StreamSubscription) {
		sub.dispatcher.On(VersionEventT, fn)
	}
}

// OnError is an optional stream parameter that specifies the listener function for an error event
// alternatively `subscription.On(f.ErrorEventT, fn)` can be used
func OnError(fn StreamEventCallback) StreamConfig {
	return func(sub *StreamSubscription) {
		sub.dispatcher.On("error", fn)
	}
}

// OnHistoryRewrite is an optional stream parameter that specifies the listener function for a history rewrite event
// alternatively `subscription.On(f.HistoryRewriteEventT, fn)` can be used
func OnHistoryRewrite(fn StreamEventCallback) StreamConfig {
	return func(sub *StreamSubscription) {
		sub.dispatcher.On("history_rewrite", fn)
	}
}
