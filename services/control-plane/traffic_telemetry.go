package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type TrafficQualityObservation struct {
	ServiceID  string       `json:"service_id"`
	CustomerID string       `json:"customer_id,omitempty"`
	PathID     string       `json:"path_id"`
	LatencyMs  float64      `json:"latency_ms"`
	JitterMs   float64      `json:"jitter_ms"`
	PacketLoss float64      `json:"packet_loss_percent"`
	Congestion float64      `json:"congestion"`
	Healthy    bool         `json:"healthy"`
	ObservedAt time.Time    `json:"observed_at"`
}

type TrafficQualityStore struct {
	mu      sync.RWMutex
	items   map[string]TrafficQualityObservation
	maxAge  time.Duration
	maxKeys int
}

func NewTrafficQualityStore() *TrafficQualityStore {
	return &TrafficQualityStore{items: make(map[string]TrafficQualityObservation), maxAge: 2 * time.Minute, maxKeys: 8192}
}

func (s *TrafficQualityStore) Upsert(o TrafficQualityObservation, now time.Time) error {
	o.ServiceID = strings.TrimSpace(o.ServiceID)
	o.CustomerID = strings.TrimSpace(o.CustomerID)
	o.PathID = strings.TrimSpace(o.PathID)
	if o.ServiceID == "" || o.PathID == "" {
		return errors.New("service_id_and_path_id_required")
	}
	if _, ok := trafficPolicyByID(o.ServiceID); !ok {
		return errors.New("unknown_service_id")
	}
	if o.LatencyMs < 0 || o.JitterMs < 0 || o.PacketLoss < 0 || o.PacketLoss > 100 || o.Congestion < 0 || o.Congestion > 1 {
		return errors.New("invalid_quality_metrics")
	}
	if o.ObservedAt.IsZero() {
		o.ObservedAt = now
	}
	if now.Sub(o.ObservedAt) > s.maxAge || o.ObservedAt.After(now.Add(10*time.Second)) {
		return errors.New("stale_observation")
	}
	key := o.ServiceID + "\x00" + o.CustomerID + "\x00" + o.PathID
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[key]; !exists && len(s.items) >= s.maxKeys {
		return errors.New("quality_store_limit")
	}
	s.items[key] = o
	return nil
}

func (s *TrafficQualityStore) Snapshot(now time.Time) []TrafficQualityObservation {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]TrafficQualityObservation, 0, len(s.items))
	for key, o := range s.items {
		if now.Sub(o.ObservedAt) > s.maxAge {
			delete(s.items, key)
			continue
		}
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ServiceID != out[j].ServiceID { return out[i].ServiceID < out[j].ServiceID }
		if out[i].CustomerID != out[j].CustomerID { return out[i].CustomerID < out[j].CustomerID }
		return out[i].PathID < out[j].PathID
	})
	return out
}

func (a *App) trafficQuality(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) { return }
	if !requirePermission(a, "network.read", w, r) { return }
	var o TrafficQualityObservation
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
	if err := dec.Decode(&o); err != nil { jsonResponse(w, http.StatusBadRequest, map[string]string{"error":"invalid_json"}); return }
	if err := a.trafficQuality.Upsert(o, time.Now().UTC()); err != nil { jsonResponse(w, http.StatusBadRequest, map[string]string{"error":err.Error()}); return }
	jsonResponse(w, http.StatusAccepted, map[string]string{"status":"accepted"})
}

func (a *App) trafficQualitySnapshot(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) { return }
	if !requirePermission(a, "network.read", w, r) { return }
	jsonResponse(w, http.StatusOK, map[string]any{"observations": a.trafficQuality.Snapshot(time.Now().UTC())})
}
