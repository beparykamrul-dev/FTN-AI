package control

import "testing"

func TestValidateRejectsUnknownOperation(t *testing.T) {
	if Validate(Operation("unknown.operation")) == nil {
		t.Fatal("unknown operation must be rejected")
	}
}

func TestPrivilegedOperationsRequireApproval(t *testing.T) {
	for _, op := range []Operation{OpDNSWrite, OpNetworkChange, OpServerRestart} {
		if !RequiresApproval(op) {
			t.Fatalf("%q must require approval", op)
		}
	}
}
