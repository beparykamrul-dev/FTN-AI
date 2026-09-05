package ftnftp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type NodeState string

const (
	NodeHealthy NodeState = "healthy"
	NodeDegraded NodeState = "degraded"
	NodeOffline NodeState = "offline"
)

type Node struct {
	ID            string    `json:"id"`
	State         NodeState `json:"state"`
	CapacityBytes int64     `json:"capacity_bytes"`
	UsedBytes     int64     `json:"used_bytes"`
	LastSeen      time.Time `json:"last_seen"`
}

type Runtime struct {
	mu        sync.RWMutex
	nodes     map[string]Node
	objects   map[string]Object
	transfers map[string]Transfer
}

func NewRuntime() *Runtime {
	return &Runtime{nodes: map[string]Node{}, objects: map[string]Object{}, transfers: map[string]Transfer{}}
}

func (r *Runtime) UpsertNode(n Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n.LastSeen = time.Now().UTC()
	r.nodes[n.ID] = n
}

func (r *Runtime) PutObject(o Object, data []byte) error {
	if o.ID == "" || o.TenantID == "" || o.Bucket == "" || o.Key == "" {
		return fmt.Errorf("object identity is incomplete")
	}
	sum := sha256.Sum256(data)
	o.SHA256 = hex.EncodeToString(sum[:])
	o.Checksum = o.SHA256
	o.Size = int64(len(data))
	o.Version++
	o.State = "verified"
	o.Status = "verified"
	r.mu.Lock()
	defer r.mu.Unlock()
	r.objects[o.ID] = o
	return nil
}

func (r *Runtime) StartTransfer(t Transfer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	t.StartedAt = now
	t.UpdatedAt = now
	t.State = TransferRunning
	r.transfers[t.ID] = t
}

func (r *Runtime) UpdateTransfer(id string, offset int64, state string, errText string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.transfers[id]
	if !ok {
		return
	}
	t.Offset = offset
	t.State = TransferState(state)
	t.Error = errText
	t.UpdatedAt = time.Now().UTC()
	r.transfers[id] = t
}

func (r *Runtime) Snapshot() (nodes []Node, objects []Object, transfers []Transfer) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.nodes {
		nodes = append(nodes, v)
	}
	for _, v := range r.objects {
		objects = append(objects, v)
	}
	for _, v := range r.transfers {
		transfers = append(transfers, v)
	}
	return
}
