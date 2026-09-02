package main

import (
	"context"
	"errors"
	"strings"
	"time"
)

type CoreRouterNode struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	Address  string `json:"address,omitempty"`
	Protocol string `json:"protocol"`
	Enabled  bool   `json:"enabled"`
}

type CoreRouterPeer struct {
	ID          string `json:"id"`
	LocalNode   string `json:"local_node"`
	RemoteASN   uint32 `json:"remote_asn"`
	RemoteIP    string `json:"remote_ip"`
	AddressFamily string `json:"address_family"`
	Enabled     bool   `json:"enabled"`
}

type CoreRouterAdapter interface {
	Protocol() string
	Capabilities(context.Context, CoreRouterNode) ([]string, error)
	Health(context.Context, CoreRouterNode) (CoreRouterHealth, error)
	Peers(context.Context, CoreRouterNode) ([]CoreRouterPeerState, error)
	PlanRouteChange(context.Context, CoreRouterNode, RouteChangeIntent) (RouteChangePlan, error)
}

type CoreRouterHealth struct {
	NodeID     string    `json:"node_id"`
	Healthy    bool      `json:"healthy"`
	CheckedAt  time.Time `json:"checked_at"`
	LatencyMS  float64   `json:"latency_ms"`
	Reason     string    `json:"reason"`
}

type CoreRouterPeerState struct {
	PeerID       string `json:"peer_id"`
	Established  bool   `json:"established"`
	Prefixes     uint64 `json:"prefixes"`
	AddressFamily string `json:"address_family"`
}

type RouteChangeIntent struct {
	Action      string `json:"action"`
	Prefix      string `json:"prefix"`
	NextHop     string `json:"next_hop,omitempty"`
	PeerID      string `json:"peer_id,omitempty"`
	ApprovalID  string `json:"approval_id,omitempty"`
	ChangeID    string `json:"change_id,omitempty"`
}

type RouteChangePlan struct {
	Allowed             bool     `json:"allowed"`
	RequiresApproval    bool     `json:"requires_approval"`
	PreChangeSnapshot   bool     `json:"pre_change_snapshot"`
	PostChangeVerify    bool     `json:"post_change_verify"`
	RollbackWhenSafe    bool     `json:"rollback_when_safe"`
	Risk                string   `json:"risk"`
	ValidationErrors    []string `json:"validation_errors,omitempty"`
}

func ValidateCoreRouterNode(n CoreRouterNode) error {
	if strings.TrimSpace(n.ID) == "" { return errors.New("core_node_id_required") }
	if n.Role != "primary" && n.Role != "standby" { return errors.New("invalid_core_node_role") }
	if strings.TrimSpace(n.Protocol) == "" { return errors.New("core_router_protocol_required") }
	return nil
}

func PlanCoreRouteChange(n CoreRouterNode, i RouteChangeIntent) RouteChangePlan {
	p := RouteChangePlan{RequiresApproval:true, PreChangeSnapshot:true, PostChangeVerify:true, RollbackWhenSafe:true, Risk:"medium"}
	if err := ValidateCoreRouterNode(n); err != nil { p.Allowed=false; p.ValidationErrors=[]string{err.Error()}; return p }
	if strings.TrimSpace(i.Action)=="" || strings.TrimSpace(i.Prefix)=="" { p.Allowed=false; p.ValidationErrors=[]string{"route_action_and_prefix_required"}; return p }
	if strings.TrimSpace(i.ApprovalID)=="" { p.Allowed=false; p.ValidationErrors=[]string{"approval_required"}; return p }
	if i.Action=="withdraw" { p.Risk="high" }
	if i.Action!="announce" && i.Action!="withdraw" && i.Action!="attribute-change" { p.Allowed=false; p.ValidationErrors=[]string{"unsupported_route_action"}; return p }
	p.Allowed=true
	return p
}
