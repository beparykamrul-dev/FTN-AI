package gis

import (
    "encoding/json"
    "io"
    "net/http"
)

type FiberAPI struct {
    Map *FiberMap
    Hub *Hub
}

func (a *FiberAPI) Assets(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
    if a.Map == nil { http.Error(w, "fiber map unavailable", http.StatusServiceUnavailable); return }
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(a.Map.List()); err != nil { return }
}

func (a *FiberAPI) Upsert(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost && r.Method != http.MethodPut { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
    if a.Map == nil { http.Error(w, "fiber map unavailable", http.StatusServiceUnavailable); return }
    r.Body = http.MaxBytesReader(w, r.Body, 256<<10)
    var asset FiberAsset
    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&asset); err != nil { http.Error(w, "invalid fiber asset", http.StatusBadRequest); return }
    var extra any
    if err := dec.Decode(&extra); err != io.EOF { http.Error(w, "multiple json values", http.StatusBadRequest); return }
    if err := a.Map.Upsert(asset); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    if a.Hub != nil { _ = a.Hub.Publish(Event{Type: "fiber.asset.updated", Node: &MapNode{ID: asset.ID, Name: asset.Name, Kind: string(asset.Type)}}) }
    w.WriteHeader(http.StatusNoContent)
}
