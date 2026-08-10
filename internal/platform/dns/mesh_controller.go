package dns

import (
	"sort"
	"sync"
	"time"
)

type Node struct {
	ID string `json:"id"`
	Address string `json:"address"`
	Role string `json:"role"`
	Healthy bool `json:"healthy"`
	LastSeen time.Time `json:"last_seen"`
}

type MeshController struct {
	mu sync.RWMutex
	nodes map[string]Node
}

func NewMeshController() *MeshController {
	return &MeshController{nodes: make(map[string]Node)}
}

func (m *MeshController) Upsert(n Node) {
	m.mu.Lock(); defer m.mu.Unlock()
	m.nodes[n.ID] = n
}

func (m *MeshController) HealthyNodes() []Node {
	m.mu.RLock(); defer m.mu.RUnlock()
	out := make([]Node, 0, len(m.nodes))
	for _, n := range m.nodes { if n.Healthy { out = append(out, n) } }
	sort.Slice(out, func(i, j int) bool { return out[i].LastSeen.After(out[j].LastSeen) })
	return out
}

func (m *MeshController) Get(id string) (Node, bool) {
	m.mu.RLock(); defer m.mu.RUnlock()
	n, ok := m.nodes[id]
	return n, ok
}
