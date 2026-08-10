package mesh

// Protocol identifies an optional routing or mesh backend. FTN keeps these as
// adapters so the canonical topology and policy model remains vendor-neutral.
type Protocol string

const (
    ProtocolOLSR Protocol = "olsr"
    ProtocolYggdrasil Protocol = "yggdrasil"
    ProtocolWireGuard Protocol = "wireguard"
    ProtocolNetBird Protocol = "netbird"
    ProtocolHeadscale Protocol = "headscale"
    ProtocolTailscale Protocol = "tailscale"
)

type ProtocolCapability struct {
    Protocol Protocol `json:"protocol"`
    Capabilities []string `json:"capabilities"`
    Enabled bool `json:"enabled"`
}

// DefaultProtocolCapabilities documents the FTN adapter boundary. It does not
// install or execute any third-party daemon.
func DefaultProtocolCapabilities() []ProtocolCapability {
    return []ProtocolCapability{
        {Protocol: ProtocolOLSR, Capabilities: []string{"link_state", "dynamic_routing"}},
        {Protocol: ProtocolYggdrasil, Capabilities: []string{"encrypted_mesh", "overlay_routing"}},
        {Protocol: ProtocolWireGuard, Capabilities: []string{"encrypted_tunnel", "peer_management"}},
        {Protocol: ProtocolNetBird, Capabilities: []string{"mesh_management", "peer_management"}},
        {Protocol: ProtocolHeadscale, Capabilities: []string{"tailnet_control", "identity_policy"}},
        {Protocol: ProtocolTailscale, Capabilities: []string{"tailnet_client", "peer_connectivity"}},
    }
}
