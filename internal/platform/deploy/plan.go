package deploy

import (
	"errors"
	"time"
)

type Plan struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	TargetID  string    `json:"target_id"`
	Artifact  string    `json:"artifact"`
	Strategy  string    `json:"strategy"`
	CreatedAt time.Time `json:"created_at"`
}

func NewPlan(id, projectID string, target Target, artifact, strategy string) (Plan, error) {
	if err := target.Validate(); err != nil { return Plan{}, err }
	if id == "" || projectID == "" || artifact == "" { return Plan{}, errors.New("id, project_id and artifact are required") }
	if strategy == "" { strategy = "rolling" }
	return Plan{ID: id, ProjectID: projectID, TargetID: target.ID, Artifact: artifact, Strategy: strategy, CreatedAt: time.Now().UTC()}, nil
}
