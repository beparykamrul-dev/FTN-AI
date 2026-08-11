package ftndns

import "time"

// CloudState is the provider-neutral state needed to scale FTNDNS from a
// local DNS service into a regional/global FTN Cloud control plane.
type CloudState struct {
	Region string `json:"region"`
	Epoch uint64 `json:"epoch"`
	UpdatedAt time.Time `json:"updated_at"`
	ConsensusHash string `json:"consensus_hash"`
	Nodes []MeshNode `json:"nodes"`
}

// BuildCloudState creates a deterministic control-plane snapshot from mesh
// nodes. It does not perform network replication or DNS mutation.
func BuildCloudState(region string, epoch uint64, nodes []MeshNode, consensusHash string, now time.Time) CloudState {
	return CloudState{Region: region, Epoch: epoch, UpdatedAt: now, ConsensusHash: consensusHash, Nodes: RankMeshNodes(nodes)}
}
