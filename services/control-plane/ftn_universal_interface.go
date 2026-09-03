package main

import (
	"errors"
	"sort"
	"strings"
)

// FTNUniversalInterface describes the provider-neutral contract exposed by
// an authorized device to participate in the FTN service fabric.
type FTNUniversalInterface struct {
	DeviceID     string   `json:"device_id"`
	DriverID     string   `json:"driver_id"`
	Services     []string `json:"services"`
	Protocols    []string `json:"protocols"`
	Transports   []string `json:"transports"`
	Backbone     string   `json:"backbone"`
	Mesh         string   `json:"mesh"`
	Telemetry    string   `json:"telemetry"`
	FTNOwned     bool     `json:"ftn_owned"`
	Authorized   bool     `json:"authorized"`
}

// NormalizeFTNUniversalInterface validates the provider-neutral FTN boundary.
// It does not perform network configuration or privileged mutations.
func NormalizeFTNUniversalInterface(i FTNUniversalInterface) (FTNUniversalInterface, error) {
	i.DeviceID = strings.TrimSpace(i.DeviceID)
	i.DriverID = strings.TrimSpace(i.DriverID)
	i.Backbone = strings.ToLower(strings.TrimSpace(i.Backbone))
	i.Mesh = strings.ToLower(strings.TrimSpace(i.Mesh))
	i.Telemetry = strings.ToLower(strings.TrimSpace(i.Telemetry))
	if i.DeviceID == "" || i.DriverID == "" {
		return FTNUniversalInterface{}, errors.New("ftn_interface_identity_required")
	}
	if !i.FTNOwned && !i.Authorized {
		return FTNUniversalInterface{}, errors.New("ftn_interface_authorization_required")
	}
	if i.Backbone != "ftn-buckboon" {
		return FTNUniversalInterface{}, errors.New("ftn_buckboon_required")
	}
	if i.Mesh != "ftn-mesh" {
		return FTNUniversalInterface{}, errors.New("ftn_mesh_required")
	}
	if i.Telemetry != "silk-primary" {
		return FTNUniversalInterface{}, errors.New("silk_primary_required")
	}
	i.Services = normalizeList(i.Services)
	i.Protocols = normalizeList(i.Protocols)
	i.Transports = normalizeList(i.Transports)
	if len(i.Services) == 0 || len(i.Protocols) == 0 || len(i.Transports) == 0 {
		return FTNUniversalInterface{}, errors.New("ftn_interface_capabilities_required")
	}
	return i, nil
}

func normalizeList(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

// FTNInterfaceCanProvide reports whether a normalized interface can expose a
// requested FTN service using at least one compatible protocol and transport.
func FTNInterfaceCanProvide(i FTNUniversalInterface, service string, protocols, transports []string) bool {
	service = strings.ToLower(strings.TrimSpace(service))
	if service == "" || !contains(i.Services, service) {
		return false
	}
	return intersects(i.Protocols, protocols) && intersects(i.Transports, transports)
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func intersects(a, b []string) bool {
	for _, x := range a {
		x = strings.ToLower(strings.TrimSpace(x))
		for _, y := range b {
			if x == strings.ToLower(strings.TrimSpace(y)) {
				return true
			}
		}
	}
	return false
}
