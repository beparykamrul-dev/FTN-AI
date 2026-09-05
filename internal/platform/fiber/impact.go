package fiber

import("sort";"strings")
type ImpactNode struct{ID string `json:"id"`;Kind string `json:"kind"`;Name string `json:"name,omitempty"`}
type FiberPath struct{ID string `json:"id"`;Segments []string `json:"segments"`;Impacted []ImpactNode `json:"impacted"`}
func BuildImpact(path FiberPath)[]ImpactNode{seen:=make(map[string]ImpactNode,len(path.Impacted));for _,n:=range path.Impacted{n.ID=strings.TrimSpace(n.ID);n.Kind=strings.TrimSpace(n.Kind);n.Name=strings.TrimSpace(n.Name);if n.ID!=""{if old,ok:=seen[n.ID];!ok||n.Kind<old.Kind||(n.Kind==old.Kind&&n.Name<old.Name){seen[n.ID]=n}}};out:=make([]ImpactNode,0,len(seen));for _,n:=range seen{out=append(out,n)};sort.Slice(out,func(i,j int)bool{if out[i].ID!=out[j].ID{return out[i].ID<out[j].ID};if out[i].Kind!=out[j].Kind{return out[i].Kind<out[j].Kind};return out[i].Name<out[j].Name});return out}
