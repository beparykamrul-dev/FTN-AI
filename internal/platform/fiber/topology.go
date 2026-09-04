package fiber

import "sync"

type TopologyNode struct {
	ID string `json:"id"`
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
	Latitude float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Status string `json:"status"`
}

type TopologyEdge struct {
	ID string `json:"id,omitempty"`
	From string `json:"from"`
	To string `json:"to"`
	Kind string `json:"kind,omitempty"`
	Source DiscoverySource `json:"source,omitempty"`
	DistanceMeters float64 `json:"distance_meters,omitempty"`
	Status string `json:"status,omitempty"`
}

type Topology struct {
	mu sync.RWMutex
	nodes map[string]TopologyNode
	edges map[string]TopologyEdge
}

func NewTopology() *Topology { return &Topology{nodes: make(map[string]TopologyNode), edges: make(map[string]TopologyEdge)} }
func (t *Topology) UpsertNode(n TopologyNode) { t.mu.Lock(); defer t.mu.Unlock(); t.nodes[n.ID] = n }
func (t *Topology) UpsertEdge(e TopologyEdge) { t.mu.Lock(); defer t.mu.Unlock(); t.edges[e.ID] = e }
func (t *Topology) Snapshot() (nodes []TopologyNode, edges []TopologyEdge) {
	t.mu.RLock(); defer t.mu.RUnlock()
	for _, n := range t.nodes { nodes = append(nodes, n) }
	for _, e := range t.edges { edges = append(edges, e) }
	return
}
