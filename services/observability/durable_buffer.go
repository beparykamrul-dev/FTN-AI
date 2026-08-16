package observability

// DurableBuffer defines the persistence boundary for telemetry waiting for delivery.
type DurableBuffer struct {
	MaxItems uint64
	Items uint64
}

func (b DurableBuffer) CanAccept() bool {
	return b.MaxItems > 0 && b.Items < b.MaxItems
}
