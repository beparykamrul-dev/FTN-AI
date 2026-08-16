package security

// Remediation records a proposed fix without automatically changing source code.
type Remediation struct {
	FindingRule string
	Owner       string
	Action      string
	Status      string
}

func (r Remediation) Valid() bool {
	return r.FindingRule != "" && r.Action != ""
}
