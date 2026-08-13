package adapters

import (
	"context"
	"fmt"
	"net"
	"time"

	ftndns "github.com/beparykamrul-dev/FTN-AI/services/ftn-dns"
)

// GoDNSAdapter provides the lightweight FTN-native GoDNS transport adapter.
type GoDNSAdapter struct {
	Endpoint string
	Timeout  time.Duration
}

func (g GoDNSAdapter) Name() string { return "godns" }

func (g GoDNSAdapter) Health(ctx context.Context) (ftndns.Health, error) {
	timeout := g.Timeout
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	start := time.Now()
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", g.Endpoint)
	latency := time.Since(start)
	if err != nil {
		return ftndns.Health{LatencyMS: latency.Milliseconds()}, fmt.Errorf("godns endpoint unreachable: %w", err)
	}
	_ = conn.Close()
	return ftndns.Health{Reachable: true, LatencyMS: latency.Milliseconds(), LossRatio: 0}, nil
}

func (g GoDNSAdapter) Query(context.Context, string, string) (ftndns.Response, error) {
	return ftndns.Response{}, fmt.Errorf("godns query adapter requires the configured DNS transport")
}
