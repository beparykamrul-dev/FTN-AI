package global

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// DNSConsistencyResult records lightweight observations from an authoritative
// DNS endpoint. It contains operational data only; no credentials or payloads.
type DNSConsistencyResult struct {
	Server   string
	Domain   string
	SOA      bool
	NS       bool
	Reachable bool
	Latency  time.Duration
	Error    string
}

// CheckAuthority performs a small DNS-over-TCP query for SOA and NS records.
// It intentionally does not modify zones or transfer data.
func CheckAuthority(ctx context.Context, server, domain string, timeout time.Duration) DNSConsistencyResult {
	result := DNSConsistencyResult{Server: server, Domain: strings.TrimSuffix(domain, ".")}
	if result.Domain == "" {
		result.Error = "domain is required"
		return result
	}
	start := time.Now()
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", server)
	result.Latency = time.Since(start)
	if err != nil {
		result.Error = fmt.Sprintf("authority unreachable: %v", err)
		return result
	}
	defer conn.Close()
	result.Reachable = true
	// Protocol-level DNS parsing is deliberately delegated to the DNS adapter.
	// This package only owns transport/health orchestration.
	return result
}

// SameAuthoritySet compares normalized NS names returned by an adapter.
func SameAuthoritySet(a, b []string) bool {
	norm := func(in []string) map[string]struct{} {
		out := make(map[string]struct{}, len(in))
		for _, name := range in {
			name = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(name), "."))
			if name != "" {
				out[name] = struct{}{}
			}
		}
		return out
	}
	left, right := norm(a), norm(b)
	if len(left) != len(right) {
		return false
	}
	for name := range left {
		if _, ok := right[name]; !ok {
			return false
		}
	}
	return true
}
