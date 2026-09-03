package main

import (
    "encoding/json"
    "errors"
    "net"
    "net/http"
    "sort"
    "strings"
    "sync"
    "time"
)

type ManagedEndpoint struct {
    ServiceID string    `json:"service_id"`
    CIDR      string    `json:"cidr"`
    Region    string    `json:"region,omitempty"`
    Provider  string    `json:"provider,omitempty"`
    ExpiresAt time.Time `json:"expires_at,omitempty"`
    network   *net.IPNet
}

type TrafficFlowObservation struct {
    FlowRecord
    ServiceID string       `json:"service_id"`
    Class     TrafficClass `json:"class"`
    PathID    string       `json:"path_id"`
    Region    string       `json:"region,omitempty"`
    Provider  string       `json:"provider,omitempty"`
    ObservedAt time.Time   `json:"observed_at"`
}

type TrafficRuntime struct {
    mu        sync.RWMutex
    endpoints []ManagedEndpoint
    flows     []TrafficFlowObservation
    decisions map[string]TrafficDecision
    controllers map[string]*TrafficPathController
}

func NewTrafficRuntime() *TrafficRuntime {
    return &TrafficRuntime{
        decisions: make(map[string]TrafficDecision),
        controllers: make(map[string]*TrafficPathController),
    }
}

func (t *TrafficRuntime) UpsertEndpoint(e ManagedEndpoint) error {
    if strings.TrimSpace(e.ServiceID) == "" || strings.TrimSpace(e.CIDR) == "" {
        return errors.New("service_id_and_cidr_required")
    }
    _, n, err := net.ParseCIDR(e.CIDR)
    if err != nil { return errors.New("invalid_cidr") }
    e.network = n
    t.mu.Lock()
    defer t.mu.Unlock()
    for i := range t.endpoints {
        if t.endpoints[i].ServiceID == e.ServiceID && t.endpoints[i].CIDR == e.CIDR {
            t.endpoints[i] = e
            return nil
        }
    }
    t.endpoints = append(t.endpoints, e)
    return nil
}

func (t *TrafficRuntime) Classify(f FlowRecord, now time.Time) (TrafficFlowObservation, bool) {
    dst := net.ParseIP(strings.TrimSpace(f.DestinationIP))
    src := net.ParseIP(strings.TrimSpace(f.SourceIP))
    if dst == nil && src == nil { return TrafficFlowObservation{}, false }
    t.mu.RLock()
    endpoints := append([]ManagedEndpoint(nil), t.endpoints...)
    t.mu.RUnlock()
    var best *ManagedEndpoint
    bestBits := -1
    for i := range endpoints {
        e := &endpoints[i]
        if !e.ExpiresAt.IsZero() && !now.Before(e.ExpiresAt) { continue }
        matched := false
        if dst != nil && e.network.Contains(dst) { matched = true }
        if src != nil && e.network.Contains(src) { matched = true }
        if !matched { continue }
        ones, _ := e.network.Mask.Size()
        if ones > bestBits { best, bestBits = e, ones }
    }
    if best == nil { return TrafficFlowObservation{}, false }
    policy, ok := trafficPolicyByID(best.ServiceID)
    if !ok { return TrafficFlowObservation{}, false }
    return TrafficFlowObservation{FlowRecord:f, ServiceID:best.ServiceID, Class:policy.Class, Region:best.Region, Provider:best.Provider, ObservedAt:now}, true
}

func trafficPolicyByID(id string) (TrafficServicePolicy, bool) {
    id = strings.TrimSpace(id)
    for _, p := range DefaultTrafficServicePolicies() {
        if p.ID == id { return p, true }
    }
    return TrafficServicePolicy{}, false
}

func (t *TrafficRuntime) Ingest(flows []FlowRecord, now time.Time) int {
    added := 0
    t.mu.Lock()
    defer t.mu.Unlock()
    for _, f := range flows {
        t.mu.Unlock()
        obs, ok := t.Classify(f, now)
        t.mu.Lock()
        if !ok { continue }
        t.flows = append(t.flows, obs)
        added++
    }
    cutoff := now.Add(-2 * time.Minute)
    keep := t.flows[:0]
    for _, f := range t.flows {
        if f.ObservedAt.After(cutoff) { keep = append(keep, f) }
    }
    t.flows = keep
    if len(t.flows) > 4096 { t.flows = t.flows[len(t.flows)-4096:] }
    return added
}

