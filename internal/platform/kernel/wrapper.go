package kernel

import (
	"context"
	"errors"
	"fmt"
)

// Tool is the narrow execution boundary used by FTN's notebook/network kernel.
type Tool interface {
	Name() string
	Operations() []string
	Execute(ctx context.Context, request Request) (Result, error)
}

type Registry struct { tools map[string]Tool }

func NewRegistry(tools ...Tool) (*Registry, error) {
	r := &Registry{tools: make(map[string]Tool)}
	for _, tool := range tools {
		if tool == nil || tool.Name() == "" { return nil, errors.New("kernel tool must have a name") }
		if _, exists := r.tools[tool.Name()]; exists { return nil, fmt.Errorf("duplicate kernel tool: %s", tool.Name()) }
		r.tools[tool.Name()] = tool
	}
	return r, nil
}

func (r *Registry) Execute(ctx context.Context, request Request) (Result, error) {
	if request.Tool == "" || request.Operation == "" { return Result{}, errors.New("tool and operation are required") }
	tool, ok := r.tools[request.Tool]
	if !ok { return Result{}, fmt.Errorf("kernel tool not registered: %s", request.Tool) }
	return tool.Execute(ctx, request)
}

func (r *Registry) DoExecute(ctx context.Context, request Request) (Result, error) { return r.Execute(ctx, request) }
