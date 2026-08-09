package gis

import (
	"encoding/json"
	"net/http"
)

type FiberAPI struct {
	Map *FiberMap
	Hub *Hub
}

func (a *FiberAPI) Assets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(a.Map.List())
}

func (a *FiberAPI) Upsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	var asset FiberAsset
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil { http.Error(w, "invalid fiber asset", http.StatusBadRequest); return }
	if err := a.Map.Upsert(asset); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	if a.Hub != nil { _ = a.Hub.Publish(Event{Type: "fiber.asset.updated", Node: &MapNode{ID: asset.ID, Name: asset.Name, Kind: string(asset.Type)}}) }
	w.WriteHeader(http.StatusNoContent)
}
