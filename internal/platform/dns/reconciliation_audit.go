package dns

// ReconciliationAudit is the immutable approval record consumed by the
// privileged DNS reconciliation boundary.
type ReconciliationAudit struct {
	ID           string `json:"id"`
	Zone         string `json:"zone"`
	ExpectedHash string `json:"expected_hash"`
	ObservedHash string `json:"observed_hash"`
	Action       ReconcileAction `json:"action"`
	Approved     bool `json:"approved"`
}

func (a ReconciliationAudit) CanExecute() bool {
	return a.Approved && a.Action == ReconcileSync && a.ID != ""
}
