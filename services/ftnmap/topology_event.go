package ftnmap

import "time"

// TopologyEvent represents a normalized FTN infrastructure change for the map layer.
type TopologyEvent struct {
	ID        string
	NodeID    string
	ServiceID string
	Region    string
	State     string
	ObservedAt time.Time
}

func (e TopologyEvent) Valid() bool {
	return e.ID != "" && e.NodeID != "" && e.ServiceID != "" && !e.ObservedAt.IsZero()
}
