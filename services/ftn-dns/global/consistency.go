package global

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type DNSConsistencyResult struct {
	Server string
	Domain string
	SOA bool
	NS bool
	Reachable bool
	Latency time.Duration
	Error string
}

func CheckAuthority(ctx context.Context, server, domain string, timeout time.Duration) DNSConsistencyResult {
	result := DNSConsistencyResult{Server: strings.TrimSpace(server), Domain: strings.TrimSuffix(strings.TrimSpace(domain), ".")}
	if result.Domain == "" { result.Error = "domain is required"; return result }
	if result.Server == "" { result.Error = "server is required"; return result }
	if ctx == nil { result.Error = "context is required"; return result }
	if timeout <= 0 { timeout = 3 * time.Second }
	start := time.Now()
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", result.Server)
	result.Latency = time.Since(start)
	if err != nil { result.Error = fmt.Sprintf("authority unreachable: %v", err); return result }
	defer conn.Close()
	result.Reachable = true
	return result
}

func SameAuthoritySet(a, b []string) bool {
	norm := func(in []string) map[string]struct{} {
		out := make(map[string]struct{}, len(in))
		for _, name := range in {
			name = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(name), "."))
			if name != "" { out[name] = struct{}{} }
		}
		return out
	}
	left, right := norm(a), norm(b)
	if len(left) != len(right) { return false }
	for name := range left { if _, ok := right[name]; !ok { return false } }
	return true
}
