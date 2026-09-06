package agent

import (
	"context"
	"testing"
)

func TestDecideRejectsMissingCapability(t *testing.T) {
	r := NewLayerRegistry([]Layer{{ID: "native", Categories: []Category{"chat"}, Enabled: true, Native: true}})
	if _, err := Decide(context.Background(), r, Category("chat"), "", true); err == nil {
		t.Fatal("expected capability validation error")
	}
}
