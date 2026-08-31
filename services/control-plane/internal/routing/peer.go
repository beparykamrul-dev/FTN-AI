package routing

type Peer struct {
	ID        string `json:"id"`
	Protocol  string `json:"protocol"`
	Address   string `json:"address"`
	ASN       uint32 `json:"asn"`
	State     string `json:"state"`
}
