package main

import (
	"errors"
	"sort"
	"strings"
)

type FTNCoreRouterHAState struct {
	ActiveNode string `json:"active_node"`
	StandbyNode string `json:"standby_node,omitempty"`
	Failover bool `json:"failover"`
	Reason string `json:"reason"`
}

// SelectFTNCoreRouterHA chooses only an enabled, healthy, BGP-ready node whose
// BFD state is not down. Existing-node preference prevents unnecessary churn.
func SelectFTNCoreRouterHA(nodes []FTNCoreNode, current string) (FTNCoreRouterHAState, error) {
	if len(nodes) == 0 { return FTNCoreRouterHAState{}, errors.New("core_router_nodes_required") }
	ordered := append([]FTNCoreNode(nil), nodes...)
	sort.SliceStable(ordered, func(i,j int) bool {
		if ordered[i].ID == current { return true }
		if ordered[j].ID == current { return false }
		return ordered[i].ID < ordered[j].ID
	})
	for _, n := range ordered {
		if strings.TrimSpace(n.ID) == "" || !n.Healthy || !n.BGPReady || n.BFDState == FTNBFDDown { continue }
		state := FTNCoreRouterHAState{ActiveNode:n.ID, Failover:n.ID != current, Reason:"healthy_core_selected"}
		for _, s := range ordered { if s.ID != n.ID && s.Healthy && s.BGPReady && s.BFDState != FTNBFDDown { state.StandbyNode=s.ID; break } }
		if state.Failover { state.Reason="failover_to_healthy_core" } else { state.Reason="current_core_remains_active" }
		return state,nil
	}
	return FTNCoreRouterHAState{Reason:"no_eligible_core"}, errors.New("no_eligible_core")
}
