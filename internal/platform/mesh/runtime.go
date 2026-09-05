package mesh

import (
	"sort"
	"sync"
	"time"
)

// RuntimeNode is the mesh-facing runtime identity for a server or FTN device.
type RuntimeNode struct {
	ID       string    `json:"id"`
	DeviceID string    `json:"device_id"`
	Address  string    `json:"address"`
	Runtime  string    `json:"runtime"` // docker, podman, k8s, linux, virtual
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

// RuntimeRegistry keeps the mesh view separate from device inventory. This
// lets Docker/Kubernetes/FTN OS adapters evolve without changing the graph.
type RuntimeRegistry struct {
	mu    sync.RWMutex
	nodes map[string]RuntimeNode
}

func NewRuntimeRegistry() *RuntimeRegistry {
	return &RuntimeRegistry{nodes: make(map[string]RuntimeNode)}
}

func (r *RuntimeRegistry) Upsert(n RuntimeNode) {
	if r == nil || n.ID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[n.ID] = n
}

func (r *RuntimeRegistry) Get(id string) (RuntimeNode, bool) {
	if r == nil {
		return RuntimeNode{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.nodes[id]
	return n, ok
}

func (r *RuntimeRegistry) Snapshot() []RuntimeNode {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]RuntimeNode, 0, len(r.nodes))
	for _, n := range r.nodes {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
