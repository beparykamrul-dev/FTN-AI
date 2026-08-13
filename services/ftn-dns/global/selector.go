package global

import "sort"

// SelectProvider chooses the best healthy DNS provider without forcing a
// particular vendor. The caller supplies already-measured health states.
func SelectProvider(states []ProviderHealth) (ProviderHealth, bool) {
	candidates := make([]ProviderHealth, 0, len(states))
	for _, state := range states {
		if state.Healthy() {
			candidates = append(candidates, state)
		}
	}
	if len(candidates) == 0 {
		return ProviderHealth{}, false
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		si, sj := candidates[i].Score(), candidates[j].Score()
		if si != sj {
			return si > sj
		}
		return candidates[i].Provider < candidates[j].Provider
	})
	return candidates[0], true
}
