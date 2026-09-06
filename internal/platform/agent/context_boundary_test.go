package agent

import (
	"context"
	"testing"
)

func TestCheckLayerHonorsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got := CheckLayer(ctx, Layer{ID: "layer-1"}, nil)
	if got.Reason == "" {
		t.Fatal("expected cancelled-context reason")
	}
}
