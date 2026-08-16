package observability

// MigrationChunk tracks resumable progress for a large replica transfer.
type MigrationChunk struct {
	MigrationID string
	Offset uint64
	Length uint64
	Checksum string
	Completed bool
}

func (c MigrationChunk) Valid() bool {
	return c.MigrationID != "" && c.Length > 0
}
