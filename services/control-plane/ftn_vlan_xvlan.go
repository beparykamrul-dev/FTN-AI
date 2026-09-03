package main

import (
	"fmt"
	"net/netip"
	"sort"
	"strings"
)

type FTNVLANMode string
const (
	FTNVLANAccess FTNVLANMode = "access"
	FTNVLANTrunk FTNVLANMode = "trunk"
	FTNVLANQinQ FTNVLANMode = "qinq"
)

type FTNVLAN struct {
	ID uint16 `json:"id"`
	Name string `json:"name"`
	Mode FTNVLANMode `json:"mode"`
	OuterID uint16 `json:"outer_id,omitempty"`
	InnerID uint16 `json:"inner_id,omitempty"`
	MTU uint16 `json:"mtu,omitempty"`
	Enabled bool `json:"enabled"`
}

type FTNVLANBinding struct {
	VLANID uint16 `json:"vlan_id"`
	DeviceID string `json:"device_id"`
	Interface string `json:"interface"`
	Service string `json:"service"`
	IPPrefix string `json:"ip_prefix,omitempty"`
}

func NormalizeFTNVLAN(v FTNVLAN) (FTNVLAN, error) {
	if v.ID == 0 || v.ID > 4094 { return FTNVLAN{}, fmt.Errorf("invalid VLAN id") }
	v.Name = strings.TrimSpace(v.Name)
	v.Mode = FTNVLANMode(strings.ToLower(strings.TrimSpace(string(v.Mode))))
	if v.Mode == "" { v.Mode = FTNVLANAccess }
	switch v.Mode { case FTNVLANAccess, FTNVLANTrunk: case FTNVLANQinQ:
		if v.OuterID == 0 || v.OuterID > 4094 || v.InnerID == 0 || v.InnerID > 4094 { return FTNVLAN{}, fmt.Errorf("invalid QinQ VLAN ids") }
	default: return FTNVLAN{}, fmt.Errorf("unsupported VLAN mode") }
	if v.MTU != 0 && (v.MTU < 576 || v.MTU > 9216) { return FTNVLAN{}, fmt.Errorf("invalid VLAN MTU") }
	return v, nil
}

func NormalizeFTNVLANBinding(b FTNVLANBinding) (FTNVLANBinding, error) {
	if b.VLANID == 0 || b.VLANID > 4094 { return FTNVLANBinding{}, fmt.Errorf("invalid VLAN binding id") }
	b.DeviceID = strings.TrimSpace(b.DeviceID); b.Interface = strings.TrimSpace(b.Interface); b.Service = strings.TrimSpace(b.Service)
	if b.DeviceID == "" || b.Interface == "" || b.Service == "" { return FTNVLANBinding{}, fmt.Errorf("device, interface and service are required") }
	b.IPPrefix = strings.TrimSpace(b.IPPrefix)
	if b.IPPrefix != "" { p, err := netip.ParsePrefix(b.IPPrefix); if err != nil { return FTNVLANBinding{}, fmt.Errorf("invalid IP prefix: %w", err) }; b.IPPrefix = p.Masked().String() }
	return b, nil
}

func SortFTNVLANs(v []FTNVLAN) []FTNVLAN {
	out := append([]FTNVLAN(nil), v...)
	sort.Slice(out, func(i,j int) bool { return out[i].ID < out[j].ID })
	return out
}
