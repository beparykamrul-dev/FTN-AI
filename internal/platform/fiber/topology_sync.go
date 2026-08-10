package fiber

import "time"

type DiscoverySource string

const (
	SourceOLT DiscoverySource = "olt"
	SourceRouter DiscoverySource = "router"
	SourcePPPoE DiscoverySource = "pppoe"
	SourceGIS DiscoverySource = "gis"
)

type DiscoveredEntity struct {
	Source DiscoverySource `json:"source"`
	ExternalID string `json:"external_id"`
	Kind string `json:"kind"`
	ParentID string `json:"parent_id,omitempty"`
	Name string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	ObservedAt time.Time `json:"observed_at"`
}

type TopologySnapshot struct {
	ObservedAt time.Time `json:"observed_at"`
	Entities []DiscoveredEntity `json:"entities"`
}

// BuildTopologySnapshot normalizes discovery results from OLT/router/PPPoE/GIS
// adapters into one topology stream. It performs no configuration changes.
func BuildTopologySnapshot(entities []DiscoveredEntity) TopologySnapshot {
	out := make([]DiscoveredEntity, 0, len(entities))
	for _, e := range entities {
		if e.ExternalID == "" || e.Kind == "" { continue }
		if e.ObservedAt.IsZero() { e.ObservedAt = time.Now().UTC() }
		out = append(out, e)
	}
	return TopologySnapshot{ObservedAt: time.Now().UTC(), Entities: out}
}
