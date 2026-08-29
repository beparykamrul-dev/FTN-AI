package billing

import "time"

type Invoice struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	DueAt     time.Time `json:"due_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
