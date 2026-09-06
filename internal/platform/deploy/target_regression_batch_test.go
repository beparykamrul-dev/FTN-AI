package deploy

import "testing"

func TestTargetRejectsMissingIdentity(t *testing.T) {
	target := Target{Name: "node"}
	if target.Validate() == nil { t.Fatal("target without ID must be rejected") }
}

func TestTargetRejectsUnboundIdentity(t *testing.T) {
	target := Target{ID: "n1", Name: "node"}
	if target.Validate() == nil { t.Fatal("target without serial or authenticated agent ID must be rejected") }
}
