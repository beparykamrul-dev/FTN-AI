package main

import (
    "errors"
    "strings"
    "sync"
    "time"
)

// FTNPathOrchestrator connects measured traffic health, service policy,
// core HA eligibility and ECMP planning. It is plan-only: route mutation
// remains behind the existing approval-gated routed engine.
type FTNPathOrchestrator struct {
    mu          sync.Mutex
    controllers map[string]*TrafficPathController
}

type FTNPathOrchestratorInput struct {
    ServiceID       string
    Observations    []TrafficPathObservation
    ECMPCandidates  []FTNRouteCandidate
    CurrentECMP     FTNECMPSelection
    MaxECMPPaths    int
    CoreNodes       []FTNCoreNode
    CurrentCoreNode string
    Now             time.Time
}

type FTNPathOrchestrationPlan struct {
    ServiceID         string
    Core              FTNCoreFailoverDecision
    Traffic           TrafficDecision
    ECMP              FTNECMPSelection
    Reconcile         FTNECMPReconcilePlan
    HasTraffic        bool
    HasECMP           bool
    RequiresApproval  bool
    Reason            string
}

func NewFTNPathOrchestrator() *FTNPathOrchestrator {
    return &FTNPathOrchestrator{controllers: make(map[string]*TrafficPathController)}
}

func (o *FTNPathOrchestrator) controller(serviceID string) *TrafficPathController {
    o.mu.Lock()
    defer o.mu.Unlock()
    c := o.controllers[serviceID]
    if c == nil {
        c = &TrafficPathController{}
        o.controllers[serviceID] = c
    }
    return c
}

func findTrafficPolicy(serviceID string) (TrafficServicePolicy, error) {
    serviceID = strings.TrimSpace(serviceID)
    for _, p := range DefaultTrafficServicePolicies() {
        if p.ID == serviceID {
            return p, nil
        }
    }
    return TrafficServicePolicy{}, errors.New("traffic_service_policy_not_found")
}

// Plan never mutates routing state. It produces the single plan consumed by
// the approval/execution layer.
func (o *FTNPathOrchestrator) Plan(in FTNPathOrchestratorInput) (FTNPathOrchestrationPlan, error) {
    if o == nil {
        return FTNPathOrchestrationPlan{}, errors.New("path_orchestrator_required")
    }
    service, err := findTrafficPolicy(in.ServiceID)
    if err != nil {
        return FTNPathOrchestrationPlan{}, err
    }
    now := in.Now
    if now.IsZero() {
        now = time.Now()
    }

    core := EvaluateFTNCoreFailover(in.CoreNodes, in.CurrentCoreNode)
    if core.ActiveNode == "" {
        return FTNPathOrchestrationPlan{ServiceID: service.ID, Core: core, Reason: "no_healthy_core_available"}, errors.New("no_healthy_core_available")
    }

    traffic, ok := o.controller(service.ID).Decide(in.Observations, service, now)
    if !ok {
        return FTNPathOrchestrationPlan{ServiceID: service.ID, Core: core, Reason: "no_usable_traffic_path"}, errors.New("no_usable_traffic_path")
    }

    selected := make([]FTNRouteCandidate, 0, len(in.ECMPCandidates))
    for _, c := range in.ECMPCandidates {
        if c.PathID == traffic.PathID {
            selected = append(selected, c)
        }
    }
    if len(selected) == 0 {
        return FTNPathOrchestrationPlan{ServiceID: service.ID, Core: core, Traffic: traffic, HasTraffic: true, Reason: "selected_path_has_no_route_candidate"}, errors.New("selected_path_has_no_route_candidate")
    }

    selection, err := SelectFTNECMPPaths(in.ECMPCandidates, in.MaxECMPPaths)
    if err != nil {
        return FTNPathOrchestrationPlan{ServiceID: service.ID, Core: core, Traffic: traffic, HasTraffic: true, Reason: "ecmp_selection_failed"}, err
    }
    plan, err := BuildFTNECMPReconcilePlan(in.CurrentECMP, selection)
    if err != nil {
        return FTNPathOrchestrationPlan{ServiceID: service.ID, Core: core, Traffic: traffic, ECMP: selection, HasTraffic: true, HasECMP: true, Reason: "ecmp_reconcile_plan_failed"}, err
    }

    return FTNPathOrchestrationPlan{
        ServiceID: service.ID,
        Core: core,
        Traffic: traffic,
        ECMP: selection,
        Reconcile: plan,
        HasTraffic: true,
        HasECMP: true,
        RequiresApproval: plan.RequiresApproval,
        Reason: "health_to_exrouter_to_ecmp_plan_ready",
    }, nil
}
