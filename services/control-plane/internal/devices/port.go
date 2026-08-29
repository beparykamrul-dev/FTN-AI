package devices

type Port struct {
	ID       string `json:"id"`
	DeviceID string `json:"device_id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	RXBps    uint64 `json:"rx_bps"`
	TXBps    uint64 `json:"tx_bps"`
}
