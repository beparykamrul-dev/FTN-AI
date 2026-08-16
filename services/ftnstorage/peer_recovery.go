package ftnstorage

// PeerRecoveryPlan describes a server-to-server recovery transfer.
type PeerRecoveryPlan struct {
	ChunkID    string
	SourceNode string
	TargetNode string
	Reason     string
}

// PlanPeerRecovery builds a recovery plan only when both endpoints are
// healthy and verified. Actual transfer remains behind the execution layer.
func PlanPeerRecovery(chunkID string, source, target StoragePeer) (PeerRecoveryPlan, bool) {
	if chunkID == "" || !source.Healthy || !source.Verified || !target.Healthy || !target.Verified || source.NodeID == target.NodeID {
		return PeerRecoveryPlan{}, false
	}
	return PeerRecoveryPlan{ChunkID: chunkID, SourceNode: source.NodeID, TargetNode: target.NodeID, Reason: "verified-server-to-server-recovery"}, true
}
