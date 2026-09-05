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
	candidates:=make([]Candidate,0); for _,item:=range evidence{if strings.TrimSpace(item.Summary)==""||strings.TrimSpace(item.ID)==""{continue};candidates=append(candidates,Candidate{Cause:item.Summary,Confidence:0.25,Evidence:[]string{item.ID}})}
	for _,dep:=range dependencies{if dep.Critical&&!dep.Healthy{candidates=append(candidates,Candidate{Cause:"critical dependency degraded",Confidence:0.5,Impact:[]string{dep.From,dep.To}})}}
	sort.SliceStable(candidates,func(i,j int)bool{return candidates[i].Confidence>candidates[j].Confidence});return candidates
}
func (Analyzer) Validate(candidates []Candidate,evidence []model.Evidence) model.Diagnosis {
	if len(candidates)==0{return model.Diagnosis{Confidence:0}}
	validEvidence:=make(map[string]struct{},len(evidence));for _,e:=range evidence{if e.ID!=""{validEvidence[e.ID]=struct{}{}}}
	for _,best:=range candidates{if strings.TrimSpace(best.Cause)==""||math.IsNaN(best.Confidence)||math.IsInf(best.Confidence,0)||best.Confidence<0||best.Confidence>1{continue};bound:=make([]string,0,len(best.Evidence));for _,id:=range best.Evidence{if _,ok:=validEvidence[id];ok{bound=append(bound,id)}};if len(best.Evidence)>0&&len(bound)==0{continue};return model.Diagnosis{Cause:best.Cause,Confidence:best.Confidence,Evidence:bound,Impact:append([]string(nil),best.Impact...)}}
	return model.Diagnosis{Confidence:0}
}
