package ftnstorage

import "time"

// SnapshotPolicy controls retention of recoverable storage checkpoints.
type SnapshotPolicy struct {
	Interval time.Duration
	Keep     uint16
	Verified bool
}

func (p SnapshotPolicy) Valid() bool {
	return p.Interval > 0 && p.Keep > 0 && p.Verified
}
