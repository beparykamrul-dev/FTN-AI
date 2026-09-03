package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// QKDNode describes QKD/KMS inventory metadata only. Raw key material is never
// represented by this contract.
type QKDNode struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	Vendor string `json:"vendor,omitempty"`
	Model string `json:"model,omitempty"`
	EndpointRef string `json:"endpoint_ref,omitempty"`
	KMSRef string `json:"kms_ref,omitempty"`
	Status string `json:"status"`
}

type QKDLink struct {
	ID string `json:"id"`
	SourceNodeID string `json:"source_node_id"`
	TargetNodeID string `json:"target_node_id"`
	QuantumChannel string `json:"quantum_channel"`
	ClassicalChannel string `json:"classical_channel"`
	Authenticated bool `json:"authenticated"`
	Status string `json:"status"`
}

type QKDStatus struct {
	NodeID string `json:"node_id"`
	KMSRef string `json:"kms_ref,omitempty"`
	PoolState string `json:"pool_state"`
	AvailableKeys uint64 `json:"available_keys"`
	GenerationRateBps uint64 `json:"generation_rate_bps"`
	ConsumptionRateBps uint64 `json:"consumption_rate_bps"`
	Healthy bool `json:"healthy"`
	FallbackMode string `json:"fallback_mode"`
}

type QKDIntent struct {
	Action string `json:"action"`
	NodeID string `json:"node_id,omitempty"`
	Consumer string `json:"consumer,omitempty"`
	Policy string `json:"policy,omitempty"`
	KMSRef string `json:"kms_ref,omitempty"`
}

func (a *App) qkdNodes(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !requirePermission(a, "qkd.read", w, r) { return }
	jsonResponse(w, http.StatusOK, map[string]any{"nodes": []QKDNode{}})
}

func (a *App) qkdLinks(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !requirePermission(a, "qkd.read", w, r) { return }
	jsonResponse(w, http.StatusOK, map[string]any{"links": []QKDLink{}})
}

func (a *App) qkdStatus(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !requirePermission(a, "qkd.read", w, r) { return }
	jsonResponse(w, http.StatusOK, map[string]any{"status": []QKDStatus{}})
}

func (a *App) qkdKMS(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) || !requirePermission(a, "qkd.read", w, r) { return }
	jsonResponse(w, http.StatusOK, map[string]any{"kms": []map[string]any{}})
}

func (a *App) qkdIntent(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !requirePermission(a, "qkd.change", w, r) { return }
	var req QKDIntent
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"}); return
	}
	req.Action = strings.TrimSpace(req.Action)
	req.NodeID = strings.TrimSpace(req.NodeID)
	req.Consumer = strings.TrimSpace(req.Consumer)
	req.Policy = strings.TrimSpace(req.Policy)
	req.KMSRef = strings.TrimSpace(req.KMSRef)
	if req.Action == "" { jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "action_required"}); return }
	if req.KMSRef == "" && req.NodeID == "" { jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "node_or_kms_required"}); return }
	// QKD changes are intent-only. Execution must use the existing approval and
	// durable-job path; this endpoint never accepts or stores raw key material.
	a.audit(r, "qkd.intent", req.NodeID, "accepted", req)
	jsonResponse(w, http.StatusAccepted, map[string]any{"status": "intent_accepted", "execution": "approval_gated", "raw_key_material": false})
}
