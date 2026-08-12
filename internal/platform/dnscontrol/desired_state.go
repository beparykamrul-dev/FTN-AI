package dnscontrol

import "time"

type Record struct {
	Name string `json:"name"`
	Type string `json:"type"`
	TTL uint32 `json:"ttl"`
	Values []string `json:"values"`
}

type Zone struct {
	Name string `json:"name"`
	Serial uint64 `json:"serial"`
	DNSSECEnabled bool `json:"dnssec_enabled"`
	Records []Record `json:"records"`
	Revision string `json:"revision"`
}

type NodeState struct {
	NodeID string `json:"node_id"`
	ZoneRevision string `json:"zone_revision"`
	Healthy bool `json:"healthy"`
	ObservedAt time.Time `json:"observed_at"`
}

type ConsistencyState string

const (
	Consistent ConsistencyState = "consistent"
	Drifted ConsistencyState = "drifted"
	Unavailable ConsistencyState = "unavailable"
)

func CompareRevision(desired string, node NodeState) ConsistencyState {
	if !node.Healthy { return Unavailable }
	if desired == node.ZoneRevision { return Consistent }
	return Drifted
}
