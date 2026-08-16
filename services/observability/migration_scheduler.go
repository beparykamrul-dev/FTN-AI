package observability

// MigrationScheduler provides bounded scheduling decisions for resumable transfers.
type MigrationScheduler struct {
	MaxConcurrent uint32
	ChunkSizeBytes uint64
}

func (s MigrationScheduler) Valid() bool {
	return s.MaxConcurrent > 0 && s.ChunkSizeBytes > 0
}
