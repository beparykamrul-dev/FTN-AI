package dns

// TrafficClass defines the traffic a self-hosted FTN DNS/Control Plane may
// legitimately need. Keeping these classes explicit makes firewall/QoS policy
// auditable instead of allowing an unrestricted "DNS server" exception.
type TrafficClass string

const (
	TrafficDNS TrafficClass = "dns"
	TrafficMeshControl TrafficClass = "mesh_control"
	TrafficTLS TrafficClass = "tls"
	TrafficMonitoring TrafficClass = "monitoring"
	TrafficReplication TrafficClass = "replication"
	TrafficManagement TrafficClass = "management"
)

type TrafficProfile struct {
	Name string `json:"name"`
	Classes []TrafficClass `json:"classes"`
	InboundDNS bool `json:"inbound_dns"`
	OutboundDNS bool `json:"outbound_dns"`
	Recursive bool `json:"recursive"`
	Authoritative bool `json:"authoritative"`
	CloudSync bool `json:"cloud_sync"`
}

func LocalAuthoritativeProfile() TrafficProfile {
	return TrafficProfile{
		Name: "ftn-local-authoritative",
		Classes: []TrafficClass{TrafficDNS, TrafficTLS, TrafficMonitoring, TrafficManagement},
		InboundDNS: true, OutboundDNS: false, Authoritative: true,
	}
}

func LocalRecursiveProfile() TrafficProfile {
	return TrafficProfile{
		Name: "ftn-local-recursive",
		Classes: []TrafficClass{TrafficDNS, TrafficTLS, TrafficMonitoring, TrafficManagement},
		InboundDNS: true, OutboundDNS: true, Recursive: true,
	}
}
