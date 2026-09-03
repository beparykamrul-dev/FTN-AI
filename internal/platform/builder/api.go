package builder

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type API struct { Engine *Engine }

func (a *API) CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	var spec Project
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil { http.Error(w, "invalid project spec", http.StatusBadRequest); return }
	if err := a.Engine.CreateProject(spec); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(spec)
}

func (a *API) BuildProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	var req struct { ProjectID string `json:"project_id"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProjectID == "" { http.Error(w, "project_id is required", http.StatusBadRequest); return }
	job, err := a.Engine.CreateBuild(req.ProjectID, uuid.NewString())
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(job)
}
