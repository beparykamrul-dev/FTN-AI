package rootcause

import (
	"math"
	"sort"
	"strings"
	"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"
)

type Candidate struct { Cause string; Confidence float64; Evidence []string; Impact []string }
type Analyzer struct{}

func (Analyzer) Find(evidence []model.Evidence, dependencies []model.Dependency) []Candidate {
	candidates := make([]Candidate,0)
	for _, item := range evidence { item=item.Normalize(); if !item.Valid() || item.Summary=="" { continue }; candidates=append(candidates,Candidate{Cause:item.Summary,Confidence:0.25,Evidence:[]string{item.ID}}) }
	for _, dep := range dependencies { dep=dep.Normalize(); if dep.Valid() && dep.Critical && !dep.Healthy { candidates=append(candidates,Candidate{Cause:"critical dependency degraded",Confidence:0.5,Impact:[]string{dep.From,dep.To}}) } }
	sort.SliceStable(candidates,func(i,j int)bool{if candidates[i].Confidence!=candidates[j].Confidence{return candidates[i].Confidence>candidates[j].Confidence};return candidates[i].Cause<candidates[j].Cause})
	return candidates
}

func (Analyzer) Validate(candidates []Candidate, evidence []model.Evidence) model.Diagnosis {
	validEvidence:=make(map[string]struct{},len(evidence)); for _,e:=range evidence { e=e.Normalize(); if e.Valid(){validEvidence[e.ID]=struct{}{}} }
	for _,best:=range candidates { cause:=strings.TrimSpace(best.Cause); if cause==""||math.IsNaN(best.Confidence)||math.IsInf(best.Confidence,0)||best.Confidence<0||best.Confidence>1{continue}; bound:=make([]string,0,len(best.Evidence)); seen:=map[string]struct{}{}; for _,id:=range best.Evidence{id=strings.TrimSpace(id);if _,ok:=validEvidence[id];ok{if _,dup:=seen[id];!dup{seen[id]=struct{}{};bound=append(bound,id)}}}; if len(best.Evidence)>0&&len(bound)==0{continue}; return model.Diagnosis{Cause:cause,Confidence:best.Confidence,Evidence:bound,Impact:append([]string(nil),best.Impact...)} }
	return model.Diagnosis{Confidence:0}
}
