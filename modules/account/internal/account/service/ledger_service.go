package service

import (
 "context"
 "errors"
 "github.com/beparykamrul-dev/FTN-AI/internal/account/domain"
)

type EntryRepository interface { InsertEntry(context.Context, domain.LedgerEntry) error }

type LedgerService struct { repo EntryRepository }
func NewLedgerService(repo EntryRepository) *LedgerService { return &LedgerService{repo: repo} }

func (s *LedgerService) Validate(entries []domain.LedgerEntry) error {
 if len(entries) < 2 { return errors.New("journal requires at least two entries") }
 var debit, credit int64
 for _, e := range entries {
  if e.AmountMinor <= 0 { return errors.New("ledger amount must be positive") }
  switch e.Side { case domain.Debit: debit += e.AmountMinor; case domain.Credit: credit += e.AmountMinor; default: return errors.New("invalid ledger side") }
 }
 if debit != credit { return errors.New("unbalanced ledger journal") }
 return nil
}

func (s *LedgerService) Post(ctx context.Context, entries []domain.LedgerEntry) error {
 if err := s.Validate(entries); err != nil { return err }
 for _, e := range entries { if err := s.repo.InsertEntry(ctx, e); err != nil { return err } }
 return nil
}
