package infrastructure

import (
	"sync"
	"time"
)

type NodeKind string

const (
	NodeSite NodeKind = "site"
	NodeServer NodeKind = "server"
	NodeNetworkDevice NodeKind = "network_device"
	NodeDNS NodeKind = "dns"
	NodeApplication NodeKind = "application"
	NodeContainer NodeKind = "container"
	NodeCloudResource NodeKind = "cloud_resource"
	NodeDatabase NodeKind = "database"
)

type Node struct {
	ID         string            `json:"id"`
	Kind       NodeKind          `json:"kind"`
	Name       string            `json:"name"`
	Provider   string            `json:"provider,omitempty"`
	Address    string            `json:"address,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	Desired    map[string]any    `json:"desired,omitempty"`
	Observed   map[string]any    `json:"observed,omitempty"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type Edge struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Kind   string `json:"kind"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Store struct {
	mu    sync.RWMutex
	nodes map[string]Node
	edges map[string]Edge
}

func NewStore() *Store {
	return &Store{nodes: make(map[string]Node), edges: make(map[string]Edge)}
}

func (s *Store) UpsertNode(n Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n.UpdatedAt.IsZero() {
		n.UpdatedAt = time.Now().UTC()
	}
	s.nodes[n.ID] = n
}

func (s *Store) UpsertEdge(e Edge) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.edges[e.ID] = e
}

func (s *Store) Snapshot() Graph {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g := Graph{Nodes: make([]Node, 0, len(s.nodes)), Edges: make([]Edge, 0, len(s.edges))}
	for _, n := range s.nodes { g.Nodes = append(g.Nodes, n) }
	for _, e := range s.edges { g.Edges = append(g.Edges, e) }
	return g
}
