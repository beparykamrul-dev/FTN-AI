package main

import "sort"

var networkCapabilityCatalog = map[string][]string{
	"routeros-api": {"identity", "interfaces", "routes", "vlans", "pppoe", "health"},
	"snmp":         {"identity", "interfaces", "counters", "health"},
	"olt":          {"identity", "pon", "onu", "optical", "health"},
	"onu":          {"identity", "optical", "ethernet", "health"},
	"bgp":          {"neighbors", "routes", "health"},
	"ospf":         {"neighbors", "routes", "health"},
	"bfd":          {"sessions", "health"},
	"vrf":          {"tables", "routes", "health"},
	"ecmp":         {"paths", "health"},
	"netflow":      {"flows"},
	"ipfix":        {"flows"},
}

func networkCapabilities(protocol string) []string {
	items := append([]string(nil), networkCapabilityCatalog[normalizeProtocol(protocol)]...)
	sort.Strings(items)
	return items
}
