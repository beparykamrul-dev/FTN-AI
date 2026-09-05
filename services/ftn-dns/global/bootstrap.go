package global

import "strings"
type BootstrapConfig struct{Providers []DNSProvider;Portals []PortalDNSProvider}
func NormalizeName(name string)string{return strings.ToLower(strings.TrimSpace(name))}
func(c BootstrapConfig)EnabledProviders()[]DNSProvider{return FederatedCandidates(c.Providers)}
func(c BootstrapConfig)EnabledPortals()[]PortalDNSProvider{return PortalCandidates(c.Portals)}
func(c BootstrapConfig)Valid()bool{return len(c.EnabledProviders())>0||len(c.EnabledPortals())>0}
