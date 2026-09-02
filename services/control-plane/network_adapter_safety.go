package main

import (
	"errors"
	"strings"
)

type NetworkActionClass string

const (
	NetworkRead        NetworkActionClass = "read"
	NetworkHealth      NetworkActionClass = "health"
	NetworkPlan        NetworkActionClass = "plan"
	NetworkApply       NetworkActionClass = "apply"
	NetworkRollback    NetworkActionClass = "rollback"
	NetworkConfiguration NetworkActionClass = "configuration"
)

type NetworkExecutionIntent struct {
	Device       NetworkDevice       `json:"device"`
	Action       NetworkActionClass  `json:"action"`
	ChangeID     string              `json:"change_id,omitempty"`
	ApprovalID   string              `json:"approval_id,omitempty"`
	PreSnapshot  bool                `json:"pre_snapshot"`
	PostVerify   bool                `json:"post_verify"`
	RollbackSafe bool                `json:"rollback_safe"`
}

func ValidateNetworkExecutionIntent(i NetworkExecutionIntent) error {
	if strings.TrimSpace(i.Device.ID) == "" {
		return errors.New("device_id_required")
	}
	if strings.TrimSpace(i.Device.Protocol) == "" {
		return errors.New("protocol_required")
	}
	if !i.Device.Healthy && i.Action != NetworkRead && i.Action != NetworkHealth && i.Action != NetworkPlan {
		return errors.New("unhealthy_device_requires_plan_or_read")
	}
	if i.Action == NetworkApply || i.Action == NetworkRollback || i.Action == NetworkConfiguration {
		if strings.TrimSpace(i.ApprovalID) == "" {
			return errors.New("approval_required")
		}
		if !i.PreSnapshot || !i.PostVerify {
			return errors.New("snapshot_and_post_verification_required")
		}
	}
	return nil
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
