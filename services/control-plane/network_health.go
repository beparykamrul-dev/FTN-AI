package main

import (
	"net/http"
	"sort"
	"time"
)

type PathHealth struct {
	PathID       string  `json:"path_id"`
	Provider     string  `json:"provider"`
	Region       string  `json:"region,omitempty"`
	LatencyMs    float64 `json:"latency_ms"`
	PacketLoss   float64 `json:"packet_loss_percent"`
	JitterMs     float64 `json:"jitter_ms,omitempty"`
	Healthy      bool    `json:"healthy"`
	LastSeen     time.Time `json:"last_seen,omitempty"`
}

func pathHealthFromNode(n Node) PathHealth {
	return PathHealth{
		PathID: n.ID, Provider: n.Provider, Region: n.Region,
		LatencyMs: n.LatencyMs, PacketLoss: n.PacketLoss,
		Healthy: n.Healthy, LastSeen: n.LastSeen,
	}
}

func (a *App) networkHealth(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) { return }
	ns, err := a.loadNodes(r.Context())
	if err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"node_query_failed"}); return }
	paths := make([]PathHealth, 0, len(ns))
	for _, n := range ns { paths = append(paths, pathHealthFromNode(n)) }
	sort.SliceStable(paths, func(i, j int) bool {
		if paths[i].Healthy != paths[j].Healthy { return paths[i].Healthy }
		if paths[i].LatencyMs != paths[j].LatencyMs { return paths[i].LatencyMs < paths[j].LatencyMs }
		return paths[i].PacketLoss < paths[j].PacketLoss
	})
	jsonResponse(w, http.StatusOK, map[string]any{"paths": paths, "policy":"observe_only"})
}
