package ftndns

import "sort"

type Upstream struct { ID string; Healthy bool; LatencyMS int64; Load float64 }

type UpstreamSelector struct{}

func (UpstreamSelector) Select(upstreams []Upstream) (Upstream, bool) {
	candidates := make([]Upstream,0,len(upstreams))
	for _, u := range upstreams { if u.Healthy { candidates=append(candidates,u) } }
	if len(candidates)==0 { return Upstream{},false }
	sort.SliceStable(candidates,func(i,j int) bool {
		if candidates[i].LatencyMS != candidates[j].LatencyMS { return candidates[i].LatencyMS < candidates[j].LatencyMS }
		if candidates[i].Load != candidates[j].Load { return candidates[i].Load < candidates[j].Load }
		return candidates[i].ID < candidates[j].ID
	})
	return candidates[0],true
}
