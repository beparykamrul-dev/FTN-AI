package ftnstorage

// ErasureRecoveryPlan links a chunk to a space-efficient shard layout.
type ErasureRecoveryPlan struct {
	Chunk ChunkRef
	Plan  ErasurePlan
}

func (p ErasureRecoveryPlan) Valid() bool {
	return p.Chunk.Size > 0 && p.Plan.Valid()
}
