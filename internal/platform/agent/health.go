package agent

import (
	"context"
	"fmt"
)

type LayerHealth struct {
	LayerID string
	Healthy bool
	Reason  string
}

type HealthChecker interface {
	Check(context.Context) error
}

func CheckLayer(ctx context.Context, layer Layer, checker HealthChecker) LayerHealth {
	if ctx == nil {
		return LayerHealth{LayerID: layer.ID, Reason: "context is required"}
	}
	if err := ctx.Err(); err != nil {
		return LayerHealth{LayerID: layer.ID, Reason: err.Error()}
	}
	if checker == nil {
		return LayerHealth{LayerID: layer.ID, Healthy: false, Reason: "health checker unavailable"}
	}
	if err := checker.Check(ctx); err != nil {
		return LayerHealth{LayerID: layer.ID, Healthy: false, Reason: err.Error()}
	}
	return LayerHealth{LayerID: layer.ID, Healthy: true, Reason: fmt.Sprintf("layer %s healthy", layer.ID)}
}
