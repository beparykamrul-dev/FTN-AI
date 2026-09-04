package kernel

import "context"

// Request is the provider-neutral boundary between a kernel client and the
// FTN backend. Payloads are structured data; they are never shell commands.
type Request struct {
	ID           string
	ServerID     string
	Tool         string
	Capability   string
	Arguments    map[string]string
	Mutating     bool
	ApprovalID   string
	Idempotency  string
}

type Result struct {
	RequestID string
	Status    string
	Output    map[string]string
	Observed  bool
	Verified  bool
}

// Backend is the minimal contract implemented by the FTN control plane.
type Backend interface {
	Execute(context.Context, Request) (Result, error)
}

// Policy is intentionally evaluated before a backend adapter executes.
type Policy interface {
	Authorize(context.Context, Request) error
}
