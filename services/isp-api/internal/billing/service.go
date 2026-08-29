package billing

import "context"

type Invoice struct {
	ID         string `json:"id"`
	AccountID  string `json:"account_id"`
	AmountMinor int64  `json:"amount_minor"`
	Currency   string `json:"currency"`
	Status     string `json:"status"`
}

type Repository interface {
	GetInvoice(ctx context.Context, id string) (Invoice, error)
	ListAccountInvoices(ctx context.Context, accountID string) ([]Invoice, error)
}

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) GetInvoice(ctx context.Context, id string) (Invoice, error) {
	return s.repo.GetInvoice(ctx, id)
}

func (s *Service) ListAccountInvoices(ctx context.Context, accountID string) ([]Invoice, error) {
	return s.repo.ListAccountInvoices(ctx, accountID)
}
