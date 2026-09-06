package control

import "testing"

func TestOperationAllowlistAndApproval(t *testing.T) {
 if err := Validate(OpReadHealth); err != nil { t.Fatalf("read operation rejected: %v", err) }
 if err := Validate(Operation(" unknown ")); err == nil { t.Fatal("unknown operation accepted") }
 if !RequiresApproval(OpDNSWrite) || !RequiresApproval(OpNetworkChange) || !RequiresApproval(OpServerRestart) { t.Fatal("privileged operations not approval-gated") }
 if RequiresApproval(OpReadHealth) { t.Fatal("read operation incorrectly requires approval") }
}
