package ftnstorage

// ReconstructionPlan describes erasure-based reconstruction without executing it.
type ReconstructionPlan struct {
	ChunkHash   string
	DataShards  uint16
	ParityShards uint16
	Sources     []string
	TargetNode  string
}

func (p ReconstructionPlan) Valid() bool {
	return p.ChunkHash != "" && p.DataShards > 0 && p.ParityShards > 0 && p.TargetNode != "" && len(p.Sources) >= int(p.DataShards)
}
