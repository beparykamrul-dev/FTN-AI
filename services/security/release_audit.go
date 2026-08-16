package security

// ReleaseAudit links a rollout decision to its immutable security evidence.
type ReleaseAudit struct {
	ReleaseID   string
	RunID       string
	Fingerprint string
	Policy      PolicyVersion
	Allowed     bool
	Reason      string
}

func (a ReleaseAudit) Valid() bool {
	return a.ReleaseID != "" && a.RunID != "" && a.Fingerprint != "" && a.Policy.Valid() && a.Reason != ""
}
