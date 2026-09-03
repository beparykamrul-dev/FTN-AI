package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
)

// FTNDeviceDriver is the provider-neutral boundary between the control plane
// and FTN-owned devices. Read-only discovery/probing is automatic; privileged
// configuration changes remain approval-gated by the caller.
type FTNDeviceDriver interface {
	Identity() DeviceDriverIdentity
	Capabilities() []string
	Probe(DeviceDriverTarget) (DeviceDriverObservation, error)
	Validate(DeviceDriverTarget) error
}

type DeviceDriverIdentity struct {
	ID string `json:"id"`
	Vendor string `json:"vendor"`
	Family string `json:"family"`
	Version string `json:"version"`
}

type DeviceDriverTarget struct {
	DeviceID string `json:"device_id"`
	Address string `json:"address"`
	Class string `json:"class"`
	Model string `json:"model"`
	Serial string `json:"serial,omitempty"`
}

type DeviceDriverObservation struct {
	DeviceID string `json:"device_id"`
	DriverID string `json:"driver_id"`
	Reachable bool `json:"reachable"`
	Health string `json:"health"`
	Capabilities []string `json:"capabilities"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Fingerprint string `json:"fingerprint"`
}

func NormalizeDeviceDriverTarget(t DeviceDriverTarget) (DeviceDriverTarget, error) {
	t.DeviceID = strings.TrimSpace(t.DeviceID)
	t.Address = strings.TrimSpace(t.Address)
	t.Class = strings.ToLower(strings.TrimSpace(t.Class))
	t.Model = strings.TrimSpace(t.Model)
	t.Serial = strings.TrimSpace(t.Serial)
	if t.DeviceID == "" || t.Address == "" || t.Class == "" {
		return DeviceDriverTarget{}, errors.New("device_driver_target_required")
	}
	allowed := map[string]bool{"router": true, "switch": true, "olt": true, "onu": true, "server": true, "storage": true, "camera": true, "tv": true, "mobile": true, "embedded": true, "virtual_node": true}
	if !allowed[t.Class] {
		return DeviceDriverTarget{}, errors.New("unsupported_device_class")
	}
	return t, nil
}

func DeviceDriverFingerprint(t DeviceDriverTarget, driver DeviceDriverIdentity) string {
	parts := []string{strings.ToLower(strings.TrimSpace(t.DeviceID)), strings.ToLower(strings.TrimSpace(t.Class)), strings.ToLower(strings.TrimSpace(t.Model)), strings.ToLower(strings.TrimSpace(t.Serial)), strings.ToLower(strings.TrimSpace(driver.ID)), strings.ToLower(strings.TrimSpace(driver.Version))}
	h := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(h[:])
}

func NormalizeDeviceDriverObservation(o DeviceDriverObservation) (DeviceDriverObservation, error) {
	o.DeviceID = strings.TrimSpace(o.DeviceID)
	o.DriverID = strings.TrimSpace(o.DriverID)
	o.Health = strings.ToLower(strings.TrimSpace(o.Health))
	if o.DeviceID == "" || o.DriverID == "" {
		return DeviceDriverObservation{}, errors.New("device_driver_observation_required")
	}
	if o.Health == "" {
		o.Health = "unknown"
	}
	o.Capabilities = append([]string(nil), o.Capabilities...)
	sort.Strings(o.Capabilities)
	return o, nil
}
