package ftnmap

import "sort"

// EventAggregator keeps the latest topology event per node/service identity.
type EventAggregator struct {
	latest map[string]TopologyEvent
}

func NewEventAggregator() *EventAggregator { return &EventAggregator{latest: make(map[string]TopologyEvent)} }

func (a *EventAggregator) Apply(e TopologyEvent) bool {
	if !e.Valid() { return false }
	key := e.NodeID + ":" + e.ServiceID
	old, ok := a.latest[key]
	if ok && !e.ObservedAt.After(old.ObservedAt) { return false }
	a.latest[key] = e
	return true
}

func (a *EventAggregator) Snapshot() []TopologyEvent {
	out := make([]TopologyEvent, 0, len(a.latest))
	for _, e := range a.latest { out = append(out, e) }
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
