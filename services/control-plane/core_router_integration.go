package main

import (
	"context"
	"errors"
)

// CoreRouterIntegration is the control-plane boundary for a concrete router
// adapter. The adapter owns transport/authentication; the control plane owns
// policy, approval and lifecycle state.
type CoreRouterIntegration struct {
	Adapter CoreRouterAdapter
}

func (i CoreRouterIntegration) Inspect(ctx context.Context, node CoreRouterNode) (CoreRouterHealth, []CoreRouterPeerState, error) {
	if i.Adapter == nil { return CoreRouterHealth{}, nil, errors.New("core_router_adapter_required") }
	if err := ValidateCoreRouterNode(node); err != nil { return CoreRouterHealth{}, nil, err }
	h, err := i.Adapter.Health(ctx, node)
	if err != nil { return CoreRouterHealth{}, nil, err }
	peers, err := i.Adapter.Peers(ctx, node)
	if err != nil { return CoreRouterHealth{}, nil, err }
	return h, peers, nil
}

func (i CoreRouterIntegration) Plan(ctx context.Context, node CoreRouterNode, intent RouteChangeIntent) (RouteChangePlan, error) {
	if i.Adapter == nil { return RouteChangePlan{}, errors.New("core_router_adapter_required") }
	if err := ValidateCoreRouterNode(node); err != nil { return RouteChangePlan{}, err }
	return i.Adapter.PlanRouteChange(ctx, node, intent)
}

// Execute remains deliberately separate from Plan. A concrete adapter may
// expose execution through the durable-job worker, but this boundary refuses
// to turn an unapproved plan into a mutation.
func RequireApprovedRoutePlan(plan RouteChangePlan, approvalID string) error {
	if !plan.Allowed { return errors.New("route_plan_not_allowed") }
	if !plan.RequiresApproval || approvalID == "" { return errors.New("approval_required") }
	if !plan.PreChangeSnapshot || !plan.PostChangeVerify || !plan.RollbackWhenSafe {
		return errors.New("route_safety_gates_incomplete")
	}
	return nil
}
