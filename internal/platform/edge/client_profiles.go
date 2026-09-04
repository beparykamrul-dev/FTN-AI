package edge

// PC and TV are client platforms used by the network-profile layer. The
// canonical ClientPlatform type is declared by client_enrollment.go.
const (
	PlatformPC ClientPlatform = "pc"
	PlatformTV ClientPlatform = "tv"
)

type ClientNetworkProfile struct {
	DeviceID string `json:"device_id"`
	Platform ClientPlatform `json:"platform"`
	UserID string `json:"user_id"`
	FTNNode string `json:"ftn_node"`
	IP string `json:"ip,omitempty"`
	MAC string `json:"mac,omitempty"`
	DNS string `json:"dns,omitempty"`
	Online bool `json:"online"`
}
