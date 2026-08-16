package observability

// PathReassignment describes moving only pending chunks away from an unusable path.
type PathReassignment struct {
	MigrationID string
	FromPath string
	ToPath string
	PendingChunks uint64
	Reason string
}

func (r PathReassignment) Valid() bool {
	return r.MigrationID != "" && r.FromPath != "" && r.ToPath != "" && r.FromPath != r.ToPath && r.Reason != ""
}
