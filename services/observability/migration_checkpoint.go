package observability

// MigrationCheckpoint persists verified transfer progress for restart/resume.
type MigrationCheckpoint struct {
	MigrationID string
	NextOffset uint64
	VerifiedBytes uint64
	LastChecksum string
}

func (c MigrationCheckpoint) Valid() bool {
	return c.MigrationID != "" && c.NextOffset >= c.VerifiedBytes
}
