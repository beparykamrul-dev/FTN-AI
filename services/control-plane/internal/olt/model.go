package olt

type OLT struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Vendor   string `json:"vendor"`
	Address  string `json:"address"`
	Status   string `json:"status"`
	PONCount int    `json:"pon_count"`
}
