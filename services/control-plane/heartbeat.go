package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// nodeHeartbeat accepts a complete, current resource snapshot from an enrolled node.
// It intentionally does not perform deployment or privileged infrastructure changes.
func (a *App) nodeHeartbeat(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !requirePermission(a, "node.manage", w, r) {
		return
	}
	rc := requestInfo(r)
	if rc.TenantID == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "tenant_required"})
		return
	}
	var n Node
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10))
	if err := dec.Decode(&n); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"})
		return
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "multiple_json_values"})
		return
	}
	if !validNode(n) {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid_node"})
		return
	}

	n.TenantID = rc.TenantID
	n.LastSeen = time.Now().UTC()
	if a.db == nil {
		upsertMemoryNode(n)
		jsonResponse(w, http.StatusOK, map[string]any{"status": "heartbeat_accepted", "node_id": n.ID, "last_seen": n.LastSeen})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	result, err := a.db.Exec(ctx, `update control_nodes set provider=$3,region=$4,services=$5,cpu_percent=$6,ram_percent=$7,ssd_percent=$8,hdd_percent=$9,net_mbps=$10,latency_ms=$11,packet_loss_percent=$12,healthy=$13,updated_at=now() where id=$1 and tenant_id=$2::uuid`, n.ID, n.TenantID, n.Provider, n.Region, n.Services, n.CPUPercent, n.RAMPercent, n.SSDPercent, n.HDDPercent, n.NetMbps, n.LatencyMs, n.PacketLoss, n.Healthy)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "heartbeat_update_failed"})
		return
	}
	if result.RowsAffected() == 0 {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": "node_not_registered"})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"status": "heartbeat_accepted", "node_id": n.ID, "last_seen": n.LastSeen})
}
