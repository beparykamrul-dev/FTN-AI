package observability

// PathHysteresis avoids rapid failover/failback oscillation.
type PathHysteresis struct {
	FailureThreshold uint32
	RecoveryThreshold uint32
}

func (h PathHysteresis) Valid() bool { return h.FailureThreshold > 0 && h.RecoveryThreshold > 0 }

func (h PathHysteresis) ShouldFailover(p PathHealth) bool {
	return h.Valid() && !p.Healthy && p.ConsecutiveFailures >= h.FailureThreshold
}

func (h PathHysteresis) ShouldFailback(p PathHealth) bool {
	return h.Valid() && p.Healthy && p.ConsecutiveSuccesses >= h.RecoveryThreshold
}
