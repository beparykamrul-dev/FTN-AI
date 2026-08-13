package global

// PortalDNSProvider is the normalized registry entry for a DNS portal or
// managed DNS service. The registry keeps provider integration separate from
// FTN's authoritative DNS control plane.
type PortalDNSProvider struct {
	Name      string
	Kind      string
	Enabled   bool
	Health    ProviderHealth
	Endpoint  string
}

// PortalCandidates returns enabled, healthy portal providers. It deliberately
// does not contain credentials or provider-specific secrets.
func PortalCandidates(providers []PortalDNSProvider) []PortalDNSProvider {
	out := make([]PortalDNSProvider, 0, len(providers))
	for _, p := range providers {
		if p.Enabled && p.Health.Healthy() {
			out = append(out, p)
		}
	}
	return out
}
