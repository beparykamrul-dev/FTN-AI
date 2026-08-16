package ftnstorage

// StoragePeer describes a server-to-server storage relationship.
type StoragePeer struct {
	NodeID       string
	Region       string
	Healthy      bool
	Verified     bool
	LatencyMS    uint32
	CapacityFree uint64
}

// PeerPolicy controls which peers may exchange storage chunks.
type PeerPolicy struct {
	RequireHealthy  bool
	RequireVerified bool
	MaxLatencyMS    uint32
}

func (p PeerPolicy) Allows(peer StoragePeer) bool {
	if p.RequireHealthy && !peer.Healthy { return false }
	if p.RequireVerified && !peer.Verified { return false }
	if p.MaxLatencyMS > 0 && peer.LatencyMS > p.MaxLatencyMS { return false }
	return peer.NodeID != ""
}
