package ftnstorage

// CapacityNode is a storage placement candidate.
type CapacityNode struct {
	ID        string
	Healthy   bool
	UsedBytes uint64
	TotalBytes uint64
}

// BalanceCandidate returns the healthiest node with the lowest utilization
// while requiring positive capacity.
func BalanceCandidate(nodes []CapacityNode) (CapacityNode, bool) {
	var best CapacityNode
	var bestUsed, bestTotal uint64
	found := false
	for _, n := range nodes {
		if !n.Healthy || n.TotalBytes == 0 || n.UsedBytes > n.TotalBytes { continue }
		if !found || n.UsedBytes*bestTotal < bestUsed*n.TotalBytes ||
			(n.UsedBytes*bestTotal == bestUsed*n.TotalBytes && n.ID < best.ID) {
			best, bestUsed, bestTotal, found = n, n.UsedBytes, n.TotalBytes, true
		}
	}
	return best, found
}
