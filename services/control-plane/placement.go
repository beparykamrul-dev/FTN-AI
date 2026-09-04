package main

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type Node struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"-"`
	Provider   string    `json:"provider"`
	Region     string    `json:"region,omitempty"`
	Services   []string  `json:"services"`
	CPUPercent float64   `json:"cpu_percent"`
	RAMPercent float64   `json:"ram_percent"`
	SSDPercent float64   `json:"ssd_percent"`
	HDDPercent float64   `json:"hdd_percent"`
	NetMbps    float64   `json:"net_mbps"`
	LatencyMs  float64   `json:"latency_ms"`
	PacketLoss float64   `json:"packet_loss_percent"`
	Healthy    bool      `json:"healthy"`
	LastSeen   time.Time `json:"last_seen,omitempty"`
}

type PlacementRequest struct {
	ServiceID string `json:"service_id"`
	Region    string `json:"region,omitempty"`
	Provider  string `json:"provider,omitempty"`
}

var (
	nodes   = []Node{}
	nodesMu sync.RWMutex
)

const nodeHeartbeatTTL = 90 * time.Second

func nodeHasService(n Node, serviceID string) bool {
	for _, s := range n.Services {
		if s == serviceID || s == "*" {
			return true
		}
	}
	return false
}

func validNode(n Node) bool {
	finite := func(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
	return strings.TrimSpace(n.ID) != "" && strings.TrimSpace(n.Provider) != "" &&
		finite(n.CPUPercent) && n.CPUPercent >= 0 && n.CPUPercent <= 100 &&
		finite(n.RAMPercent) && n.RAMPercent >= 0 && n.RAMPercent <= 100 &&
		finite(n.SSDPercent) && n.SSDPercent >= 0 && n.SSDPercent <= 100 &&
		finite(n.HDDPercent) && n.HDDPercent >= 0 && n.HDDPercent <= 100 &&
		finite(n.NetMbps) && n.NetMbps >= 0 && finite(n.LatencyMs) && n.LatencyMs >= 0 &&
		finite(n.PacketLoss) && n.PacketLoss >= 0 && n.PacketLoss <= 100
}

func serviceExists(catalog []Service, serviceID string) bool {
	for _, s := range catalog {
		if s.ID == serviceID {
			return true
		}
	}
	return false
}

func nodeScore(n Node, req PlacementRequest) float64 {
	if !n.Healthy || !nodeHasService(n, req.ServiceID) {
		return -1
	}
	score := 100.0 - n.CPUPercent*.20 - n.RAMPercent*.20 - n.SSDPercent*.10 - n.HDDPercent*.05 - n.LatencyMs*.50 - n.PacketLoss*2
	if req.Region != "" && strings.EqualFold(req.Region, n.Region) {
		score += 20
	}
	if req.Provider != "" && strings.EqualFold(req.Provider, n.Provider) {
		score += 10
	}
	return score
}

func refreshNodeHealth(list []Node, now time.Time) []Node {
	out := make([]Node, len(list))
	copy(out, list)
	for i := range out {
		if out[i].LastSeen.IsZero() || now.Sub(out[i].LastSeen) > nodeHeartbeatTTL {
			out[i].Healthy = false
		}
	}
	return out
}

func (a *App) loadNodes(ctx context.Context, tenantID string) ([]Node, error) {
	if tenantID == "" {
		return nil, context.Canceled
	}
	if a.db == nil {
		nodesMu.RLock()
		defer nodesMu.RUnlock()
		out := make([]Node, 0, len(nodes))
		for _, n := range nodes {
			if n.TenantID == tenantID {
				out = append(out, n)
			}
		}
		return refreshNodeHealth(out, time.Now().UTC()), nil
	}
	qctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	rows, err := a.db.Query(qctx, `select id,provider,region,services,cpu_percent,ram_percent,ssd_percent,hdd_percent,net_mbps,latency_ms,packet_loss_percent,healthy,updated_at from control_nodes where tenant_id=$1::uuid order by id`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Node{}
	for rows.Next() {
		var n Node
		n.TenantID = tenantID
		if err := rows.Scan(&n.ID, &n.Provider, &n.Region, &n.Services, &n.CPUPercent, &n.RAMPercent, &n.SSDPercent, &n.HDDPercent, &n.NetMbps, &n.LatencyMs, &n.PacketLoss, &n.Healthy, &n.LastSeen); err != nil {
			return nil, err
		}
		if validNode(n) {
			out = append(out, n)
		}
	}
	return refreshNodeHealth(out, time.Now().UTC()), rows.Err()
}

func (a *App) nodeCatalog(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !requirePermission(a, "node.read", w, r) {
		return
	}
	rc := requestInfo(r)
	out, err := a.loadNodes(r.Context(), rc.TenantID)
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "node_query_failed"})
		return
	}
	jsonResponse(w, 200, map[string]any{"nodes": out})
}

func upsertMemoryNode(n Node) {
	nodesMu.Lock()
	defer nodesMu.Unlock()
	for i := range nodes {
		if nodes[i].ID == n.ID && nodes[i].TenantID == n.TenantID {
			nodes[i] = n
			return
		}
	}
	nodes = append(nodes, n)
}

func (a *App) registerNode(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !requirePermission(a, "node.manage", w, r) {
		return
	}
	rc := requestInfo(r)
	if rc.TenantID == "" {
		jsonResponse(w, 400, map[string]string{"error": "tenant_required"})
		return
	}
	var n Node
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10)).Decode(&n); err != nil {
		jsonResponse(w, 400, map[string]string{"error": "invalid_json"})
		return
	}
	if !validNode(n) {
		jsonResponse(w, 400, map[string]string{"error": "invalid_node"})
		return
	}
	n.TenantID = rc.TenantID
	n.LastSeen = time.Now().UTC()
	if a.db == nil {
		upsertMemoryNode(n)
		jsonResponse(w, 202, map[string]any{"status": "registered", "node": n})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	res, err := a.db.Exec(ctx, `insert into control_nodes(tenant_id,id,provider,region,services,cpu_percent,ram_percent,ssd_percent,hdd_percent,net_mbps,latency_ms,packet_loss_percent,healthy,updated_at) values($1::uuid,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,now()) on conflict(id) do update set tenant_id=excluded.tenant_id,provider=excluded.provider,region=excluded.region,services=excluded.services,cpu_percent=excluded.cpu_percent,ram_percent=excluded.ram_percent,ssd_percent=excluded.ssd_percent,hdd_percent=excluded.hdd_percent,net_mbps=excluded.net_mbps,latency_ms=excluded.latency_ms,packet_loss_percent=excluded.packet_loss_percent,healthy=excluded.healthy,updated_at=now() where control_nodes.tenant_id=excluded.tenant_id`, n.TenantID, n.ID, n.Provider, n.Region, n.Services, n.CPUPercent, n.RAMPercent, n.SSDPercent, n.HDDPercent, n.NetMbps, n.LatencyMs, n.PacketLoss, n.Healthy)
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "node_register_failed"})
		return
	}
	if res.RowsAffected() != 1 {
		jsonResponse(w, 409, map[string]string{"error": "node_id_conflict"})
		return
	}
	jsonResponse(w, 202, map[string]any{"status": "registered", "node": n})
}

func (a *App) placement(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !requirePermission(a, "node.read", w, r) {
		return
	}
	var req PlacementRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&req); err != nil {
		jsonResponse(w, 400, map[string]string{"error": "invalid_json"})
		return
	}
	req.ServiceID, req.Region, req.Provider = strings.TrimSpace(req.ServiceID), strings.TrimSpace(req.Region), strings.TrimSpace(req.Provider)
	if req.ServiceID == "" {
		jsonResponse(w, 400, map[string]string{"error": "service_id_required"})
		return
	}
	if !serviceExists(a.catalog, req.ServiceID) {
		jsonResponse(w, 404, map[string]string{"error": "service_not_found"})
		return
	}
	rc := requestInfo(r)
	available, err := a.loadNodes(r.Context(), rc.TenantID)
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "node_query_failed"})
		return
	}
	candidates := make([]Node, 0, len(available))
	for _, n := range available {
		if nodeScore(n, req) >= 0 {
			candidates = append(candidates, n)
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		si, sj := nodeScore(candidates[i], req), nodeScore(candidates[j], req)
		if si == sj {
			return candidates[i].ID < candidates[j].ID
		}
		return si > sj
	})
	if len(candidates) == 0 {
		jsonResponse(w, 503, map[string]any{"status": "no_eligible_node", "service_id": req.ServiceID})
		return
	}
	best := candidates[0]
	jsonResponse(w, 200, map[string]any{"status": "placement_ready", "service_id": req.ServiceID, "node_id": best.ID, "provider": best.Provider, "region": best.Region, "score": nodeScore(best, req), "candidates": candidates, "execution": "approval_required"})
}
