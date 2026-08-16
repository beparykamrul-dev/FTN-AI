package observability

// BufferStore is the durable-buffer persistence boundary.
type BufferStore interface {
	Append(event LogEvent) error
	Pending(limit uint32) ([]LogEvent, error)
	Ack(event LogEvent) error
}
