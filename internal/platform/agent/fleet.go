package agent

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Scope struct {
	ServiceID string
	UserID    string
	Role      string
	TenantID  string
}

type Mode string

const (
	ModeService   Mode = "service"
	ModeUser      Mode = "user"
	ModeDeveloper Mode = "developer"
	ModeAssistant Mode = "assistant"
)

type Request struct {
	Scope Scope
	Mode  Mode
	Input string
}

type Response struct {
	Text          string
	NeedsApproval bool
	Action        string
}

type Runtime interface {
	Run(ctx context.Context, request Request) (Response, error)
}

type Policy interface {
	Allow(ctx context.Context, request Request, response Response) error
}

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
	if f == nil {
		return Response{}, errors.New("agent fleet is required")
	}
	if ctx == nil {
		return Response{}, errors.New("context is required")
	}
	request.Scope.ServiceID = strings.TrimSpace(request.Scope.ServiceID)
	request.Scope.UserID = strings.TrimSpace(request.Scope.UserID)
	request.Scope.Role = strings.TrimSpace(request.Scope.Role)
	request.Scope.TenantID = strings.TrimSpace(request.Scope.TenantID)
	if request.Scope.ServiceID == "" && request.Scope.UserID == "" {
		return Response{}, errors.New("service_id or user_id is required")
	}
	request.Input = strings.TrimSpace(request.Input)
	if request.Input == "" {
		return Response{}, errors.New("input is required")
	}
	if request.Mode == "" {
		request.Mode = ModeAssistant
	}
	f.mu.RLock()
	runtime := f.runtime
	policy := f.policy
	f.mu.RUnlock()
	if runtime == nil || policy == nil {
		return Response{}, errors.New("agent fleet dependencies unavailable")
	}
	response, err := runtime.Run(ctx, request)
	if err != nil {
		return Response{}, fmt.Errorf("agent runtime: %w", err)
	}
	if err := policy.Allow(ctx, request, response); err != nil {
		return Response{}, fmt.Errorf("agent policy: %w", err)
	}
	return response, nil
}
