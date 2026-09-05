package ai

import("math";"strings";"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model")
type Recommendation struct{Summary string `json:"summary"`;Action string `json:"action,omitempty"`;Confidence float64 `json:"confidence"`;Evidence []string `json:"evidence,omitempty"`;ApprovalRequired bool `json:"approval_required"`}
type Advisor struct{}
func(Advisor)Advise(d model.Diagnosis)Recommendation{confidence:=d.Confidence;if math.IsNaN(confidence)||math.IsInf(confidence,0)||confidence<0||confidence>1{confidence=0};evidence:=make([]string,0,len(d.Evidence));seen:=map[string]struct{}{};for _,e:=range d.Evidence{e=strings.TrimSpace(e);if e==""||len(e)>2048{continue};if _,ok:=seen[e];ok{continue};seen[e]=struct{}{};evidence=append(evidence,e)};cause:=strings.TrimSpace(d.Cause);if len(cause)>4096{cause=cause[:4096]};if cause==""{return Recommendation{Summary:"Insufficient evidence for a reliable root-cause recommendation.",Confidence:confidence,Evidence:evidence,ApprovalRequired:true}};return Recommendation{Summary:cause,Action:"Review evidence and apply an approved remediation policy.",Confidence:confidence,Evidence:evidence,ApprovalRequired:true}}
