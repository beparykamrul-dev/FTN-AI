package main

import (
	"strings"
	"testing"
	"time"
)

func TestNormalizeFlowObservationUTC(t *testing.T) {
	f := FlowRecord{ExporterID: "192.0.2.1", Version: 10, SourceIP: "192.0.2.10", DestinationIP: "198.51.100.20", Bytes: 10, Packets: 1, SamplingRate: 1}
	observed := time.Date(2026, 9, 4, 12, 0, 0, 0, time.FixedZone("test", 6*60*60))
	o, err := normalizeFlowObservation(f, observed, time.Time{}, time.Time{})
	if err != nil { t.Fatalf("normalize: %v", err) }
	if o.ObservedAt.Location() != time.UTC { t.Fatalf("expected UTC, got %v", o.ObservedAt.Location()) }
	if o.ObservedAt.Hour() != 6 { t.Fatalf("unexpected UTC conversion: %v", o.ObservedAt) }
}

func TestFlowObservationRejectsInvalidTimeRange(t *testing.T) {
	f := FlowRecord{ExporterID: "192.0.2.1", Version: 10, SourceIP: "192.0.2.10", DestinationIP: "198.51.100.20", SamplingRate: 1}
	_, err := normalizeFlowObservation(f, time.Unix(100, 0), time.Unix(200, 0), time.Unix(100, 0))
	if err == nil || !strings.Contains(err.Error(), "invalid_flow_time_range") { t.Fatalf("expected invalid time range, got %v", err) }
}

func TestFlowObservationRequiresObservedAt(t *testing.T) {
	o := FlowObservation{FlowRecord: FlowRecord{ExporterID: "192.0.2.1", Version: 10, SourceIP: "192.0.2.10", DestinationIP: "198.51.100.20", SamplingRate: 1}}
	if err := o.Validate(); err == nil || !strings.Contains(err.Error(), "observed_at_required") { t.Fatalf("expected observed_at_required, got %v", err) }
}
