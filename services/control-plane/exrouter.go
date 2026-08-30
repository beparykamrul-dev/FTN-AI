package main

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type ExRouterRequest struct {
	ServiceID string `json:"service_id"`
	Region string `json:"region,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type ExRoute struct {
	PathID string `json:"path_id"`
	Provider string `json:"provider"`
	Region string `json:"region,omitempty"`
	Score float64 `json:"score"`
	LatencyMs float64 `json:"latency_ms"`
	PacketLoss float64 `json:"packet_loss_percent"`
	JitterMs float64 `json:"jitter_ms,omitempty"`
	ThroughputMbps float64 `json:"throughput_mbps,omitempty"`
	Retransmissions float64 `json:"retransmissions,omitempty"`
	Healthy bool `json:"healthy"`
	LastSeen time.Time `json:"last_seen,omitempty"`
}

var exRouteState = struct { sync.Mutex; selected map[string]string }{selected: map[string]string{}}

func exRouteScore(n Node, req ExRouterRequest) float64 {
	if !n.Healthy || !nodeHasService(n, req.ServiceID) { return -1 }
	// Lower latency/loss/jitter/retransmissions are better; higher throughput is better.
	score := 100.0 - n.LatencyMs - n.PacketLoss*4 - n.JitterMs*0.5 - n.Retransmissions*2
	score += math.Min(n.NetMbps/100.0, 20.0)
	if req.Region != "" && strings.EqualFold(req.Region, n.Region) { score += 15 }
	if req.Provider != "" && strings.EqualFold(req.Provider, n.Provider) { score += 8 }
	return score
}

func (a *App) exRouter(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) { return }
	var req ExRouterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ServiceID == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error":"service_id_required"}); return
	}
	if !serviceExists(a.catalog, req.ServiceID) {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error":"service_not_found"}); return
	}
	ns, err := a.loadNodes(r.Context())
	if err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"node_query_failed"}); return }
	routes := make([]ExRoute, 0, len(ns))
	for _, n := range ns {
		score := exRouteScore(n, req)
		if score < 0 { continue }
		routes = append(routes, ExRoute{PathID:n.ID, Provider:n.Provider, Region:n.Region, Score:score, LatencyMs:n.LatencyMs, PacketLoss:n.PacketLoss, JitterMs:n.JitterMs, ThroughputMbps:n.NetMbps, Retransmissions:n.Retransmissions, Healthy:n.Healthy, LastSeen:n.LastSeen})
	}
	sort.SliceStable(routes, func(i, j int) bool { return routes[i].Score > routes[j].Score })
	if len(routes) == 0 { jsonResponse(w, http.StatusServiceUnavailable, map[string]any{"status":"no_healthy_path","service_id":req.ServiceID}); return }
	// Hysteresis: retain the current path unless the new path is materially better.
	exRouteState.Lock()
	current := exRouteState.selected[req.ServiceID]
	selected := routes[0]
	if current != "" {
		for _, candidate := range routes { if candidate.PathID == current && candidate.Score+10 >= selected.Score { selected = candidate; break } }
	}
	exRouteState.selected[req.ServiceID] = selected.PathID
	exRouteState.Unlock()
	jsonResponse(w, http.StatusOK, map[string]any{"status":"path_ready", "service_id":req.ServiceID, "selected_path":selected, "paths":routes, "service_unchanged":true, "execution":"policy_controlled", "hysteresis":true})
}
