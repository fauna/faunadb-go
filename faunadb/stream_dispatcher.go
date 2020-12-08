package faunadb

type streamDispatcher struct {
	OnError   StreamEventCallback
	OnHistory StreamEventCallback
	OnStart   StreamEventCallback
	OnVersion StreamEventCallback
}

func (dispatch *streamDispatcher) On(eventType StreamEventType, callback StreamEventCallback) {
	switch eventType {
	case ErrorEventT:
		dispatch.OnError = callback
	case StartEventT:
		dispatch.OnStart = callback
	case VersionEventT:
		dispatch.OnVersion = callback
	case HistoryRewriteEventT:
		dispatch.OnHistory = callback
	}
}

func (dispatch streamDispatcher) Dispatch(event StreamEvent) {
	var fn StreamEventCallback
	switch event.Type() {
	case StartEventT:
		fn = dispatch.OnStart
	case VersionEventT:
		fn = dispatch.OnVersion
	case ErrorEventT:
		fn = dispatch.OnError
	case HistoryRewriteEventT:
		fn = dispatch.OnHistory
	default:
		return
	}
	fn(event)
}
