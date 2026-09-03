package edge

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
