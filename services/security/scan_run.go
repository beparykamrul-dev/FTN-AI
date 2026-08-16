package security

// ScanRun identifies one reproducible security analysis run.
type ScanRun struct {
	RunID string
	Scanner string
	Commit string
	StartedUnix int64
	FinishedUnix int64
}

func (r ScanRun) Valid() bool {
	return r.RunID != "" && r.Scanner != "" && r.Commit != "" && r.FinishedUnix >= r.StartedUnix
}
