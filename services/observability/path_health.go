package observability

// PathHealth tracks consecutive path health observations for failover decisions.
type PathHealth struct {
	PathID string
	Healthy bool
	ConsecutiveFailures uint32
	ConsecutiveSuccesses uint32
}

func (h PathHealth) Valid() bool { return h.PathID != "" }
