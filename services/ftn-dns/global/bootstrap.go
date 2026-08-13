package global

import "strings"

// BootstrapConfig defines the provider-neutral bootstrap inventory used by
// FTN DNS. Provider names are labels only; actual endpoints and credentials
// belong in environment-specific configuration/secrets.
type BootstrapConfig struct {
	Providers []DNSProvider
	Portals   []PortalDNSProvider
}

// NormalizeName provides deterministic provider identifiers for registry use.
func NormalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// EnabledProviders returns the federation providers that are enabled and
// healthy. It is intentionally read-only and safe for the control plane.
func (c BootstrapConfig) EnabledProviders() []DNSProvider {
	return FederatedCandidates(c.Providers)
}

// EnabledPortals returns healthy DNS portal integrations.
func (c BootstrapConfig) EnabledPortals() []PortalDNSProvider {
	return PortalCandidates(c.Portals)
}
