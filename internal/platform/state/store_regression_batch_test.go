package state

import "testing"

func TestValidateDecisionAllowsOptionalTransport(t *testing.T) {
	d := Decision{ServiceID: "svc", PolicyVersion: "v1", Status: "approved"}
	if err := ValidateDecision(d); err != nil { t.Fatal(err) }
}

func TestValidateDecisionRejectsOversizedReason(t *testing.T) {
	d := Decision{ServiceID: "svc", PolicyVersion: "v1", Status: "approved", Reason: string(make([]byte, 4097))}
	if ValidateDecision(d) == nil { t.Fatal("oversized reason must be rejected") }
}
