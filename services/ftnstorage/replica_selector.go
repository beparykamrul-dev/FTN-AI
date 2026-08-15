package ftnstorage

// Replica describes a possible source for a repaired chunk.
type Replica struct {
	NodeID    string
	Healthy   bool
	LatencyMS uint32
	Verified  bool
}

// SelectReplica deterministically chooses the healthiest verified replica
// with the lowest latency.
func SelectReplica(replicas []Replica) (Replica, bool) {
	var best Replica
	found := false
	for _, r := range replicas {
		if !r.Healthy || !r.Verified {
			continue
		}
		if !found || r.LatencyMS < best.LatencyMS ||
			(r.LatencyMS == best.LatencyMS && r.NodeID < best.NodeID) {
			best = r
			found = true
		}
	}
	return best, found
}
