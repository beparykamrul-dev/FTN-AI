package main

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
)

func (a *App) nodeHeartbeat(w http.ResponseWriter, r *http.Request) {
    if !method(w, r, http.MethodPost) { return }
    var n Node
    if err := json.NewDecoder(r.Body).Decode(&n); err != nil { jsonResponse(w, http.StatusBadRequest, map[string]string{"error":"invalid_json"}); return }
    if !validNode(n) { jsonResponse(w, http.StatusBadRequest, map[string]string{"error":"invalid_node"}); return }
    n.LastSeen = time.Now().UTC()
    if a.db == nil { upsertMemoryNode(n); jsonResponse(w,http.StatusOK,map[string]any{"status":"heartbeat_accepted","node_id":n.ID,"last_seen":n.LastSeen}); return }
    ctx,cancel := context.WithTimeout(r.Context(),2*time.Second); defer cancel()
    result,err := a.db.Exec(ctx, `update control_nodes set provider=$2,region=$3,services=$4,cpu_percent=$5,ram_percent=$6,ssd_percent=$7,hdd_percent=$8,net_mbps=$9,latency_ms=$10,packet_loss_percent=$11,jitter_ms=$12,retransmissions=$13,capacity_mbps=$14,utilization_percent=$15,bgp_up=$16,bfd_up=$17,isis_up=$18,evpn_ready=$19,anycast_ready=$20,rpki_valid=$21,prefix_count=$22,healthy=$23,updated_at=now() where id=$1`, n.ID,n.Provider,n.Region,n.Services,n.CPUPercent,n.RAMPercent,n.SSDPercent,n.HDDPercent,n.NetMbps,n.LatencyMs,n.PacketLoss,n.JitterMs,n.Retransmissions,n.CapacityMbps,n.UtilizationPercent,n.BGPUp,n.BFDUp,n.ISISUp,n.EVPNReady,n.AnycastReady,n.RPKIValid,n.PrefixCount,n.Healthy)
    if err != nil { jsonResponse(w,http.StatusInternalServerError,map[string]string{"error":"heartbeat_update_failed"}); return }
    if result.RowsAffected()==0 { jsonResponse(w,http.StatusNotFound,map[string]string{"error":"node_not_registered"}); return }
    jsonResponse(w,http.StatusOK,map[string]any{"status":"heartbeat_accepted","node_id":n.ID,"last_seen":n.LastSeen})
}
