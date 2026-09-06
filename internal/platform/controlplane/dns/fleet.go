package dns
import("math";"sort";"time")
type Fleet struct{Servers []ServerHealth `json:"servers"`}
type RankedServer struct{ServerHealth;Score float64 `json:"score"`}
func(f Fleet)Rank(now time.Time,maxAge time.Duration)[]RankedServer{out:=make([]RankedServer,0,len(f.Servers));for _,s:=range f.Servers{if !Available(s,now,maxAge){continue};latencyScore:=100.0-min(s.LatencyMs,100.0);servfailScore:=100.0-s.ServfailRate;score:=latencyScore*.55+servfailScore*.25+min(s.QPS/100.0,100.0)*.20;if math.IsNaN(score)||math.IsInf(score,0){continue};out=append(out,RankedServer{ServerHealth:s,Score:score})};sort.SliceStable(out,func(i,j int)bool{if out[i].Score!=out[j].Score{return out[i].Score>out[j].Score};return out[i].NodeID<out[j].NodeID});return out}
func min(a,b float64)float64{if a<b{return a};return b}
