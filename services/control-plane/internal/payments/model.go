package payments

import "time"

type Payment struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"account_id"`
	InvoiceID     string    `json:"invoice_id"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	Provider      string    `json:"provider"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
}
