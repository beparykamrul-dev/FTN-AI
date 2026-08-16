package ftnstorage

import "sort"

// Checkpoint is a verified storage recovery point.
type Checkpoint struct {
	ID      string
	NodeID  string
	Version uint64
	Valid   bool
}

// LatestCheckpoint returns the newest valid checkpoint deterministically.
func LatestCheckpoint(items []Checkpoint) (Checkpoint, bool) {
	valid := make([]Checkpoint, 0, len(items))
	for _, c := range items {
		if c.ID != "" && c.NodeID != "" && c.Valid {
			valid = append(valid, c)
		}
	}
	if len(valid) == 0 {
		return Checkpoint{}, false
	}
	sort.SliceStable(valid, func(i, j int) bool {
		if valid[i].Version != valid[j].Version {
			return valid[i].Version > valid[j].Version
		}
		return valid[i].ID < valid[j].ID
	})
	return valid[0], true
}
