package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// NetworkIntegrationRequest is the common control-plane contract used by
// device, telemetry and routing adapters. It describes intent only; adapters
// must perform their own capability and safety validation before execution.
type NetworkIntegrationRequest struct {
	DeviceID   string          `json:"device_id,omitempty"`
	DeviceType string          `json:"device_type,omitempty"`
	Protocol   string          `json:"protocol,omitempty"`
	Action     string          `json:"action"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}

type NetworkIntegrationResponse struct {
	Status    string `json:"status"`
	DeviceID  string `json:"device_id,omitempty"`
	Protocol  string `json:"protocol,omitempty"`
	Action    string `json:"action,omitempty"`
	Execution string `json:"execution"`
}

func validNetworkProtocol(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "routeros-api", "snmp", "netflow", "ipfix", "bgp", "ospf", "bfd", "vrf", "ecmp", "olt", "onu":
		return true
	default:
		return false
	}
}

func (a *App) networkIntegration(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !requirePermission(a, "network.intent", w, r) {
		return
	}
	var req NetworkIntegrationRequest
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10))
	if err := dec.Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"})
		return
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "multiple_json_values"})
		return
	}
	req.DeviceID, req.DeviceType, req.Protocol, req.Action = strings.TrimSpace(req.DeviceID), strings.TrimSpace(req.DeviceType), strings.ToLower(strings.TrimSpace(req.Protocol)), strings.TrimSpace(req.Action)
	if req.Action == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "action_required"})
		return
	}
	if !validNetworkProtocol(req.Protocol) {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "unsupported_protocol"})
		return
	}
	if len(req.Payload) > 0 && !json.Valid(req.Payload) {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid_payload"})
		return
	}
	// This endpoint creates an adapter intent only. It never mutates a router,
	// OLT or route table directly. Privileged execution must go through the
	// approval + durable-job path.
	jsonResponse(w, http.StatusAccepted, NetworkIntegrationResponse{
		Status: "intent_accepted", DeviceID: req.DeviceID, Protocol: req.Protocol,
		Action: req.Action, Execution: "approval_gated",
	})
}
