package main

import "net/http"

type apiCapabilities struct {
	Version string   `json:"version"`
	Modules []string `json:"modules"`
}

func (a *App) apiRoot(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"name":    "FTN Control Plane API",
		"version": "v1",
		"status":  "ok",
		"request_id": requestInfo(r).RequestID,
		"correlation_id": requestInfo(r).CorrelationID,
	})
}

func (a *App) apiMe(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	rc := requestInfo(r)
	jsonResponse(w, http.StatusOK, map[string]any{
		"principal_id":   rc.PrincipalID,
		"tenant_id":      rc.TenantID,
		"request_id":     rc.RequestID,
		"correlation_id": rc.CorrelationID,
	})
}

func (a *App) apiCapabilities(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	if !requirePermission(a, "service.read", w, r) {
		return
	}
	jsonResponse(w, http.StatusOK, apiCapabilities{
		Version: "v1",
		Modules: []string{
			"auth", "rbac", "services", "devices", "network", "dns",
			"jobs", "approvals", "audit", "monitoring", "ai", "control-panel",
		},
	})
}

func (a *App) apiDashboard(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	if !requirePermission(a, "service.read", w, r) {
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"status": "ok",
		"services": map[string]any{
			"total":    len(a.catalog),
			"available": len(a.catalog),
		},
		"control_plane": map[string]any{
			"api":        "online",
			"approval_first": true,
			"audit":      true,
		},
		"request_id": requestInfo(r).RequestID,
	})
}
