package global

import (
	"context"
	"fmt"
	"net"
	"time"
)

// ResolverProbe describes a DNS endpoint that FTN may health-check.
type ResolverProbe struct {
	Name    string
	Address string
	Network string
}

// Result contains only operational health data; it must never contain secrets.
type Result struct {
	Name      string
	Address   string
	Reachable bool
	Latency   time.Duration
	Error     string
}

// Probe performs a lightweight transport reachability/latency check.
// DNS protocol validation remains separate so the transport layer stays small.
func Probe(ctx context.Context, p ResolverProbe, timeout time.Duration) Result {
	start := time.Now()
	result := Result{Name: p.Name, Address: p.Address}
	if p.Address == "" {
		result.Error = "resolver address is required"
		return result
	}
	if p.Network == "" {
		p.Network = "tcp"
	}
	if p.Network != "tcp" && p.Network != "udp" {
		result.Error = fmt.Sprintf("unsupported network: %s", p.Network)
		return result
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	if ctx == nil {
		ctx = context.Background()
	}
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, p.Network, p.Address)
	result.Latency = time.Since(start)
	if err != nil {
		result.Error = fmt.Sprintf("probe failed: %v", err)
		return result
	}
	if err := conn.Close(); err != nil {
		result.Error = fmt.Sprintf("probe close failed: %v", err)
		return result
	}
	result.Reachable = true
	return result
}

// RankByLatency returns probes ordered by measured latency, placing
// unreachable endpoints last. The input is copied so callers retain ownership.
func RankByLatency(results []Result) []Result {
	out := append([]Result(nil), results...)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0; j-- {
			left, right := out[j-1], out[j]
			if betterOrEqual(left, right) {
				break
			}
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}

func betterOrEqual(a, b Result) bool {
	if a.Reachable != b.Reachable {
		return a.Reachable
	}
	if !a.Reachable {
		return true
	}
	return a.Latency <= b.Latency
}
