package main

import (
	"net"
	"testing"
	"time"
)

func TestSummarizeTrafficProbe(t *testing.T) {
	now := time.Now().UTC()
	state := trafficProbeState{Samples: []trafficProbeSample{
		{Latency: 20 * time.Millisecond, Success: true, At: now.Add(-3 * time.Second)},
		{Latency: 30 * time.Millisecond, Success: true, At: now.Add(-2 * time.Second)},
		{Latency: 40 * time.Millisecond, Success: false, At: now.Add(-1 * time.Second)},
		{Latency: 40 * time.Millisecond, Success: true, At: now},
	}, ConsecutiveFailure: 0}
	o := summarizeTrafficProbe(TrafficProbeTarget{ServiceID: "pubg", PathID: "path-a", Address: "127.0.0.1:1"}, state, now)
	if !o.Healthy || o.PacketLoss != 25 { t.Fatalf("unexpected health/loss: %+v", o) }
	if o.LatencyMs != 30 { t.Fatalf("unexpected latency: %v", o.LatencyMs) }
	if o.JitterMs <= 0 { t.Fatalf("expected jitter: %v", o.JitterMs) }
}

func TestSummarizeTrafficProbeMarksThreeConsecutiveFailuresUnhealthy(t *testing.T) {
	now := time.Now().UTC()
	state := trafficProbeState{Samples: []trafficProbeSample{{Latency: time.Millisecond, Success: true, At: now.Add(-3 * time.Second)}, {Success: false, At: now.Add(-2 * time.Second)}, {Success: false, At: now.Add(-time.Second)}, {Success: false, At: now}}, ConsecutiveFailure: 3}
	o := summarizeTrafficProbe(TrafficProbeTarget{ServiceID: "freefire", PathID: "path-b", Address: "127.0.0.1:1"}, state, now)
	if o.Healthy { t.Fatalf("expected unhealthy observation: %+v", o) }
}

func TestTrafficProbeCanReachLocalTCPService(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil { t.Fatal(err) }
	defer ln.Close()
	runtime := &TrafficRuntime{quality: NewTrafficQualityStore()}
	probe := NewTrafficActiveProbe(runtime, time.Hour, time.Second)
	probe.probeTarget(t.Context(), TrafficProbeTarget{ServiceID: "telegram", PathID: "local", Address: ln.Addr().String()})
	got := runtime.QualitySnapshot("telegram", time.Now().UTC())
	if len(got) != 1 || !got[0].Healthy || got[0].LatencyMs <= 0 { t.Fatalf("unexpected quality: %+v", got) }
}
