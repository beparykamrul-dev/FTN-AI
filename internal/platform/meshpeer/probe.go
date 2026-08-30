package meshpeer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// Probe performs a bounded TCP/TLS connectivity check against an explicitly
// supplied endpoint. It does not create tunnels, mutate routes, or change
// firewall state.
func Probe(ctx context.Context, endpoint string, serverName string, caPool *tls.Config, timeout time.Duration) ProbeResult {
	if endpoint == "" {
		return ProbeResult{Error: "empty endpoint"}
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", endpoint)
	if err != nil {
		return ProbeResult{Error: fmt.Sprintf("tcp connect: %v", err)}
	}
	defer conn.Close()

	result := ProbeResult{Reachable: true, LatencyMS: time.Since(start).Milliseconds()}
	if caPool == nil {
		result.ReadOK = true
		return result
	}

	cfg := caPool.Clone()
	cfg.ServerName = serverName
	tlsConn := tls.Client(conn, cfg)
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		return ProbeResult{Reachable: true, LatencyMS: result.LatencyMS, Error: fmt.Sprintf("tls verify: %v", err)}
	}
	result.ReadOK = true
	return result
}
