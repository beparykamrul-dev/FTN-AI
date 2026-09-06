package controlplane

import("sort";"strings")

func NormalizeServices(in []string)[]string{seen:=make(map[string]struct{},len(in));out:=make([]string,0,len(in));for _,v:=range in{v=strings.TrimSpace(v);if v==""{continue};if _,ok:=seen[v];ok{continue};seen[v]=struct{}{};out=append(out,v)};sort.Strings(out);return out}
