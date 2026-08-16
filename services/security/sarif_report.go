package security

// SARIFReport is a minimal normalized ingestion envelope for static-analysis reports.
type SARIFReport struct {
	Tool string
	RunID string
	Findings []Finding
}

func (r SARIFReport) Valid() bool {
	return r.Tool != "" && r.RunID != ""
}
