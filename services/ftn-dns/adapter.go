package dns
import("context";"math";"strings")
type Adapter interface{Name()string;Health(context.Context)(Health,error);Query(context.Context,string,string)(Response,error)}
type Health struct{Reachable bool;Secure bool;LatencyMS int64;LossRatio float64}
type Response struct{Name string;RecordType string;Values []string;Secure bool}
func(h Health)Valid()bool{return h.LatencyMS>=0&&h.LatencyMS<=24*60*60*1000&&!math.IsNaN(h.LossRatio)&&!math.IsInf(h.LossRatio,0)&&h.LossRatio>=0&&h.LossRatio<=1}
func(r Response)Valid()bool{n:=r.Normalized();return n.Name!=""&&n.RecordType!=""&&len(n.Values)>0&&len(n.Name)<=253&&len(n.RecordType)<=32}
func(r Response)Normalized()Response{r.Name=strings.TrimSuffix(strings.TrimSpace(r.Name),".");r.RecordType=strings.ToUpper(strings.TrimSpace(r.RecordType));values:=make([]string,0,len(r.Values));seen:=make(map[string]struct{},len(r.Values));for _,value:=range r.Values{value=strings.TrimSpace(value);if value==""||len(value)>65535{continue};if _,ok:=seen[value];ok{continue};seen[value]=struct{}{};values=append(values,value);if len(values)>=1024{break}};r.Values=values;return r}
