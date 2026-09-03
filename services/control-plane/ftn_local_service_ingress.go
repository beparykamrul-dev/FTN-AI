package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/netip"
	"strings"
)

type FTNLocalServiceIngress struct {
	RouterID string `json:"router_id"`
	Tunnel string `json:"tunnel"`
	SourcePrefix string `json:"source_prefix"`
	ServicePrefix string `json:"service_prefix"`
	Service string `json:"service"`
	SILKExport bool `json:"silk_export"`
	FirewallEnforced bool `json:"firewall_enforced"`
	Approved bool `json:"approved"`
}

func NormalizeFTNLocalServiceIngress(in FTNLocalServiceIngress) (FTNLocalServiceIngress, error) {
	in.RouterID = strings.TrimSpace(in.RouterID)
	in.Tunnel = strings.ToLower(strings.TrimSpace(in.Tunnel))
	in.SourcePrefix = strings.TrimSpace(in.SourcePrefix)
	in.ServicePrefix = strings.TrimSpace(in.ServicePrefix)
	in.Service = strings.TrimSpace(in.Service)
	if in.RouterID == "" || in.Tunnel == "" || in.Service == "" { return FTNLocalServiceIngress{}, fmt.Errorf("router, tunnel and service are required") }
	if in.Tunnel != "wireguard" && in.Tunnel != "gre" && in.Tunnel != "vxlan" && in.Tunnel != "ipsec" { return FTNLocalServiceIngress{}, fmt.Errorf("unsupported local ingress tunnel") }
	if in.SourcePrefix != "" { p, err := netip.ParsePrefix(in.SourcePrefix); if err != nil { return FTNLocalServiceIngress{}, fmt.Errorf("invalid source prefix: %w", err) }; in.SourcePrefix = p.Masked().String() }
	if in.ServicePrefix != "" { p, err := netip.ParsePrefix(in.ServicePrefix); if err != nil { return FTNLocalServiceIngress{}, fmt.Errorf("invalid service prefix: %w", err) }; in.ServicePrefix = p.Masked().String() }
	if !in.Approved { return FTNLocalServiceIngress{}, fmt.Errorf("local ingress approval required") }
	if !in.FirewallEnforced { return FTNLocalServiceIngress{}, fmt.Errorf("firewall enforcement required") }
	return in, nil
}

func FTNLocalIngressHash(in FTNLocalServiceIngress) string {
	v := strings.Join([]string{in.RouterID,in.Tunnel,in.SourcePrefix,in.ServicePrefix,in.Service,fmt.Sprint(in.SILKExport),fmt.Sprint(in.FirewallEnforced)}, "|")
	s := sha256.Sum256([]byte(v))
	return hex.EncodeToString(s[:])
}
