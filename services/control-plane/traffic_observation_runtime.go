package main

import (
	"context"
	"log"
	"time"
)

// IngestObservations is the timestamp-preserving ingestion path for imported
// flow sources such as nfcapd. Subscriber/main-server identity comes only from
// the managed endpoint registry; imported payloads cannot override it.
func (t *TrafficRuntime) IngestObservations(ctx context.Context, observations []FlowObservation) int {
	if t == nil || len(observations) == 0 { return 0 }
	if ctx == nil { ctx = context.Background() }
	accepted := make([]TrafficFlowObservation, 0, len(observations))
	for _, o := range observations {
		if err := o.Validate(); err != nil { continue }
		obs, ok := t.Classify(o.FlowRecord, o.ObservedAt)
		if !ok { continue }
		obs.ObservedAt = o.ObservedAt.UTC()
		obs.FirstSeen = o.FirstSeen.UTC()
		obs.LastSeen = o.LastSeen.UTC()
		if obs.FirstSeen.IsZero() { obs.FirstSeen = obs.ObservedAt }
		if obs.LastSeen.IsZero() { obs.LastSeen = obs.ObservedAt }
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
		type groupKey struct { serviceID, subscriberID, mainServerID string }
		groups := make(map[groupKey][]FlowObservation)
		for _, f := range accepted {
			o, err := normalizeFlowObservation(f.FlowRecord, f.ObservedAt, f.FirstSeen, f.LastSeen)
			if err != nil { continue }
			key := groupKey{serviceID: f.ServiceID, subscriberID: f.SubscriberID, mainServerID: f.MainServerID}
			groups[key] = append(groups[key], o)
		}
		for key, batch := range groups {
			if err := silk.InsertBatchObservations(ctx, batch, key.serviceID, key.subscriberID, key.mainServerID); err != nil {
				log.Printf("FTN SiLK imported flow persistence failed for service %s: %v", key.serviceID, err)
			}
		}
	}
	return len(accepted)
}
