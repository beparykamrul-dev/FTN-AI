package agent

import (
	"context"
	"fmt"
)

// Category is the bounded responsibility of a specialized FTN AI.
type Category string

const (
	CategoryStudio      Category = "studio"
	CategoryCallCenter  Category = "call-center"
	CategoryBilling     Category = "billing"
	CategoryNetwork     Category = "network"
	CategoryDeveloper   Category = "developer"
	CategoryCustomer    Category = "customer"
	CategoryExecutive   Category = "executive"
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
	fleet *Fleet
	runtimes map[Category]CategoryRuntime
}

func NewOrchestrator(fleet *Fleet, runtimes map[Category]CategoryRuntime) (*Orchestrator, error) {
	if fleet == nil { return nil, fmt.Errorf("agent fleet is required") }
	if len(runtimes) == 0 { return nil, fmt.Errorf("agent categories are required") }
	return &Orchestrator{fleet: fleet, runtimes: runtimes}, nil
}

func (o *Orchestrator) Handle(ctx context.Context, request CategoryRequest) (Response, error) {
	runtime, ok := o.runtimes[request.Category]
	if !ok { return Response{}, fmt.Errorf("unsupported agent category: %s", request.Category) }
	return runtime.Run(ctx, request)
}
