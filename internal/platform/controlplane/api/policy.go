package api

import "fmt"

type PolicyDecision string

const (
	DecisionAllow   PolicyDecision = "allow"
	DecisionRequire PolicyDecision = "require_approval"
	DecisionDeny   PolicyDecision = "deny"
)

type Policy struct {
	AllowObserve bool
	AllowPlan bool
	RequireApprovalForExecute bool
}

func (p Policy) Evaluate(c Command) (PolicyDecision, error) {
	if c.ID == "" || c.Target == "" || c.Action == "" { return DecisionDeny, fmt.Errorf("invalid command") }
	switch c.Kind {
	case CommandObserve:
		if p.AllowObserve { return DecisionAllow, nil }
	case CommandPlan:
		return DecisionAllow, nil
	case CommandExecute:
		if p.RequireApprovalForExecute && !c.Approved { return DecisionRequire, nil }
		if c.Approved { return DecisionAllow, nil }
	}
	return DecisionDeny, fmt.Errorf("command denied by policy")
}
