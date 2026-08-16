package ftnstorage

// SelectStoragePeer chooses the best allowed peer by latency, then free capacity.
func SelectStoragePeer(peers []StoragePeer, policy PeerPolicy) (StoragePeer, bool) {
	var best StoragePeer
	found := false
	for _, p := range peers {
		if !policy.Allows(p) { continue }
		if !found || p.LatencyMS < best.LatencyMS ||
			(p.LatencyMS == best.LatencyMS && p.CapacityFree > best.CapacityFree) ||
			(p.LatencyMS == best.LatencyMS && p.CapacityFree == best.CapacityFree && p.NodeID < best.NodeID) {
			best, found = p, true
		}
	}
	return best, found
}
