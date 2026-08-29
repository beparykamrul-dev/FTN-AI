package main

import (
	"net/http"
	"sort"
	"strings"
)

type Node struct {
	ID         string   `json:"id"`
	Provider   string   `json:"provider"`
	Region     string   `json:"region,omitempty"`
	Services   []string `json:"services"`
	CPUPercent float64  `json:"cpu_percent"`
	RAMPercent float64  `json:"ram_percent"`
	SSDPercent float64  `json:"ssd_percent"`
	HDDPercent float64  `json:"hdd_percent"`
	NetMbps    float64  `json:"net_mbps"`
	LatencyMs  float64  `json:"latency_ms"`
	PacketLoss float64  `json:"packet_loss_percent"`
	Healthy    bool     `json:"healthy"`
}

type PlacementRequest struct {
	ServiceID string `json:"service_id"`
	Region    string `json:"region,omitempty"`
	Provider  string `json:"provider,omitempty"`
}

var nodes = []Node{}

func nodeHasService(n Node, serviceID string) bool {
	for _, s := range n.Services {
		if s == serviceID || s == "*" {
			return true
		}
	}
	return false
}

func nodeScore(n Node, req PlacementRequest) float64 {
	if !n.Healthy || !nodeHasService(n, req.ServiceID) {
		return -1
	}
	score := 100.0
	score -= n.CPUPercent * 0.20
	score -= n.RAMPercent * 0.20
	score -= n.SSDPercent * 0.10
	score -= n.HDDPercent * 0.05
	score -= n.LatencyMs * 0.50
	score -= n.PacketLoss * 2
	if req.Region != "" && strings.EqualFold(req.Region, n.Region) {
		score += 20
	}
	if req.Provider != "" && strings.EqualFold(req.Provider, n.Provider) {
		score += 10
	}
	return score
}

func (a *App) nodeCatalog(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"nodes": nodes})
}

func (a *App) placement(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var req PlacementRequest
	if err := decodeJSON(r, &req); err != nil || req.ServiceID == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "service_id_required"})
		return
	}
	candidates := make([]Node, 0, len(nodes))
	for _, n := range nodes {
		if nodeScore(n, req) >= 0 {
			candidates = append(candidates, n)
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool { return nodeScore(candidates[i], req) > nodeScore(candidates[j], req) })
	if len(candidates) == 0 {
		jsonResponse(w, http.StatusServiceUnavailable, map[string]any{"status": "no_eligible_node", "service_id": req.ServiceID})
		return
	}
	best := candidates[0]
	jsonResponse(w, http.StatusOK, map[string]any{
		"status":     "placement_ready",
		"service_id": req.ServiceID,
		"node_id":    best.ID,
		"provider":   best.Provider,
		"region":     best.Region,
		"score":      nodeScore(best, req),
		"candidates": candidates,
		"execution":  "approval_required",
	})
}
