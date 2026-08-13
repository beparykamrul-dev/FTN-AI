package adapters

import (
	"context"
	"fmt"
	"net"
	"time"

	ftndns "github.com/beparykamrul-dev/FTN-AI/services/ftn-dns"
)

// PowerDNSAdapter provides the first provider-neutral PowerDNS transport
// adapter. Management credentials are intentionally not stored here.
type PowerDNSAdapter struct {
	Endpoint string
	Timeout  time.Duration
}

func (p PowerDNSAdapter) Name() string { return "powerdns" }

func (p PowerDNSAdapter) Health(ctx context.Context) (ftndns.Health, error) {
	timeout := p.Timeout
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	start := time.Now()
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", p.Endpoint)
	latency := time.Since(start)
	if err != nil {
		return ftndns.Health{LatencyMS: latency.Milliseconds()}, fmt.Errorf("powerdns endpoint unreachable: %w", err)
	}
	_ = conn.Close()
	return ftndns.Health{Reachable: true, LatencyMS: latency.Milliseconds(), LossRatio: 0}, nil
}

func (p PowerDNSAdapter) Query(context.Context, string, string) (ftndns.Response, error) {
	return ftndns.Response{}, fmt.Errorf("powerdns query adapter requires the configured DNS transport")
}
