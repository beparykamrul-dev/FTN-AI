package dns

import (
 "context"
 "testing"
 "time"
)

func TestProbePolicyValidationAndTransitions(t *testing.T) {
 p := ProbePolicy{FailureThreshold: 2, SuccessThreshold: 2, MaxLatency: time.Second}
 if err := p.Validate(); err != nil { t.Fatalf("valid policy rejected: %v", err) }
 if err := (ProbePolicy{FailureThreshold: 0, SuccessThreshold: 1, MaxLatency: time.Second}).Validate(); err == nil { t.Fatal("zero failure threshold accepted") }
 c := NewProbePolicyController(); ctx := context.Background()
 healthy, err := c.Observe(ctx, "dns-1", DNSProbeResult{Reachable: true, Latency: time.Millisecond}, p); if err != nil || healthy { t.Fatalf("first success healthy=%v err=%v", healthy, err) }
 healthy, err = c.Observe(ctx, "dns-1", DNSProbeResult{Reachable: true, Latency: time.Millisecond}, p); if err != nil || !healthy { t.Fatalf("second success healthy=%v err=%v", healthy, err) }
 if _, err := c.Observe(nil, "dns-1", DNSProbeResult{}, p); err == nil { t.Fatal("nil context accepted") }
}
