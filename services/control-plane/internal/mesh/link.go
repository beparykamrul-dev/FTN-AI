package mesh

type Link struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Status  string `json:"status"`
	Latency int64  `json:"latency_ms"`
}
