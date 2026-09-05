package correlate

import("sort";"strings";"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model")
type Correlator struct{}
func(Correlator)Collect(incident model.Incident,evidence []model.Evidence)[]model.Evidence{allowed:=make(map[string]struct{},len(incident.EvidenceIDs));for _,id:=range incident.EvidenceIDs{id=strings.TrimSpace(id);if id!=""{allowed[id]=struct{}{}}};if len(allowed)==0{return nil};out:=make([]model.Evidence,0,len(evidence));seen:=make(map[string]struct{},len(evidence));for _,item:=range evidence{item=item.Normalize();item.ID=strings.TrimSpace(item.ID);if _,ok:=allowed[item.ID];!ok||!item.Valid(){continue};if _,ok:=seen[item.ID];ok{continue};seen[item.ID]=struct{}{};out=append(out,item)};sort.SliceStable(out,func(i,j int)bool{if out[i].Timestamp!=out[j].Timestamp{return out[i].Timestamp<out[j].Timestamp};return out[i].ID<out[j].ID});return out}
