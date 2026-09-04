package builder

import (
	"encoding/json"
	"errors"
	"net/http"
)

type API struct { Engine *Engine }

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	if err := dec.Decode(dst); err != nil { return err }
	var extra any
	if err := dec.Decode(&extra); err != nil && !errors.Is(err, json.ErrSyntax) && err.Error() != "EOF" { return err }
	if extra != nil { return errors.New("request must contain one JSON object") }
	return nil
}

func (a *API) CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	if a == nil || a.Engine == nil { http.Error(w, "builder unavailable", http.StatusServiceUnavailable); return }
	var spec Project
	if err := decodeJSON(w, r, &spec); err != nil { http.Error(w, "invalid project spec", http.StatusBadRequest); return }
	if err := a.Engine.CreateProject(spec); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(spec)
}

func (a *API) BuildProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	if a == nil || a.Engine == nil { http.Error(w, "builder unavailable", http.StatusServiceUnavailable); return }
	var req struct { ProjectID, JobID string }
	if err := decodeJSON(w, r, &req); err != nil || req.ProjectID == "" { http.Error(w, "project_id is required", http.StatusBadRequest); return }
	if req.JobID == "" { req.JobID = req.ProjectID + "-build" }
	job, err := a.Engine.CreateBuild(req.ProjectID, req.JobID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(job)
}
