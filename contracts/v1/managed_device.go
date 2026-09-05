package contracts

// ManagedDevice is the common identity/health contract for routers and cameras.
type ManagedDevice struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	Vendor       string `json:"vendor"`
	Model        string `json:"model,omitempty"`
	MAC          string `json:"mac"`
	SerialNumber string `json:"serial_number"`
	SiteID       string `json:"site_id"`
	GatewayID    string `json:"gateway_id,omitempty"`
	PowerOnly    bool   `json:"power_only"`
	Status       string `json:"status"`
}

type DeviceHealth struct {
	Reachable          bool    `json:"reachable"`
	UptimeSeconds      uint64  `json:"uptime_seconds"`
	CPUPercent         float64 `json:"cpu_percent"`
	MemoryPercent      float64 `json:"memory_percent"`
	TemperatureC       float64 `json:"temperature_c,omitempty"`
	PacketLossPercent  float64 `json:"packet_loss_percent"`
	LatencyMs          float64 `json:"latency_ms"`
	StoragePercent     float64 `json:"storage_percent,omitempty"`
	LastSeen           string  `json:"last_seen"`
}

type DeviceMetrics struct {
	RxBps    uint64 `json:"rx_bps"`
	TxBps    uint64 `json:"tx_bps"`
	Errors   uint64 `json:"errors"`
	Drops    uint64 `json:"drops"`
	Sessions uint64 `json:"sessions"`
	Events   uint64 `json:"events"`
}
