package gis

import "sync"

type TopologyLink struct {
	ID     string `json:"id"`
	FromID string `json:"from_id"`
	ToID   string `json:"to_id"`
	Kind   string `json:"kind"`
	Status string `json:"status,omitempty"`
}

type TopologyGraph struct {
	mu    sync.RWMutex
	links map[string]TopologyLink
}

func NewTopologyGraph() *TopologyGraph {
	return &TopologyGraph{links: make(map[string]TopologyLink)}
}

func (g *TopologyGraph) Upsert(l TopologyLink) {
	if l.ID == "" {
		return
	}
	g.mu.Lock()
	g.links[l.ID] = l
	g.mu.Unlock()
}

func (g *TopologyGraph) List() []TopologyLink {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]TopologyLink, 0, len(g.links))
	for _, l := range g.links {
		out = append(out, l)
	}
	return out
}

func (g *TopologyGraph) Neighbors(id string) []TopologyLink {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]TopologyLink, 0)
	for _, l := range g.links {
		if l.FromID == id || l.ToID == id {
			out = append(out, l)
		}
	}
	return out
}
