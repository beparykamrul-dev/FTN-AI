package ftnmap

// SyncEnvelope is the transport-neutral boundary for propagating map state
// between FTN global and local map consumers.
type SyncEnvelope struct {
	Origin   string
	Revision uint64
	Events   []TopologyEvent
}

func (e SyncEnvelope) Valid() bool {
	if e.Origin == "" || e.Revision == 0 { return false }
	for _, event := range e.Events { if !event.Valid() { return false } }
	return true
}
