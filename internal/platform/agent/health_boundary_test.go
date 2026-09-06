package agent

import (
	"context"
	"errors"
	"testing"
)

type failingHealthChecker struct{}

func (failingHealthChecker) Check(context.Context) error { return errors.New("unhealthy") }

func TestCheckLayerPropagatesHealthFailure(t *testing.T) {
	got := CheckLayer(context.Background(), Layer{ID: "layer-1"}, failingHealthChecker{})
	if got.Healthy || got.Reason != "unhealthy" {
		t.Fatalf("unexpected health result: %+v", got)
	}
}
