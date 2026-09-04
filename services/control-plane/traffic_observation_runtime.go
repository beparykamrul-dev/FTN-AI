package main

import (
	"context"
	"log"
	"time"
)

// IngestObservations is the timestamp-preserving ingestion path for imported
// flow sources such as nfcapd. Live collectors may continue using Ingest().
func (t *TrafficRuntime) IngestObservations(ctx context.Context, observations []FlowObservation) int {
	if t == nil || len(observations) == 0 { return 0 }
	if ctx == nil { ctx = context.Background() }
	accepted := make([]TrafficFlowObservation, 0, len(observations))
	for _, o := range observations {
		if err := o.Validate(); err != nil { continue }
		obs, ok := t.Classify(o.FlowRecord, o.ObservedAt)
		if !ok { continue }
		obs.ObservedAt = o.ObservedAt.UTC()
		accepted = append(accepted, obs)
	}
	if len(accepted) == 0 { return 0 }

	now := time.Now().UTC()
	t.mu.Lock()
	silk := t.silk
	t.flows = append(t.flows, accepted...)
	cutoff := now.Add(-2 * time.Minute)
	keep := t.flows[:0]
	for _, f := range t.flows { if f.ObservedAt.After(cutoff) { keep = append(keep, f) } }
	t.flows = keep
	if len(t.flows) > 4096 { t.flows = t.flows[len(t.flows)-4096:] }
	t.mu.Unlock()

	if silk != nil {
		byService := make(map[string][]FlowObservation)
		for _, f := range accepted {
			o, err := normalizeFlowObservation(f.FlowRecord, f.ObservedAt, time.Time{}, time.Time{})
			if err != nil { continue }
			byService[f.ServiceID] = append(byService[f.ServiceID], o)
		}
		for serviceID, batch := range byService {
			if err := silk.InsertBatchObservations(ctx, batch, serviceID, "", ""); err != nil {
				log.Printf("FTN SiLK imported flow persistence failed for service %s: %v", serviceID, err)
			}
		}
	}
	return len(accepted)
}
