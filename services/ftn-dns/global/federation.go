package global

import "sort"

// DNSProvider describes an external or FTN DNS implementation exposed through
// a common adapter. It is a federation registry, not a forced traffic policy.
type DNSProvider struct {
	Name       string
	Kind       string
	Enabled    bool
	Priority   int
	Health     ProviderHealth
}

// FederatedCandidates returns enabled, healthy DNS providers ordered by
// operational score and then configured priority. FTN-owned providers can be
// included alongside external providers without changing the adapter contract.
func FederatedCandidates(providers []DNSProvider) []DNSProvider {
	out := make([]DNSProvider, 0, len(providers))
	for _, provider := range providers {
		if provider.Enabled && provider.Health.Healthy() {
			out = append(out, provider)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		si, sj := out[i].Health.Score(), out[j].Health.Score()
		if si != sj {
			return si > sj
		}
		if out[i].Priority != out[j].Priority {
			return out[i].Priority < out[j].Priority
		}
		return out[i].Name < out[j].Name
	})
	return out
}
