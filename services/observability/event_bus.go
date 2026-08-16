package observability

// EventBus is the application-level boundary for publishing telemetry events.
type EventBus interface {
	Publish(LogEvent) error
}

// EventPublisher adapts a function to EventBus without coupling the package to a broker.
type EventPublisher struct { PublishFunc func(LogEvent) error }

func (p EventPublisher) Publish(e LogEvent) error {
	if p.PublishFunc == nil { return nil }
	return p.PublishFunc(e)
}
