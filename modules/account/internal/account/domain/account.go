package domain

import "time"

type AccountType string
const ( AccountAsset AccountType = "asset"; AccountLiability AccountType = "liability"; AccountEquity AccountType = "equity"; AccountIncome AccountType = "income"; AccountExpense AccountType = "expense" )

type Account struct { ID string `json:"id"`; Name string `json:"name"`; Type AccountType `json:"type"`; Currency string `json:"currency"`; Active bool `json:"active"`; CreatedAt time.Time `json:"created_at"`; UpdatedAt time.Time `json:"updated_at"` }

type EntrySide string
const ( Debit EntrySide = "debit"; Credit EntrySide = "credit" )

type LedgerEntry struct { ID string `json:"id"`; JournalID string `json:"journal_id"`; AccountID string `json:"account_id"`; Side EntrySide `json:"side"`; AmountMinor int64 `json:"amount_minor"`; Currency string `json:"currency"`; SourceType string `json:"source_type"`; SourceID string `json:"source_id"`; CreatedAt time.Time `json:"created_at"` }
