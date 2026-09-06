package agent

import (
	"context"
	"testing"
)

type quotaTestStore struct{ usage Usage }

func (s *quotaTestStore) Get(context.Context, Scope) (Usage, error) { return s.usage, nil }
func (s *quotaTestStore) Put(_ context.Context, _ Scope, u Usage) error { s.usage = u; return nil }

func TestQuotaGateRejectsNilStore(t *testing.T) {
	if err := (&QuotaGate{}).CheckAndConsume(context.Background(), Scope{}, Plans["free"], 1); err == nil {
		t.Fatal("expected missing store error")
	}
}
