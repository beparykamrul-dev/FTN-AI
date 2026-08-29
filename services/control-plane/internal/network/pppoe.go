package network

type PPPoESession struct {
	ID       string `json:"id"`
	AccountID string `json:"account_id"`
	DeviceID string `json:"device_id"`
	Username string `json:"username"`
	Address  string `json:"address"`
	Status   string `json:"status"`
}
