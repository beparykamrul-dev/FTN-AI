package mesh

import (
	"sort"
	"strings"
	"sync"
)

type Scope string

const (
	Local  Scope = "local"
	Global Scope = "global"
)

type FullMesh struct {
	mu    sync.RWMutex
	nodes map[string]Node
	links map[string]Link
}

func NewFullMesh() *FullMesh { return &FullMesh{nodes: map[string]Node{}, links: map[string]Link{}} }

func linkKey(a, b string) string {
	a, b = strings.TrimSpace(a), strings.TrimSpace(b)
	if a > b { a, b = b, a }
	return a + "|" + b
}

func (m *FullMesh) UpsertNode(n Node) bool {
	n.ID, n.Region, n.Endpoint = strings.TrimSpace(n.ID), strings.TrimSpace(n.Region), strings.TrimSpace(n.Endpoint)
	if n.ID == "" || n.Scope == "" { return false }
	m.mu.Lock(); m.nodes[n.ID] = n; m.mu.Unlock()
	return true
}

func (m *FullMesh) ObserveLink(l Link) bool {
	if l.From == "" || l.To == "" || l.From == l.To || l.RTTMillis < 0 || l.Loss < 0 || l.Loss > 1 || l.JitterMs < 0 { return false }
	m.mu.Lock(); m.links[linkKey(l.From, l.To)] = l; m.mu.Unlock()
	return true
}

func (m *FullMesh) Nodes(scope Scope) []Node {
	m.mu.RLock(); defer m.mu.RUnlock()
	out := make([]Node, 0)
	for _, n := range m.nodes { if n.Enabled && (scope == "" || n.Scope == scope) { out = append(out, n) } }
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (m *FullMesh) Links() []Link {
	m.mu.RLock(); defer m.mu.RUnlock()
	out := make([]Link, 0, len(m.links))
	for _, l := range m.links { out = append(out, l) }
	return out
}
