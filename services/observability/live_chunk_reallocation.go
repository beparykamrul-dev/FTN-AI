package observability

// LiveChunkReallocation identifies only pending chunks that should move after path conditions change.
type LiveChunkReallocation struct {
	MigrationID string
	FromPath string
	ToPath string
	PendingChunks uint64
	OldCapacityMbps float64
	NewCapacityMbps float64
}

func (r LiveChunkReallocation) Valid() bool {
	return r.MigrationID != "" && r.FromPath != "" && r.ToPath != "" && r.FromPath != r.ToPath && r.PendingChunks > 0 && r.NewCapacityMbps >= 0
}