func (t *TrafficRuntime) Decisions(now time.Time, nodes []Node) []TrafficDecision {
    t.mu.Lock()
    defer t.mu.Unlock()
    byService := make(map[string][]TrafficPathObservation)
    for _, f := range t.flows {
        if f.PathID == "" { continue }
        byService[f.ServiceID] = append(byService[f.ServiceID], TrafficPathObservation{PathID:f.PathID, ServiceID:f.ServiceID, Class:f.Class, LatencyMs:0, JitterMs:0, PacketLoss:0, Congestion:0, Healthy:true, ObservedAt:f.ObservedAt})
    }
    // Current node telemetry is the authoritative path-health source until
    // dedicated active probes populate jitter/loss/congestion observations.
    for _, n := range nodes {
        for serviceID, observations := range byService {
            _ = observations
            p, ok := trafficPolicyByID(serviceID); if !ok { continue }
            byService[serviceID] = append(byService[serviceID], TrafficPathObservation{PathID:n.ID, ServiceID:serviceID, Class:p.Class, LatencyMs:n.LatencyMs, PacketLoss:n.PacketLoss, Healthy:n.Healthy, ObservedAt:now})
        }
    }
    out := make([]TrafficDecision, 0, len(byService))
    for serviceID, observations := range byService {
        p, ok := trafficPolicyByID(serviceID); if !ok { continue }
        c := t.controllers[serviceID]
        if c == nil { c = &TrafficPathController{}; t.controllers[serviceID] = c }
        d, ok := c.Decide(observations, p, now); if !ok { continue }
        t.decisions[serviceID] = d
        out = append(out, d)
    }
    sort.Slice(out, func(i,j int) bool { return out[i].ServiceID < out[j].ServiceID })
    return out
}

func (a *App) trafficEndpoints(w http.ResponseWriter, r *http.Request) {
    if !method(w,r,http.MethodPost) { return }
    if !requirePermission(a,"network.read",w,r) { return }
    var e ManagedEndpoint
    if err := json.NewDecoder(r.Body).Decode(&e); err != nil { jsonResponse(w,400,map[string]string{"error":"invalid_json"}); return }
    if err := a.traffic.UpsertEndpoint(e); err != nil { jsonResponse(w,400,map[string]string{"error":err.Error()}); return }
    a.audit(r,"traffic.endpoint",e.ServiceID,"accepted",e)
    jsonResponse(w,http.StatusAccepted,map[string]any{"status":"accepted","service_id":e.ServiceID,"cidr":e.CIDR,"source":"managed-endpoint-feed"})
}

func (a *App) trafficFlowIngest(w http.ResponseWriter, r *http.Request) {
    if !method(w,r,http.MethodPost) { return }
    if !requirePermission(a,"network.read",w,r) { return }
    var req struct { Flows []FlowRecord `json:"flows"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { jsonResponse(w,400,map[string]string{"error":"invalid_json"}); return }
    now := time.Now().UTC()
    accepted := a.traffic.Ingest(req.Flows, now)
    jsonResponse(w,http.StatusAccepted,map[string]any{"accepted":accepted,"received":len(req.Flows),"observed_at":now})
}

func (a *App) trafficDecisions(w http.ResponseWriter, r *http.Request) {
    if !method(w,r,http.MethodGet) { return }
    if !requirePermission(a,"network.read",w,r) { return }
    nodes, err := a.loadNodes(r.Context())
    if err != nil { jsonResponse(w,500,map[string]string{"error":"node_query_failed"}); return }
    decisions := a.traffic.Decisions(time.Now().UTC(), nodes)
    jsonResponse(w,200,map[string]any{"decisions":decisions,"execution":"read-only","configuration_changes":"approval-required"})
}
