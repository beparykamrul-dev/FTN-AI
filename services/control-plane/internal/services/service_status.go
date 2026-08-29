package services

type ServiceStatus struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}
