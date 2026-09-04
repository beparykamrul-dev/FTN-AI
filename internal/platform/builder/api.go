package builder

import (
	"encoding/json"
	"net/http"
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
	var req struct { ProjectID, JobID string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProjectID == "" { http.Error(w, "project_id is required", http.StatusBadRequest); return }
	if req.JobID == "" { req.JobID = req.ProjectID + "-build" }
	job, err := a.Engine.CreateBuild(req.ProjectID, req.JobID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(job)
}
