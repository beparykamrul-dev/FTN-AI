package dns

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type AnycastNode struct {
	ID string `json:"id"`
	Address string `json:"address"`
	Region string `json:"region"`
	Enabled bool `json:"enabled"`
	Healthy bool `json:"healthy"`
	Priority uint32 `json:"priority"`
}

type MeshAnycastConfig struct {
	Enabled bool `json:"enabled"`
	Nodes []AnycastNode `json:"nodes"`
}

type MeshAnycastController struct {
	mu sync.RWMutex
	nodes map[string]AnycastNode
	enabled bool
}

func NewMeshAnycastController() *MeshAnycastController {
	return &MeshAnycastController{nodes: make(map[string]AnycastNode)}
}

func (c *MeshAnycastController) SetEnabled(enabled bool) {
	c.mu.Lock(); defer c.mu.Unlock(); c.enabled = enabled
}

func (c *MeshAnycastController) Upsert(node AnycastNode) error {
	if strings.TrimSpace(node.ID) == "" || strings.TrimSpace(node.Address) == "" { return fmt.Errorf("node ID and address are required") }
	c.mu.Lock(); defer c.mu.Unlock(); c.nodes[node.ID] = node
	return nil
}

func (c *MeshAnycastController) Config() MeshAnycastConfig {
	c.mu.RLock(); defer c.mu.RUnlock()
	nodes := make([]AnycastNode, 0, len(c.nodes))
	for _, n := range c.nodes { nodes = append(nodes, n) }
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Priority < nodes[j].Priority })
	return MeshAnycastConfig{Enabled: c.enabled, Nodes: nodes}
}

func (c *MeshAnycastController) HealthyNodes() []AnycastNode {
	c.mu.RLock(); defer c.mu.RUnlock()
	out := make([]AnycastNode, 0, len(c.nodes))
	if !c.enabled { return out }
	for _, n := range c.nodes { if n.Enabled && n.Healthy { out = append(out, n) } }
	sort.Slice(out, func(i, j int) bool { return out[i].Priority < out[j].Priority })
	return out
}
