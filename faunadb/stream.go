package faunadb

// StreamField represents a stream field
type StreamField string

const (
	DiffField     StreamField = "diff"
	PrevField     StreamField = "prev"
	DocumentField StreamField = "document"
	ActionField   StreamField = "action"
	IndexField    StreamField = "index"
)

type streamConfig struct {
	Fields []StreamField
}

// StreamConfig describes optional parameters for a stream subscription
type StreamConfig func(*StreamSubscription)

// StreamConnectionStatus is a expression shortcut to represent
// the stream connection status
type StreamConnectionStatus int

const (
	// StreamConnIdle represents an idle stream subscription
	StreamConnIdle StreamConnectionStatus = iota
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
