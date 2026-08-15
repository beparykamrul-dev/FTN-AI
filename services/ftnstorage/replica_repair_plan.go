package ftnstorage

// ReplicaRepairPlan describes replacement of an unhealthy replica.
type ReplicaRepairPlan struct {
	ChunkHash string
	SourceNode string
	TargetNode string
	VerifiedSource bool
}

func (p ReplicaRepairPlan) Valid() bool {
	return p.ChunkHash != "" && p.SourceNode != "" && p.TargetNode != "" && p.SourceNode != p.TargetNode && p.VerifiedSource
}
