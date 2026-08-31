package dns

type Record struct {
	ID       string `json:"id"`
	Zone     string `json:"zone"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	TTL      uint32 `json:"ttl"`
	Enabled  bool   `json:"enabled"`
}
