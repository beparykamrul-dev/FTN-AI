package ftnstorage

// ErasurePlan describes a space-efficient recovery layout. It is declarative;
// an executor must validate policy before writing shards.
type ErasurePlan struct {
	DataShards   int
	ParityShards int
	ChunkSize    uint64
}

func (p ErasurePlan) Valid() bool {
	return p.DataShards > 0 && p.ParityShards > 0 && p.ChunkSize > 0
}

func (p ErasurePlan) TotalShards() int { return p.DataShards + p.ParityShards }
