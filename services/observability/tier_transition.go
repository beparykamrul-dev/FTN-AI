package observability

// TierTransition describes a controlled movement of buffered telemetry between storage tiers.
type TierTransition struct {
	From StorageClass
	To StorageClass
	Reason string
	Priority uint32
}

func (t TierTransition) Valid() bool {
	return t.From != "" && t.To != "" && t.From != t.To && t.Reason != ""
}
