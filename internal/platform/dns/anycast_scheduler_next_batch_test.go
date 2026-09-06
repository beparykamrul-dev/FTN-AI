package dns

import (
 "context"
 "testing"
 "time"
)

func TestAnycastScheduleAndRunBoundaries(t *testing.T) {
 s := AnycastSchedule{Interval: time.Hour, Probe: DNSQueryProbe{Address: "127.0.0.1:53", Name: "example.com", Timeout: time.Second}, Policy: ProbePolicy{FailureThreshold: 1, SuccessThreshold: 1, MaxLatency: time.Second}}
 if err := s.Validate(); err != nil { t.Fatalf("valid schedule rejected: %v", err) }
 if err := (AnycastSchedule{}).Validate(); err == nil { t.Fatal("empty schedule accepted") }
 if err := NewAnycastScheduler().Run(nil, s, func(context.Context, DNSProbeResult) error { return nil }); err == nil { t.Fatal("nil context accepted") }
}
