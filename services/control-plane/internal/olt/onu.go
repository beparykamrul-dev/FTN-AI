package olt

type ONU struct {
	ID        string `json:"id"`
	OLTID     string `json:"olt_id"`
	PON       string `json:"pon"`
	Serial    string `json:"serial"`
	AccountID string `json:"account_id,omitempty"`
	Status    string `json:"status"`
}
