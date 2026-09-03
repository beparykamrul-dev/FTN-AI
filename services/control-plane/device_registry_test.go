package main

import (
    "context"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

type registryTestAdapter struct{}
func (registryTestAdapter) Protocol() string { return "routeros-api" }
func (registryTestAdapter) Capabilities(context.Context, NetworkDevice) ([]string,error) { return []string{"health.read"},nil }
func (registryTestAdapter) CollectInterfaceState(context.Context, NetworkDevice) ([]InterfaceState,error) { return []InterfaceState{{DeviceID:"r1",Name:"ether1",OperUp:true}},nil }
func (registryTestAdapter) CollectRoutingState(context.Context, NetworkDevice) ([]RoutingState,error) { return []RoutingState{{DeviceID:"r1",Prefix:"0.0.0.0/0",Active:true}},nil }

func TestDeviceRegistryUpsertAndTelemetry(t *testing.T) {
    adapters := NewAdapterRegistry(); adapters.RegisterAdapter(registryTestAdapter{})
    app := &App{catalog:catalog, adapters:adapters, devices:NewDeviceRegistry(adapters)}
    req := httptest.NewRequest(http.MethodPost,"/api/v1/devices",strings.NewReader(`{"id":"r1","name":"core","kind":"router","address":"127.0.0.1","protocol":"routeros-api"}`))
    rec := httptest.NewRecorder(); app.deviceRegistryAPI(rec,req)
    if rec.Code != http.StatusAccepted { t.Fatalf("register status=%d body=%s",rec.Code,rec.Body.String()) }
    req = httptest.NewRequest(http.MethodGet,"/api/v1/devices/r1/telemetry",nil)
    rec = httptest.NewRecorder(); app.deviceTelemetry(rec,req)
    if rec.Code != http.StatusOK { t.Fatalf("telemetry status=%d body=%s",rec.Code,rec.Body.String()) }
}

func TestDeviceRegistryRejectsUnknownProtocol(t *testing.T) {
    adapters := NewAdapterRegistry(); app := &App{devices:NewDeviceRegistry(adapters),adapters:adapters}
    req := httptest.NewRequest(http.MethodPost,"/api/v1/devices",strings.NewReader(`{"id":"r1","kind":"router","address":"127.0.0.1","protocol":"unknown"}`))
    rec := httptest.NewRecorder(); app.deviceRegistryAPI(rec,req)
    if rec.Code != http.StatusBadRequest { t.Fatalf("status=%d",rec.Code) }
}
