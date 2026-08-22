package agent

import "context"

// Decision describes why a layer was selected and whether escalation is allowed.
type Decision struct {
	LayerID     string
	Capability  string
	Escalate    bool
	Reason      string
}

// Decide keeps routing explainable: FTN can inspect why a layer was selected.
func Decide(ctx context.Context, registry *LayerRegistry, category Category, capability string, allowEscalation bool) (Decision, error) {
	if err := ctx.Err(); err != nil { return Decision{}, err }
	layer, err := registry.Resolve(category)
	if err != nil { return Decision{}, err }
	return Decision{LayerID: layer.ID, Capability: capability, Escalate: allowEscalation, Reason: "best enabled layer for category"}, nil
}
