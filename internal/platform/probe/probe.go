package probe

import (
	"context"
	"net"
	"strings"
	"time"
)

// Observation is a normalized measurement emitted by an FTN probe.
type Observation struct {
	ProbeID     string
	Target      string
	Kind        string
	Transport   string
	Success     bool
	Latency     time.Duration
	ErrorClass  string
	ObservedAt  time.Time
}

// DNSProbe performs a lightweight DNS resolution measurement.
type DNSProbe struct {
	Resolver *net.Resolver
	Timeout  time.Duration
}

func (p DNSProbe) Resolve(ctx context.Context, probeID, target string) Observation {
	started := time.Now()
	resolver := p.Resolver
	if resolver == nil {
		resolver = net.DefaultResolver
	}
	timeout := p.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := resolver.LookupHost(ctx, target)
	obs := Observation{ProbeID: probeID, Target: target, Kind: "dns", Transport: "system-resolver", Success: err == nil, Latency: time.Since(started), ObservedAt: time.Now().UTC()}
	if err != nil {
		obs.ErrorClass = classifyError(err)
	}
	return obs
}

func classifyError(err error) string {
	if err == nil {
		return ""
	}
	s := strings.ToLower(err.Error())
	switch {
	case strings.Contains(s, "timeout"):
		return "timeout"
	case strings.Contains(s, "no such host"):
		return "nxdomain"
	case strings.Contains(s, "refused"):
		return "refused"
	default:
		return "resolver_error"
	}
}
