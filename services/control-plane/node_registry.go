package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Node struct {
	ID string `json:"id"`
	Provider string `json:"provider"`
	Region string `json:"region"`
	Services []string `json:"services"`
	CPUPercent float64 `json:"cpu_percent"`
	RAMPercent float64 `json:"ram_percent"`
	SSDPercent float64 `json:"ssd_percent"`
	HDDPercent float64 `json:"hdd_percent"`
	NetMbps float64 `json:"net_mbps"`
	LatencyMs float64 `json:"latency_ms"`
	PacketLossPercent float64 `json:"packet_loss_percent"`
	Healthy bool `json:"healthy"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *App) nodeCatalog(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) { return }
	if a.db == nil { jsonResponse(w, http.StatusOK, map[string]any{"nodes": []Node{}}); return }
	rows, err := a.db.Query(r.Context(), `SELECT id,provider,region,services,cpu_percent,ram_percent,ssd_percent,hdd_percent,net_mbps,latency_ms,packet_loss_percent,healthy,updated_at FROM control_nodes ORDER BY healthy DESC, latency_ms ASC, id ASC`)
	if err != nil { jsonResponse(w, 500, map[string]string{"error":"nodes_query_failed"}); return }
	defer rows.Close()
	out := make([]Node, 0)
	for rows.Next() {
		var n Node
		if err := rows.Scan(&n.ID,&n.Provider,&n.Region,&n.Services,&n.CPUPercent,&n.RAMPercent,&n.SSDPercent,&n.HDDPercent,&n.NetMbps,&n.LatencyMs,&n.PacketLossPercent,&n.Healthy,&n.UpdatedAt); err != nil { jsonResponse(w,500,map[string]string{"error":"nodes_decode_failed"}); return }
		out = append(out,n)
	}
	if err := rows.Err(); err != nil { jsonResponse(w,500,map[string]string{"error":"nodes_query_failed"}); return }
	jsonResponse(w,200,map[string]any{"nodes":out})
}

func (a *App) registerNode(w http.ResponseWriter, r *http.Request) {
	if !method(w,r,http.MethodPost) { return }
	var n Node
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil { jsonResponse(w,400,map[string]string{"error":"invalid_json"}); return }
	n.ID = strings.TrimSpace(n.ID); n.Provider = strings.TrimSpace(n.Provider)
	if n.ID == "" || n.Provider == "" { jsonResponse(w,400,map[string]string{"error":"id_and_provider_required"}); return }
	if n.Services == nil { n.Services=[]string{} }
	n.UpdatedAt=time.Now().UTC()
	if a.db == nil { jsonResponse(w,202,map[string]any{"status":"accepted","node":n,"database":"disabled"}); return }
	_,err:=a.db.Exec(r.Context(),`INSERT INTO control_nodes(id,provider,region,services,cpu_percent,ram_percent,ssd_percent,hdd_percent,net_mbps,latency_ms,packet_loss_percent,healthy,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) ON CONFLICT(id) DO UPDATE SET provider=EXCLUDED.provider,region=EXCLUDED.region,services=EXCLUDED.services,cpu_percent=EXCLUDED.cpu_percent,ram_percent=EXCLUDED.ram_percent,ssd_percent=EXCLUDED.ssd_percent,hdd_percent=EXCLUDED.hdd_percent,net_mbps=EXCLUDED.net_mbps,latency_ms=EXCLUDED.latency_ms,packet_loss_percent=EXCLUDED.packet_loss_percent,healthy=EXCLUDED.healthy,updated_at=EXCLUDED.updated_at`,n.ID,n.Provider,n.Region,n.Services,n.CPUPercent,n.RAMPercent,n.SSDPercent,n.HDDPercent,n.NetMbps,n.LatencyMs,n.PacketLossPercent,n.Healthy,n.UpdatedAt)
	if err!=nil { jsonResponse(w,500,map[string]string{"error":"node_persist_failed"}); return }
	jsonResponse(w,202,map[string]any{"status":"registered","node":n})
}
