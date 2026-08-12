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

// Probe performs a lightweight TCP reachability/latency check. DNS protocol
// validation is intentionally kept separate so the transport layer stays small.
func Probe(ctx context.Context, p ResolverProbe, timeout time.Duration) Result {
	start := time.Now()
	result := Result{Name: p.Name, Address: p.Address}

	if p.Network == "" {
		p.Network = "tcp"
	}

	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, p.Network, p.Address)
	result.Latency = time.Since(start)
	if err != nil {
		result.Error = fmt.Sprintf("probe failed: %v", err)
		return result
	}
	_ = conn.Close()
	result.Reachable = true
	return result
}
