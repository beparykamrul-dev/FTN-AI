package main

import "testing"

func TestValidateRouterID(t *testing.T) {
	for _, tc := range []struct {
		name string
		value string
		want bool
	}{
		{"valid", "192.0.2.1", true},
		{"ipv6", "2001:db8::1", false},
		{"unspecified", "0.0.0.0", false},
		{"invalid", "not-an-ip", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := validateRouterID(tc.value) == nil; got != tc.want {
				t.Fatalf("validateRouterID(%q) valid=%v, want %v", tc.value, got, tc.want)
			}
		})
	}
}

func TestValidatePeerAddress(t *testing.T) {
	for _, tc := range []struct {
		name string
		value string
		want bool
	}{
		{"ipv4", "198.51.100.2", true},
		{"ipv6", "2001:db8::2", true},
		{"unspecified", "::", false},
		{"multicast", "224.0.0.1", false},
		{"invalid", "bad-peer", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := validatePeerAddress(tc.value) == nil; got != tc.want {
				t.Fatalf("validatePeerAddress(%q) valid=%v, want %v", tc.value, got, tc.want)
			}
		})
	}
}

func TestValidateIPv4Advertisement(t *testing.T) {
	for _, tc := range []struct {
		name   string
		prefix string
		nextHop string
		want   bool
	}{
		{"/24", "198.51.100.0/24", "192.0.2.1", true},
		{"/16", "198.51.0.0/16", "192.0.2.1", true},
		{"too-specific", "198.51.100.0/25", "192.0.2.1", false},
		{"ipv6-prefix", "2001:db8::/32", "192.0.2.1", false},
		{"ipv6-next-hop", "198.51.100.0/24", "2001:db8::1", false},
		{"zero-next-hop", "198.51.100.0/24", "0.0.0.0", false},
		{"multicast-next-hop", "198.51.100.0/24", "224.0.0.1", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := validateIPv4Advertisement(tc.prefix, tc.nextHop)
			if got := err == nil; got != tc.want {
				t.Fatalf("validateIPv4Advertisement(%q, %q) valid=%v, want %v; err=%v", tc.prefix, tc.nextHop, got, tc.want, err)
			}
		})
	}
}
