package main

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strings"
)

type ExRouterRequest struct { ServiceID string `json:"service_id"`; Region string `json:"region,omitempty"`; Provider string `json:"provider,omitempty"` }
type ExRoute struct { PathID string `json:"path_id"`; Provider string `json:"provider"`; Region string `json:"region,omitempty"`; Score float64 `json:"score"`; LatencyMs float64 `json:"latency_ms"`; PacketLoss float64 `json:"packet_loss_percent"`; Healthy bool `json:"healthy"` }

func exRouteScore(n Node, req ExRouterRequest) float64 {
	if !n.Healthy || n.LatencyMs < 0 || n.PacketLoss < 0 || n.PacketLoss > 100 || math.IsNaN(n.LatencyMs) || math.IsInf(n.LatencyMs,0) || math.IsNaN(n.PacketLoss) || math.IsInf(n.PacketLoss,0) { return -1 }
	score := 100.0 - n.LatencyMs - n.PacketLoss*4.0
	if req.Region != "" && strings.EqualFold(req.Region,n.Region) { score += 15 }
	if req.Provider != "" && strings.EqualFold(req.Provider,n.Provider) { score += 8 }
	if math.IsNaN(score) || math.IsInf(score,0) { return -1 }
	return score
}

func (a *App) exRouter(w http.ResponseWriter,r *http.Request) {
	if !method(w,r,http.MethodPost)||!requirePermission(a,"network.route",w,r){return}
	var req ExRouterRequest
	dec:=json.NewDecoder(http.MaxBytesReader(w,r.Body,16<<10));if err:=dec.Decode(&req);err!=nil{jsonResponse(w,http.StatusBadRequest,map[string]string{"error":"invalid_json"});return};var extra any;if err:=dec.Decode(&extra);err!=nil{jsonResponse(w,http.StatusBadRequest,map[string]string{"error":"multiple_json_values"});return}
	req.ServiceID,req.Region,req.Provider=strings.TrimSpace(req.ServiceID),strings.TrimSpace(req.Region),strings.TrimSpace(req.Provider)
	if req.ServiceID==""{jsonResponse(w,http.StatusBadRequest,map[string]string{"error":"service_id_required"});return};if !serviceExists(a.catalog,req.ServiceID){jsonResponse(w,http.StatusNotFound,map[string]string{"error":"service_not_found"});return}
	rc:=requestInfo(r);ns,err:=a.loadNodes(r.Context(),rc.TenantID);if err!=nil{jsonResponse(w,http.StatusInternalServerError,map[string]string{"error":"node_query_failed"});return}
	routes:=make([]ExRoute,0,len(ns));for _,n:=range ns{score:=exRouteScore(n,req);if score<0{continue};routes=append(routes,ExRoute{PathID:n.ID,Provider:n.Provider,Region:n.Region,Score:score,LatencyMs:n.LatencyMs,PacketLoss:n.PacketLoss,Healthy:n.Healthy})}
	sort.SliceStable(routes,func(i,j int)bool{if routes[i].Score!=routes[j].Score{return routes[i].Score>routes[j].Score};return routes[i].PathID<routes[j].PathID})
	if len(routes)==0{jsonResponse(w,http.StatusServiceUnavailable,map[string]any{"status":"no_healthy_path","service_id":req.ServiceID});return}
	jsonResponse(w,http.StatusOK,map[string]any{"status":"path_ready","service_id":req.ServiceID,"selected_path":routes[0],"paths":routes,"service_unchanged":true,"execution":"policy_controlled"})
}
