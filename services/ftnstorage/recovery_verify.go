package ftnstorage

// RecoveryResult records the outcome of a storage repair without embedding
// backend-specific write operations in the coordinator.
type RecoveryResult struct {
	Chunk       ChunkRef
	TargetNode  string
	Verified    bool
	SourceNode  string
	Error       string
}

func (r RecoveryResult) Successful() bool {
	return r.TargetNode != "" && r.SourceNode != "" && r.Verified && r.Error == ""
}
