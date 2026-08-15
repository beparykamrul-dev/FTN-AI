package ftnstorage

import "time"

// HealthSnapshot is the storage control-plane health state.
type HealthSnapshot struct {
	NodeID       string
	PoolID       string
	Epoch        uint64
	ObservedAt   time.Time
	Healthy      bool
	Capacity     uint64
	Used         uint64
	Errors       uint64
	ChecksumOK   bool
	RepairNeeded bool
}

func (h HealthSnapshot) Valid() bool {
	return h.NodeID != "" && h.PoolID != "" && h.Epoch > 0 && !h.ObservedAt.IsZero()
}
