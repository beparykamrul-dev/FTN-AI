package kernel

import (
	"context"
	"errors"
	"fmt"
)

// Request is a structured kernel operation. Raw shell execution is deliberately
// not part of this interface; operations must resolve to a registered tool.
type Request struct {
	Tool       string
	Operation  string
	Target     string
	Parameters map[string]string
}

type Result struct {
	Status string
	Output string
}

// Tool is the narrow execution boundary used by FTN's notebook/network kernel.
type Tool interface {
	Name() string
	Operations() []string
	Execute(ctx context.Context, request Request) (Result, error)
}

type Registry struct {
	tools map[string]Tool
}

func NewRegistry(tools ...Tool) (*Registry, error) {
	r := &Registry{tools: make(map[string]Tool)}
	for _, tool := range tools {
		if tool == nil || tool.Name() == "" {
			return nil, errors.New("kernel tool must have a name")
		}
		if _, exists := r.tools[tool.Name()]; exists {
			return nil, fmt.Errorf("duplicate kernel tool: %s", tool.Name())
		}
		r.tools[tool.Name()] = tool
	}
	return r, nil
}

func (r *Registry) Execute(ctx context.Context, request Request) (Result, error) {
	if request.Tool == "" || request.Operation == "" {
		return Result{}, errors.New("tool and operation are required")
	}
	tool, ok := r.tools[request.Tool]
	if !ok {
		return Result{}, fmt.Errorf("kernel tool not registered: %s", request.Tool)
	}
	return tool.Execute(ctx, request)
}

// DoExecute is the controlled equivalent of a kernel do_execute operation.
// It accepts structured tool calls only and never forwards arbitrary code or
// shell text to a host.
func (r *Registry) DoExecute(ctx context.Context, request Request) (Result, error) {
	return r.Execute(ctx, request)
}
