package security

// EvidenceStore is the persistence boundary for security evidence.
// Database adapters (PostgreSQL/CockroachDB/etc.) implement it outside this package.
type EvidenceStore interface {
	SaveReleaseAudit(ReleaseAudit) error
	SaveDecision(DecisionRecord) error
	SaveScanRun(ScanRun) error
}
