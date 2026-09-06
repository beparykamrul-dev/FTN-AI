package agent

import "testing"

func TestPlansHavePositiveLimits(t *testing.T) {
	for id, p := range Plans {
		if p.ID != id || p.RequestsPerDay <= 0 || p.TokensPerDay <= 0 {
			t.Fatalf("invalid plan %q: %+v", id, p)
		}
	}
}
