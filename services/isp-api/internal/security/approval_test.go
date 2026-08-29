package security

import "testing"

func TestRequireApproval(t *testing.T) {
    gate := Gate{}
    if err := gate.RequireApproval("restart-service", false); err != ErrApprovalRequired {
        t.Fatalf("expected approval error, got %v", err)
    }
    if err := gate.RequireApproval("restart-service", true); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestInvalidAction(t *testing.T) {
    if err := (Gate{}).RequireApproval("", true); err != ErrInvalidAction {
        t.Fatalf("expected invalid action error, got %v", err)
    }
}
