package network

type VLAN struct {
	ID       string `json:"id"`
	DeviceID string `json:"device_id"`
	VID      uint16 `json:"vid"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}
