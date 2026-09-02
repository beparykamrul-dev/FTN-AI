package main

import "testing"

func TestNormalizeProtocol(t *testing.T) {
	if got := normalizeProtocol("  RouterOS-API "); got != "routeros-api" { t.Fatalf("got %q", got) }
}

func TestNetworkProtocolAllowlist(t *testing.T) {
	for _, protocol := range []string{"routeros-api", "snmp", "netflow", "ipfix", "bgp", "ospf", "bfd", "vrf", "ecmp", "olt", "onu"} {
		if !validNetworkProtocol(protocol) { t.Fatalf("protocol %q rejected", protocol) }
	}
	if validNetworkProtocol("telnet") { t.Fatal("unsafe/unsupported protocol accepted") }
}
