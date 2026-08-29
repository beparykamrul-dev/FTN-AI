package audit

type Actor struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	IP   string `json:"ip,omitempty"`
}
