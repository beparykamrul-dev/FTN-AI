package deploy

import("errors";"strings";"time")

type Plan struct{ID string `json:"id"`;ProjectID string `json:"project_id"`;TargetID string `json:"target_id"`;Artifact string `json:"artifact"`;Strategy string `json:"strategy"`;CreatedAt time.Time `json:"created_at"`}

func NewPlan(id,projectID string,target Target,artifact,strategy string)(Plan,error){
	id=strings.TrimSpace(id);projectID=strings.TrimSpace(projectID);artifact=strings.TrimSpace(artifact);strategy=strings.ToLower(strings.TrimSpace(strategy))
	if err:=target.Validate();err!=nil{return Plan{},err}
	if id==""||projectID==""||artifact==""{return Plan{},errors.New("id, project_id and artifact are required")}
	if strategy==""{strategy="rolling"}
	switch strategy{case "rolling","recreate","blue-green":default:return Plan{},errors.New("unsupported deployment strategy")}
	return Plan{ID:id,ProjectID:projectID,TargetID:strings.TrimSpace(target.ID),Artifact:artifact,Strategy:strategy,CreatedAt:time.Now().UTC()},nil
}
