package main

import "testing"

func TestValidateNetworkExecutionIntent(t *testing.T) {
	base := NetworkExecutionIntent{Device: NetworkDevice{ID: "r1", Protocol: "routeros-api", Healthy: true}, Action: NetworkRead}
	if err := ValidateNetworkExecutionIntent(base); err != nil { t.Fatal(err) }

	apply := base
	apply.Action = NetworkApply
	if err := ValidateNetworkExecutionIntent(apply); err == nil { t.Fatal("expected approval requirement") }

	apply.ApprovalID = "approval-1"
	apply.PreSnapshot = true
	if err := ValidateNetworkExecutionIntent(apply); err == nil { t.Fatal("expected post verification requirement") }

	apply.PostVerify = true
	if err := ValidateNetworkExecutionIntent(apply); err != nil { t.Fatal(err) }
}

func TestValidateFTNOwnership(t *testing.T) {
	device := NetworkDevice{ID: "r1", Protocol: "routeros-api"}
	if err := ValidateFTNOwnership(device, false); err == nil { t.Fatal("expected ownership failure") }
	if err := ValidateFTNOwnership(device, true); err != nil { t.Fatal(err) }
}
