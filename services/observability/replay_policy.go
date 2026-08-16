package observability

// ReplayPolicy controls safe replay of buffered telemetry after recovery.
type ReplayPolicy struct {
	BatchSize uint64
	MaxAttempts uint32
}

func (p ReplayPolicy) Valid() bool {
	return p.BatchSize > 0 && p.MaxAttempts > 0
}
