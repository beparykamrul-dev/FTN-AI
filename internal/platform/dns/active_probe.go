package dns

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type DNSProbe struct {
	Address string `json:"address"`
	Name string `json:"name"`
	Timeout time.Duration `json:"timeout"`
}

type DNSProbeResult struct {
	Address string `json:"address"`
	Reachable bool `json:"reachable"`
	Latency time.Duration `json:"latency"`
	Error string `json:"error,omitempty"`
}

func (p DNSProbe) Validate() error {
	if strings.TrimSpace(p.Address) == "" { return fmt.Errorf("DNS probe address is required") }
	if p.Timeout <= 0 { return fmt.Errorf("DNS probe timeout must be positive") }
	return nil
}

// Probe performs a transport-level reachability check. DNS query semantics are
// intentionally injected by the caller so the probe can be backed by UDP/TCP
// or a DNS client without coupling the control plane to one resolver library.
func (p DNSProbe) Probe(ctx context.Context) DNSProbeResult {
	result := DNSProbeResult{Address: p.Address}
	if err := p.Validate(); err != nil { result.Error = err.Error(); return result }
	start := time.Now()
	d := net.Dialer{Timeout: p.Timeout}
	conn, err := d.DialContext(ctx, "udp", p.Address)
	result.Latency = time.Since(start)
	if err != nil { result.Error = err.Error(); return result }
	_ = conn.Close()
	result.Reachable = true
	return result
}
