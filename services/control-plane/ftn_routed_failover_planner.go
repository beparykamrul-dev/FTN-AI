package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type FTNRoutedFailoverPlan struct {
	CurrentNode string `json:"current_node"`
	TargetNode string `json:"target_node"`
	RequiresApproval bool `json:"requires_approval"`
	Allowed bool `json:"allowed"`
	Risk string `json:"risk"`
	Reason string `json:"reason"`
	DecisionHash string `json:"decision_hash"`
}

// PlanFTNRoutedFailover creates a deterministic, approval-gated failover plan.
// Planning is side-effect free: it does not alter BGP, BFD, RIB, FIB, or peers.
func PlanFTNRoutedFailover(nodes []FTNCoreNode, current string) FTNRoutedFailoverPlan {
	ordered := append([]FTNCoreNode(nil), nodes...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })
	decision := EvaluateFTNCoreFailover(ordered, current)
	plan := FTNRoutedFailoverPlan{
		CurrentNode: strings.TrimSpace(current),
		TargetNode: decision.ActiveNode,
		RequiresApproval: decision.Failover,
		Allowed: decision.ActiveNode != "",
		Risk: "low",
		Reason: decision.Reason,
	}
	if decision.Failover {
		plan.Risk = "high"
	}
	plan.DecisionHash = hashFailoverPlan(plan)
	return plan
}

func hashFailoverPlan(plan FTNRoutedFailoverPlan) string {
	payload := fmt.Sprintf("%s|%s|%t|%t|%s|%s", plan.CurrentNode, plan.TargetNode, plan.RequiresApproval, plan.Allowed, plan.Risk, plan.Reason)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
