package main

import (
	"errors"
	"strings"
	"time"
)

// ValidateNetworkExecutionIntent validates the common execution contract used
// by all FTN network adapters. It performs no network side effects.
func ValidateNetworkExecutionIntent(i NetworkExecutionIntent) error {
	if strings.TrimSpace(i.Device.ID) == "" {
		return errors.New("device_id_required")
	}
	if strings.TrimSpace(i.Device.Protocol) == "" {
		return errors.New("protocol_required")
	}
	if !i.Device.Healthy && !isReadOnlyNetworkAction(i.Action) {
		return errors.New("unhealthy_device_requires_read_or_health")
	}
	if isMutationNetworkAction(i.Action) {
		if !i.Approved {
			return errors.New("approval_required")
		}
		if !i.PrechangeSnapshot || !i.VerificationRequired || !i.RollbackSafe {
			return errors.New("snapshot_verification_and_rollback_required")
		}
		if i.Timeout <= 0 || i.Timeout > 2*time.Minute {
			return errors.New("execution_timeout_out_of_bounds")
		}
	}
	return nil
}

func isReadOnlyNetworkAction(action string) bool {
	a := strings.ToLower(strings.TrimSpace(action))
	return a == "read" || a == "health" || a == "plan" || strings.HasSuffix(a, ".read")
}

func isMutationNetworkAction(action string) bool {
	a := strings.ToLower(strings.TrimSpace(action))
	return a == "apply" || a == "configuration" || strings.HasPrefix(a, "configure") || strings.HasPrefix(a, "route") || strings.HasPrefix(a, "firewall") || strings.HasPrefix(a, "vlan") || strings.HasPrefix(a, "pppoe") || strings.HasPrefix(a, "bgp") || strings.HasPrefix(a, "ospf") || strings.HasPrefix(a, "bfd") || strings.HasPrefix(a, "service")
}

// ValidateFTNOwnership prevents an adapter from acting on an unverified target.
// Ownership is represented by the control-plane device inventory; adapters never
// infer ownership from an address, ASN, or hostname.
func ValidateFTNOwnership(device NetworkDevice, owned bool) error {
	if !owned {
		return errors.New("ftn_ownership_not_verified")
	}
	if strings.TrimSpace(device.ID) == "" {
		return errors.New("device_id_required")
	}
	return nil
}
