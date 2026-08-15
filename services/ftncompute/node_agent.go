package ftncompute

import "time"

// NodeHeartbeat is the control-plane health snapshot published by an FTN node.
type NodeHeartbeat struct {
	NodeID      string
	Epoch       uint64
	ObservedAt  time.Time
	Healthy     bool
	CPUFree     uint64
	MemoryFree  uint64
	StorageFree uint64
	LatencyMS   uint32
}

func (h NodeHeartbeat) Valid() bool {
	return h.NodeID != "" && h.Epoch > 0 && !h.ObservedAt.IsZero()
}
