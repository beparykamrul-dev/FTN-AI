package monitoring

import("sort";"strings";"github.com/beparykamrul-dev/FTN-AI/internal/platform/module")
func Definition()module.Definition{caps:=[]string{"metrics","events","health","telemetry","alerts","observability"};sort.Strings(caps);return module.Definition{Name:"monitoring",Version:"v1",Capabilities:caps}}
func ValidDefinition(d module.Definition)bool{if strings.TrimSpace(d.Name)!="monitoring"||strings.TrimSpace(d.Version)==""||len(d.Capabilities)==0{return false};seen:=map[string]struct{}{};for _,c:=range d.Capabilities{c=strings.TrimSpace(c);if c==""{return false};if _,ok:=seen[c];ok{return false};seen[c]=struct{}{}};return true}
