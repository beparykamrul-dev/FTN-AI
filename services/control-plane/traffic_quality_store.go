package main

import (
    "errors"
    "strings"
    "sync"
    "time"
)

const (
    trafficQualityTTL     = 2 * time.Minute
    trafficQualityMaxKeys = 8192
)

type TrafficQualityObservation struct {
    PathID      string       `json:"path_id"`
    ServiceID   string       `json:"service_id"`
    Class       TrafficClass `json:"class"`
    LatencyMs   float64      `json:"latency_ms"`
    JitterMs    float64      `json:"jitter_ms"`
    PacketLoss  float64      `json:"packet_loss_percent"`
    Congestion  float64      `json:"congestion"`
    Healthy     bool         `json:"healthy"`
    ObservedAt  time.Time    `json:"observed_at"`
}

type TrafficQualityStore struct {
    mu    sync.RWMutex
    items map[string]TrafficQualityObservation
}

func NewTrafficQualityStore() *TrafficQualityStore {
    return &TrafficQualityStore{items: make(map[string]TrafficQualityObservation)}
}

func trafficQualityKey(serviceID, pathID string) string {
    return strings.TrimSpace(serviceID) + "\x00" + strings.TrimSpace(pathID)
}

func validateTrafficQualityObservation(o TrafficQualityObservation, now time.Time) error {
    o.PathID = strings.TrimSpace(o.PathID)
    o.ServiceID = strings.TrimSpace(o.ServiceID)
    if o.PathID == "" || o.ServiceID == "" {
        return errors.New("invalid_traffic_quality_identity")
    }
    if _, ok := trafficPolicyByID(o.ServiceID); !ok {
        return errors.New("unknown_service_id")
    }
    if o.LatencyMs < 0 || o.JitterMs < 0 || o.PacketLoss < 0 || o.PacketLoss > 100 || o.Congestion < 0 || o.Congestion > 1 {
        return errors.New("invalid_traffic_quality_metric")
    }
    if o.ObservedAt.IsZero() || now.Sub(o.ObservedAt) > trafficQualityTTL || o.ObservedAt.After(now.Add(10*time.Second)) {
        return errors.New("stale_traffic_quality")
    }
    return nil
}

func (s *TrafficQualityStore) Upsert(o TrafficQualityObservation, now time.Time) error {
    if s == nil {
        return errors.New("traffic_quality_store_required")
    }
    if err := validateTrafficQualityObservation(o, now); err != nil {
        return err
    }
    key := trafficQualityKey(o.ServiceID, o.PathID)
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, exists := s.items[key]; !exists && len(s.items) >= trafficQualityMaxKeys {
        return errors.New("traffic_quality_registry_limit")
    }
    s.items[key] = o
    return nil
}

func (s *TrafficQualityStore) Snapshot(serviceID string, now time.Time) []TrafficQualityObservation {
    if s == nil {
        return nil
    }
    serviceID = strings.TrimSpace(serviceID)
    s.mu.RLock()
    out := make([]TrafficQualityObservation, 0, len(s.items))
    for _, o := range s.items {
        if serviceID != "" && o.ServiceID != serviceID {
            continue
        }
        if o.ObservedAt.IsZero() || now.Sub(o.ObservedAt) > trafficQualityTTL {
            continue
        }
        out = append(out, o)
    }
    s.mu.RUnlock()
    return out
}

func (s *TrafficQualityStore) Prune(now time.Time) int {
    if s == nil {
        return 0
    }
    s.mu.Lock()
    defer s.mu.Unlock()
    removed := 0
    for key, o := range s.items {
        if o.ObservedAt.IsZero() || now.Sub(o.ObservedAt) > trafficQualityTTL {
            delete(s.items, key)
            removed++
        }
    }
    return removed
}
