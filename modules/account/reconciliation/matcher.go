package reconciliation

import "strings"

type Candidate struct {
	ID         string
	Reference  string
	AmountMinor int64
	Party      string
}

type Event struct {
	Reference   string
	AmountMinor int64
	Party       string
}

// ExactMatch intentionally uses deterministic evidence only.
func ExactMatch(event Event, candidates []Candidate) *Candidate {
	ref := strings.TrimSpace(event.Reference)
	party := strings.TrimSpace(event.Party)
	for i := range candidates {
		c := &candidates[i]
		if ref != "" && c.Reference == ref && c.AmountMinor == event.AmountMinor {
			return c
		}
		if ref == "" && event.AmountMinor > 0 && c.AmountMinor == event.AmountMinor && party != "" && strings.EqualFold(c.Party, party) {
			return c
		}
	}
	return nil
}
