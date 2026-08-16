package edge

import "time"

// TrafficEvent stores privacy-minimized ingress telemetry rather than raw
// customer URLs or payloads.
type TrafficEvent struct {
	Provider  string
	Kind      string
	NodeID    string
	Status    int
	Bytes     uint64
	LatencyMS uint32
	CacheHit  bool
	At        time.Time
}

func (e TrafficEvent) Valid() bool {
	return e.Provider != "" && e.Kind != "" && e.NodeID != "" && !e.At.IsZero()
}
