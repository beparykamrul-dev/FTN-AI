package kernel

import "context"

// Request is the shared structured operation boundary for backend execution
// and registered kernel tools. It never carries arbitrary shell/code text.
type Request struct {
	ID          string
	ServerID    string
	Tool        string
	Operation   string
	Target      string
	Capability  string
	Arguments   map[string]string
	Parameters  map[string]string
	Mutating    bool
	ApprovalID  string
	Idempotency string
}

type Result struct {
	RequestID string
	Status    string
	Output    map[string]string
	Observed  bool
	Verified  bool
}

type Backend interface {
	Execute(context.Context, Request) (Result, error)
}

type Policy interface {
	Authorize(context.Context, Request) error
}
