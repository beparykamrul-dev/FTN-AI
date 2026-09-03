package main

import "testing"

func TestFlowRecordFingerprintStable(t *testing.T) {
	a := FlowRecord{ExporterID: "r1", SourceIP: "10.0.0.1", DestinationIP: "10.0.0.2", SourcePort: 1234, DestinationPort: 443, Protocol: 6, Bytes: 100, Packets: 2, SamplingRate: 1}
	if FlowRecordFingerprint(a) != FlowRecordFingerprint(a) { t.Fatal("fingerprint must be deterministic") }
}

func TestNormalizeSampledCounters(t *testing.T) {
	r := NormalizeSampledCounters(FlowRecord{Bytes: 100, Packets: 2, SamplingRate: 10})
	if r.Bytes != 1000 || r.Packets != 20 { t.Fatalf("unexpected sampled counters: %+v", r) }
}
