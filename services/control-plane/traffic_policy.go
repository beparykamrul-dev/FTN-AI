package main

import (
	"errors"
	"sort"
	"strings"
	"time"
)

type TrafficClass string

const (
	TrafficRealtime TrafficClass = "realtime"
	TrafficVoice    TrafficClass = "voice"
	TrafficVideo    TrafficClass = "video"
	TrafficGaming   TrafficClass = "gaming"
	TrafficNormal   TrafficClass = "normal"
	TrafficBulk     TrafficClass = "bulk"
)

type TrafficServicePolicy struct {
	ID       string
	Class    TrafficClass
	Priority uint8
	DSCP     uint8
}

func DefaultTrafficServicePolicies() []TrafficServicePolicy {
	return []TrafficServicePolicy{
		{ID: "pubg", Class: TrafficGaming, Priority: 95, DSCP: 46},
		{ID: "freefire", Class: TrafficGaming, Priority: 95, DSCP: 46},
		{ID: "imo", Class: TrafficVoice, Priority: 92, DSCP: 46},
		{ID: "whatsapp", Class: TrafficRealtime, Priority: 90, DSCP: 46},
		{ID: "telegram", Class: TrafficRealtime, Priority: 90, DSCP: 46},
		{ID: "realtime-generic", Class: TrafficRealtime, Priority: 85, DSCP: 46},
		{ID: "voice-generic", Class: TrafficVoice, Priority: 88, DSCP: 46},
		{ID: "video-generic", Class: TrafficVideo, Priority: 75, DSCP: 34},
		{ID: "normal", Class: TrafficNormal, Priority: 40, DSCP: 0},
		{ID: "bulk", Class: TrafficBulk, Priority: 10, DSCP: 8},
	}
}

type TrafficPathObservation struct {
	PathID     string
	ServiceID  string
	Class      TrafficClass
	LatencyMs  float64
	JitterMs   float64
	PacketLoss float64
	Congestion float64
	Healthy    bool
	ObservedAt time.Time
}

type TrafficDecision struct {
	PathID      string `json:"path_id"`
	ServiceID   string `json:"service_id"`
	Class       TrafficClass `json:"class"`
	Score       float64 `json:"score"`
	DSCP        uint8 `json:"dscp"`
	Priority    uint8 `json:"priority"`
	Failover    bool `json:"failover"`
	HoldDownSec int `json:"hold_down_sec"`
}

func trafficScore(o TrafficPathObservation, now time.Time) float64 {
	if !o.Healthy || strings.TrimSpace(o.PathID) == "" { return -1 }
	s := 100 - o.LatencyMs - 1.5*o.JitterMs - 5*o.PacketLoss - 20*o.Congestion
	if o.ObservedAt.IsZero() { return s - 10 }
	if age := now.Sub(o.ObservedAt); age > 30*time.Second { s -= 20 }
	if age > 2*time.Minute { return -1 }
	return s
}

// SelectTrafficPath performs stateless selection. Stateful hold-down is handled
// by TrafficPathController so the pure selector remains deterministic and testable.
func SelectTrafficPath(observations []TrafficPathObservation, service TrafficServicePolicy, now time.Time) (TrafficDecision, bool) {
	best := TrafficDecision{}
	bestScore := -1.0
	for _, o := range observations {
		if strings.TrimSpace(o.ServiceID) != service.ID && strings.TrimSpace(string(o.Class)) != string(service.Class) { continue }
		s := trafficScore(o, now)
		if s > bestScore || (s == bestScore && o.PathID < best.PathID) {
			bestScore = s
			best = TrafficDecision{PathID: o.PathID, ServiceID: service.ID, Class: service.Class, Score: s, DSCP: service.DSCP, Priority: service.Priority, HoldDownSec: 5}
		}
	}
	if bestScore < 0 || best.PathID == "" { return TrafficDecision{}, false }
	return best, true
}

type TrafficPathController struct {
	CurrentPath string
	CandidatePath string
	CandidateSince time.Time
	LastSwitch time.Time
}

// Decide prevents route flapping. A healthy current path is retained unless a
// materially better candidate has stayed preferred for the hold-down period.
// An unhealthy current path may fail over immediately.
func (c *TrafficPathController) Decide(observations []TrafficPathObservation, service TrafficServicePolicy, now time.Time) (TrafficDecision, bool) {
	best, ok := SelectTrafficPath(observations, service, now)
	if !ok { return TrafficDecision{}, false }
	if c.CurrentPath == "" {
		c.CurrentPath, c.LastSwitch = best.PathID, now
		best.Failover = false
		return best, true
	}
	currentHealthy := false
	currentScore := -1.0
	for _, o := range observations {
		if o.PathID == c.CurrentPath {
			currentScore = trafficScore(o, now)
			currentHealthy = currentScore >= 0
			break
		}
	}
	if best.PathID == c.CurrentPath {
		c.CandidatePath = ""
		return best, true
	}
	if !currentHealthy {
		c.CurrentPath, c.CandidatePath, c.LastSwitch = best.PathID, "", now
		best.Failover = true
		return best, true
	}
	if best.Score <= currentScore+5 {
		return TrafficDecision{PathID: c.CurrentPath, ServiceID: service.ID, Class: service.Class, Score: currentScore, DSCP: service.DSCP, Priority: service.Priority, HoldDownSec: 5}, true
	}
	if c.CandidatePath != best.PathID {
		c.CandidatePath, c.CandidateSince = best.PathID, now
		return TrafficDecision{PathID: c.CurrentPath, ServiceID: service.ID, Class: service.Class, Score: currentScore, DSCP: service.DSCP, Priority: service.Priority, HoldDownSec: 5}, true
	}
	if now.Sub(c.CandidateSince) < 5*time.Second {
		return TrafficDecision{PathID: c.CurrentPath, ServiceID: service.ID, Class: service.Class, Score: currentScore, DSCP: service.DSCP, Priority: service.Priority, HoldDownSec: 5}, true
	}
	c.CurrentPath, c.CandidatePath, c.LastSwitch = best.PathID, "", now
	best.Failover = false
	return best, true
}

func SortedTrafficPolicies(in []TrafficServicePolicy) []TrafficServicePolicy {
	out := append([]TrafficServicePolicy(nil), in...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority { return out[i].ID < out[j].ID }
		return out[i].Priority > out[j].Priority
	})
	return out
}

var errTrafficPolicyInvalid = errors.New("traffic policy is invalid")

func validateTrafficPolicy(p TrafficServicePolicy) error {
	if strings.TrimSpace(p.ID) == "" { return errTrafficPolicyInvalid }
	if p.DSCP > 63 { return errTrafficPolicyInvalid }
	return nil
}
