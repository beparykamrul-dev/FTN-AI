package observability

// MigrationChecksum records integrity evidence for a replica transfer.
type MigrationChecksum struct {
	MigrationID string
	Algorithm string
	SourceDigest string
	TargetDigest string
}

func (c MigrationChecksum) Verified() bool {
	return c.MigrationID != "" && c.Algorithm != "" && c.SourceDigest != "" && c.SourceDigest == c.TargetDigest
}
