package gis

import "sync"

type MapStore struct {
	mu    sync.RWMutex
	nodes map[string]MapNode
	edges map[string]MapEdge
}

func NewMapStore() *MapStore {
	return &MapStore{nodes: make(map[string]MapNode), edges: make(map[string]MapEdge)}
}

func (m *MapStore) UpsertNode(n MapNode) {
	m.mu.Lock(); defer m.mu.Unlock(); m.nodes[n.ID] = n
}

func (m *MapStore) UpsertEdge(e MapEdge) {
	m.mu.Lock(); defer m.mu.Unlock(); m.edges[e.ID] = e
}

func (m *MapStore) Snapshot() MapSnapshot {
	m.mu.RLock(); defer m.mu.RUnlock()
	s := MapSnapshot{Nodes: make([]MapNode, 0, len(m.nodes)), Edges: make([]MapEdge, 0, len(m.edges))}
	for _, n := range m.nodes { s.Nodes = append(s.Nodes, n) }
	for _, e := range m.edges { s.Edges = append(s.Edges, e) }
	return s
}
