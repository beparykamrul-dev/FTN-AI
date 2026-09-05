package agent

import (
	"context"
	"fmt"
)

type CategoryRequest struct {
	Scope    Scope
	Category Category
	Input    string
}

type CategoryRuntime interface {
	Run(ctx context.Context, request CategoryRequest) (Response, error)
}

// Orchestrator routes requests to specialized agents while keeping one shared,
// lightweight runtime and the existing policy/approval boundary.
type Orchestrator struct {
	fleet    *Fleet
	runtimes map[Category]CategoryRuntime
}

func NewOrchestrator(fleet *Fleet, runtimes map[Category]CategoryRuntime) (*Orchestrator, error) {
	if fleet == nil {
		return nil, fmt.Errorf("agent fleet is required")
	}
	if len(runtimes) == 0 {
		return nil, fmt.Errorf("agent categories are required")
	}
	return &Orchestrator{fleet: fleet, runtimes: runtimes}, nil
}

func (o *Orchestrator) Handle(ctx context.Context, request CategoryRequest) (Response, error) {
	if o == nil {
		return Response{}, fmt.Errorf("agent orchestrator is required")
	}
	if ctx == nil {
		return Response{}, fmt.Errorf("context is required")
	}
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}
	runtime, ok := o.runtimes[request.Category]
	if !ok || runtime == nil {
		return Response{}, fmt.Errorf("unsupported agent category: %s", request.Category)
	}
	return runtime.Run(ctx, request)
}
