package main

import (
	"errors"
	"strings"
)

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
