package main

import (
    "encoding/json"
    "net/http"
    "strings"
    "sync"
    "time"
)

type RegisteredDevice struct {
    NetworkDevice
    CredentialRef string `json:"credential_ref,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type DeviceRegistry struct {
    mu sync.RWMutex
    devices map[string]RegisteredDevice
    adapters *AdapterRegistry
}

func NewDeviceRegistry(adapters *AdapterRegistry) *DeviceRegistry {
    return &DeviceRegistry{devices: make(map[string]RegisteredDevice), adapters: adapters}
}

func (r *DeviceRegistry) Upsert(d RegisteredDevice) bool {
    if strings.TrimSpace(d.ID) == "" || strings.TrimSpace(d.Address) == "" || strings.TrimSpace(d.Protocol) == "" { return false }
    now := time.Now().UTC()
    r.mu.Lock(); defer r.mu.Unlock()
    old, exists := r.devices[d.ID]
    if exists { d.CreatedAt = old.CreatedAt } else { d.CreatedAt = now }
    d.UpdatedAt = now
    r.devices[d.ID] = d
    return true
}

func (r *DeviceRegistry) Get(id string) (RegisteredDevice, bool) {
    r.mu.RLock(); defer r.mu.RUnlock(); d, ok := r.devices[id]; return d, ok
}

func (r *DeviceRegistry) List() []RegisteredDevice {
    r.mu.RLock(); defer r.mu.RUnlock()
    out := make([]RegisteredDevice, 0, len(r.devices))
    for _, d := range r.devices { out = append(out, d) }
    return out
}

func (a *App) deviceRegistryAPI(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/api/v1/devices" && r.URL.Path != "/api/v1/devices/" { http.NotFound(w,r); return }
    switch r.Method {
    case http.MethodGet:
        jsonResponse(w, http.StatusOK, map[string]any{"devices": a.devices.List()})
    case http.MethodPost:
        var d RegisteredDevice
        if err := json.NewDecoder(r.Body).Decode(&d); err != nil { jsonResponse(w,400,map[string]string{"error":"invalid_json"}); return }
        d.Protocol = normalizeProtocol(d.Protocol)
        if !validNetworkProtocol(d.Protocol) { jsonResponse(w,400,map[string]string{"error":"unsupported_protocol"}); return }
        if !validNode(NetworkDevice{ID:d.ID,Address:d.Address,Kind:d.Kind,Protocol:d.Protocol}) { jsonResponse(w,400,map[string]string{"error":"invalid_device"}); return }
        if !isFTNDeviceKind(d.Kind) { jsonResponse(w,403,map[string]string{"error":"device_ownership_not_verified"}); return }
        if !a.devices.Upsert(d) { jsonResponse(w,400,map[string]string{"error":"invalid_device"}); return }
        jsonResponse(w,http.StatusAccepted,map[string]any{"status":"registered","device":d})
    default:
        jsonResponse(w,http.StatusMethodNotAllowed,map[string]string{"error":"method_not_allowed"})
    }
}

func (a *App) deviceTelemetry(w http.ResponseWriter, r *http.Request) {
    if !method(w,r,http.MethodGet) { return }
    id := strings.TrimPrefix(r.URL.Path, "/api/v1/devices/")
    id = strings.TrimSuffix(id, "/telemetry")
    id = strings.Trim(id, "/")
    if id == "" { jsonResponse(w,400,map[string]string{"error":"device_id_required"}); return }
    d, ok := a.devices.Get(id); if !ok { jsonResponse(w,404,map[string]string{"error":"device_not_found"}); return }
    adapter, ok := a.adapters.Adapter(d.Protocol); if !ok { jsonResponse(w,409,map[string]string{"error":"adapter_not_registered"}); return }
    interfaces, err := adapter.CollectInterfaceState(r.Context(), d.NetworkDevice); if err != nil { jsonResponse(w,502,map[string]string{"error":"interface_collection_failed"}); return }
    routes, err := adapter.CollectRoutingState(r.Context(), d.NetworkDevice); if err != nil { jsonResponse(w,502,map[string]string{"error":"routing_collection_failed"}); return }
    jsonResponse(w,200,map[string]any{"device_id":id,"protocol":d.Protocol,"interfaces":interfaces,"routes":routes,"read_only":true,"collected_at":time.Now().UTC()})
}
