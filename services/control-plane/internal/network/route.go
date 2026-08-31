package network

type Route struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id"`
	Prefix    string `json:"prefix"`
	NextHop   string `json:"next_hop"`
	Protocol  string `json:"protocol"`
	Active    bool   `json:"active"`
}
