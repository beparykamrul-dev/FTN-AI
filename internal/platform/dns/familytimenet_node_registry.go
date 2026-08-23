package dns

import (
	"sort"
	"strings"
	"sync"
)

const FamilyTimeNetZone = "familytimenet.com"

type NodeScope string

const (
	NodeScopeLocal  NodeScope = "local"
	NodeScopeGlobal NodeScope = "global"
)

type DNSNode struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	ProviderID   string       `json:"providerId"`
	Scope        NodeScope    `json:"scope"`
	Region       string       `json:"region,omitempty"`
	Endpoint     string       `json:"endpoint,omitempty"`
	Enabled      bool         `json:"enabled"`
	Capabilities []string     `json:"capabilities,omitempty"`
	Zone         string       `json:"zone"`
}

type NodeRegistry struct {
	mu    sync.RWMutex
	nodes map[string]DNSNode
}

func NewNodeRegistry() *NodeRegistry { return &NodeRegistry{nodes: make(map[string]DNSNode)} }

func (r *NodeRegistry) Upsert(n DNSNode) bool {
	n.ID = strings.TrimSpace(n.ID)
	n.Name = strings.TrimSpace(n.Name)
	n.ProviderID = strings.TrimSpace(n.ProviderID)
	n.Region = strings.TrimSpace(n.Region)
	n.Zone = strings.TrimSpace(n.Zone)
	if n.ID == "" || n.Name == "" || n.ProviderID == "" || n.Zone != FamilyTimeNetZone {
		return false
	}
	if n.Scope != NodeScopeLocal && n.Scope != NodeScopeGlobal { return false }
	r.mu.Lock()
	r.nodes[n.ID] = n
	r.mu.Unlock()
	return true
}

func (r *NodeRegistry) Get(id string) (DNSNode, bool) {
	r.mu.RLock(); defer r.mu.RUnlock()
	n, ok := r.nodes[strings.TrimSpace(id)]
	return n, ok
}

func (r *NodeRegistry) List(scope NodeScope) []DNSNode {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]DNSNode, 0)
	for _, n := range r.nodes {
		if n.Enabled && (scope == "" || n.Scope == scope) { out = append(out, n) }
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Region != out[j].Region { return out[i].Region < out[j].Region }
		return out[i].ID < out[j].ID
	})
	return out
}
