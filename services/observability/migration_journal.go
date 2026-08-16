package observability

// MigrationJournalEntry provides an append-only audit record for storage movement.
type MigrationJournalEntry struct {
	MigrationID string
	State MigrationState
	SourceNode string
	TargetNode string
	Bytes uint64
	Verified bool
	Reason string
}

func (e MigrationJournalEntry) Valid() bool {
	return e.MigrationID != "" && e.SourceNode != "" && e.TargetNode != "" && e.State != ""
}
