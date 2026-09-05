package diagnostics

import (
	"strings"
	"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"
)

// Service coordinates evidence-backed diagnosis. It has no privileged execution path.
type Service struct{}

func (s Service) Diagnose(i model.Incident) model.Diagnosis {
	i=i.Normalize(); if !i.Valid(){return model.Diagnosis{Confidence:0}}
	evidence:=make([]string,0,len(i.EvidenceIDs)); seen:=map[string]struct{}{}; for _,id:=range i.EvidenceIDs{id=strings.TrimSpace(id);if id==""{continue};if _,ok:=seen[id];ok{continue};seen[id]=struct{}{};evidence=append(evidence,id)}
	return model.Diagnosis{IncidentID:i.ID,Confidence:0,Evidence:evidence}
}
