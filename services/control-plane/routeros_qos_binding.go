package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type RouterOSPathBinding struct {
	PathID string `json:"path_id"`
	RouteMark string `json:"route_mark"`
	QueueClass string `json:"queue_class"`
}

type RouterOSPathRegistry struct { Bindings []RouterOSPathBinding `json:"bindings"` }

func (r RouterOSPathRegistry) Resolve(pathID string) (RouterOSPathBinding, error) {
	pathID = strings.TrimSpace(pathID)
	if pathID == "" { return RouterOSPathBinding{}, errors.New("path_id_required") }
	for _, b := range r.Bindings {
		if b.PathID == pathID && strings.TrimSpace(b.RouteMark) != "" && strings.TrimSpace(b.QueueClass) != "" { return b, nil }
	}
	return RouterOSPathBinding{}, fmt.Errorf("unregistered_path_id: %s", pathID)
}

func BuildRouterOSQOSPlanWithRegistry(device NetworkDevice, decisions []TrafficDecision, registry RouterOSPathRegistry) (RouterOSTrafficQoSPlan, error) {
	plan, err := BuildRouterOSTrafficQoSPlan(device, decisions)
	if err != nil { return RouterOSTrafficQoSPlan{}, err }
	for _, rule := range plan.Rules { if _, err := registry.Resolve(rule.PathID); err != nil { return RouterOSTrafficQoSPlan{}, err } }
	sort.Slice(plan.Rules, func(i, j int) bool { if plan.Rules[i].ServiceID == plan.Rules[j].ServiceID { return plan.Rules[i].PathID < plan.Rules[j].PathID }; return plan.Rules[i].ServiceID < plan.Rules[j].ServiceID })
	return plan, nil
}
