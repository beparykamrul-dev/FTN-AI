package dns

import (
	"context"
	"fmt"
	"net"
	"time"
)

type ProtocolLatency struct {
	Protocol string `json:"protocol"`
	RTT time.Duration `json:"rtt"`
	Success bool `json:"success"`
	Error string `json:"error,omitempty"`
}

// ProbeProtocolLatency measures connection setup latency for a DNS endpoint.
// It intentionally avoids pretending that a TCP/UDP socket-open is equivalent
// to a complete DNS transaction; DNSQueryProbe remains the application-level
// correctness check.
func ProbeProtocolLatency(ctx context.Context, protocol, address string, timeout time.Duration) ProtocolLatency {
	result := ProtocolLatency{Protocol: protocol}
	if protocol != "tcp" && protocol != "udp" { result.Error = fmt.Sprintf("unsupported protocol: %s", protocol); return result }
	start := time.Now()
	conn, err := (&net.Dialer{Timeout: timeout}).DialContext(ctx, protocol, address)
	result.RTT = time.Since(start)
	if err != nil { result.Error = err.Error(); return result }
	_ = conn.Close()
	result.Success = true
	return result
}
