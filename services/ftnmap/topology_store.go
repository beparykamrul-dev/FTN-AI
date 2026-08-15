package ftnmap

// TopologyStore is an in-memory logical state boundary for the current FTN map.
type TopologyStore struct { events *EventAggregator }

func NewTopologyStore() *TopologyStore { return &TopologyStore{events: NewEventAggregator()} }
func (s *TopologyStore) Apply(e TopologyEvent) bool { return s.events.Apply(e) }
func (s *TopologyStore) Snapshot() []TopologyEvent { return s.events.Snapshot() }
