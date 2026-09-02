package main

import (
	"errors"
	"strings"
	"time"
)

var (
	errExternalTarget = errors.New("external target is prohibited")
	errOwnershipRequired = errors.New("FTN asset ownership is required")
	errApprovalRequired = errors.New("approval is required")
	errInvalidAction = errors.New("invalid network action")
)

type NetworkExecutionIntent struct {
	Device NetworkDevice
	Action string
	Approved bool
	Explicit bool
	PrechangeSnapshot bool
	VerificationRequired bool
	Timeout time.Duration
}

type NetworkExecutionDecision struct {
	Allowed bool `json:"allowed"`
	RequiresApproval bool `json:"requires_approval"`
	RequiresExplicitApproval bool `json:"requires_explicit_approval"`
	Reason string `json:"reason"`
}

// EvaluateNetworkExecution is the final policy boundary before an adapter may
// touch an FTN-owned device. It deliberately contains no transport-specific
// side effects and never executes shell commands.
func EvaluateNetworkExecution(in NetworkExecutionIntent) NetworkExecutionDecision {
	if strings.TrimSpace(in.Device.ID) == "" || !in.Device.Healthy {
		return NetworkExecutionDecision{Reason: "device_invalid_or_unhealthy"}
	}
	if strings.TrimSpace(in.Device.Address) == "" {
		return NetworkExecutionDecision{Reason: "device_address_required"}
	}
	if !isFTNDeviceKind(in.Device.Kind) {
		return NetworkExecutionDecision{Reason: "target_ownership_not_verified"}
	}
	action := strings.ToLower(strings.TrimSpace(in.Action))
	if action == "" {
		return NetworkExecutionDecision{Reason: "action_required"}
	}
	if isExternalNetworkAction(action) {
		return NetworkExecutionDecision{RequiresExplicitApproval: true, Reason: errExternalTarget.Error()}
	}
	if isDestructiveNetworkAction(action) {
		return NetworkExecutionDecision{RequiresExplicitApproval: true, Reason: "destructive_network_action"}
	}
	if !in.PrechangeSnapshot || !in.VerificationRequired {
		return NetworkExecutionDecision{Reason: "snapshot_and_verification_required"}
	}
	if in.Timeout <= 0 || in.Timeout > 2*time.Minute {
		return NetworkExecutionDecision{Reason: "execution_timeout_out_of_bounds"}
	}
	if requiresNetworkApproval(action) && !in.Approved {
		return NetworkExecutionDecision{RequiresApproval: true, Reason: errApprovalRequired.Error()}
	}
	return NetworkExecutionDecision{Allowed: true, Reason: "approved_ftn_owned_execution"}
}

func isFTNDeviceKind(kind string) bool {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "router", "core-router", "edge-router", "mikrotik", "switch", "olt", "onu", "dns", "server", "firewall":
		return true
	default:
		return false
	}
}

func requiresNetworkApproval(action string) bool {
	for _, prefix := range []string{"configure", "route", "firewall", "credential", "service", "vlan", "pppoe", "bgp", "ospf", "bfd"} {
		if strings.HasPrefix(action, prefix) { return true }
	}
	return false
}

func isDestructiveNetworkAction(action string) bool {
	for _, word := range []string{"delete", "wipe", "factory-reset", "destroy"} {
		if strings.Contains(action, word) { return true }
	}
	return false
}

func isExternalNetworkAction(action string) bool {
	for _, word := range []string{"external", "retaliate", "exploit", "scan-external"} {
		if strings.Contains(action, word) { return true }
	}
	return false
}
