package observability

// MigrationResume describes the point from which an interrupted transfer can continue.
type MigrationResume struct {
	MigrationID string
	NextOffset uint64
	TotalBytes uint64
	VerifiedChunks uint64
}

func (r MigrationResume) Valid() bool {
	return r.MigrationID != "" && r.NextOffset <= r.TotalBytes
}
