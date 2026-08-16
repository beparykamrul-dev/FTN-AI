package edge

// FailoverState tracks provider-route failover without performing network I/O.
type FailoverState struct {
	ActiveProvider string
	ConsecutiveFailures uint32
	Threshold uint32
}

func (s FailoverState) ShouldFailover() bool {
	return s.Threshold > 0 && s.ConsecutiveFailures >= s.Threshold
}
