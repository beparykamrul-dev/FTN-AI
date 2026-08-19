package agent

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Scope identifies who or what an FTN agent is serving.
type Scope struct {
	ServiceID string
	UserID    string
	Role      string
	TenantID  string
}

// Mode identifies the agent responsibility without creating separate runtimes.
type Mode string

const (
	ModeService    Mode = "service"
	ModeUser       Mode = "user"
	ModeDeveloper  Mode = "developer"
	ModeAssistant  Mode = "assistant"
)

// Request is the normalized input passed to an agent.
type Request struct {
	Scope Scope
	Mode  Mode
	Input string
}

// Response is deliberately provider-neutral so local models can be used without
// changing the service/user/developer agent layer.
type Response struct {
	Text       string
	NeedsApproval bool
	Action     string
}

// Runtime is the local/private intelligence boundary. Implementations may use
// an FTN-hosted model, rules engine, retrieval system, or another approved backend.
type Runtime interface {
	Run(ctx context.Context, request Request) (Response, error)
}

// Policy prevents a conversational agent from silently executing side effects.
type Policy interface {
	Allow(ctx context.Context, request Request, response Response) error
}

// Fleet provides one logical agent identity per service/user scope while sharing
// the same runtime, policy, memory, and tool infrastructure.
type Fleet struct {
	runtime Runtime
	policy  Policy
	mu      sync.RWMutex
}

func NewFleet(runtime Runtime, policy Policy) (*Fleet, error) {
	if runtime == nil {
		return nil, errors.New("agent runtime is required")
	}
	if policy == nil {
		return nil, errors.New("agent policy is required")
	}
	return &Fleet{runtime: runtime, policy: policy}, nil
}

func (f *Fleet) Handle(ctx context.Context, request Request) (Response, error) {
	if request.Scope.ServiceID == "" && request.Scope.UserID == "" {
		return Response{}, errors.New("service_id or user_id is required")
	}
	if request.Mode == "" {
		request.Mode = ModeAssistant
	}
	f.mu.RLock()
	runtime := f.runtime
	policy := f.policy
	f.mu.RUnlock()

	response, err := runtime.Run(ctx, request)
	if err != nil {
		return Response{}, fmt.Errorf("agent runtime: %w", err)
	}
	if err := policy.Allow(ctx, request, response); err != nil {
		return Response{}, fmt.Errorf("agent policy: %w", err)
	}
	return response, nil
}
