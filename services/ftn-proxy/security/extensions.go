package security

// Capability describes an optional security control that can be enabled later
// without changing the proxy's core request path.
type Capability string

const (
	CapabilityWAF            Capability = "waf"
	CapabilityBotProtection  Capability = "bot-protection"
	CapabilityDDoSProtection Capability = "ddos-protection"
	CapabilityThreatIntel    Capability = "threat-intelligence"
	CapabilitySIEMExport     Capability = "siem-export"
	CapabilityMFAContext     Capability = "mfa-context"
	CapabilityDeviceTrust    Capability = "device-trust"
	CapabilityFraudSignals   Capability = "fraud-signals"
)

type SecurityProfile struct {
	Capabilities map[Capability]bool
}

func NewSecurityProfile() SecurityProfile {
	return SecurityProfile{Capabilities: map[Capability]bool{}}
}

func (p SecurityProfile) Enabled(c Capability) bool {
	return p.Capabilities != nil && p.Capabilities[c]
}
