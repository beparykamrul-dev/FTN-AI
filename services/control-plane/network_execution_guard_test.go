package main

import (
	"testing"
	"time"
)

func TestEvaluateNetworkExecution(t *testing.T) {
	base := NetworkExecutionIntent{
		Device:               NetworkDevice{ID: "core-a", Kind: "core-router", Address: "10.0.0.1", Healthy: true},
		Action:               "configure-interface",
		ApprovalID:           "approval-1",
		Approved:             true,
		Explicit:             true,
		PrechangeSnapshot:    true,
		VerificationRequired: true,
		RollbackSafe:         true,
		Timeout:              30 * time.Second,
	}
	if got := EvaluateNetworkExecution(base); !got.Allowed {
		t.Fatalf("expected allowed execution, got %+v", got)
	}

	base.Approved = false
	if got := EvaluateNetworkExecution(base); !got.RequiresApproval || got.Allowed {
		t.Fatalf("expected approval gate, got %+v", got)
	}

	base.Approved = true
	base.Action = "delete-route"
	if got := EvaluateNetworkExecution(base); !got.RequiresExplicitApproval || got.Allowed {
		t.Fatalf("expected explicit approval gate, got %+v", got)
	}

	base.Action = "retaliate-external-host"
	if got := EvaluateNetworkExecution(base); !got.RequiresExplicitApproval || got.Allowed {
		t.Fatalf("expected external-action prohibition, got %+v", got)
	}

	base.Action = "configure-interface"
	base.Device.Kind = "unknown"
	if got := EvaluateNetworkExecution(base); got.Allowed {
		t.Fatalf("expected ownership rejection, got %+v", got)
	}
}
