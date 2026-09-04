package kernel

import (
	"context"
	"errors"
	"fmt"
)

// ToolRequest is the structured execution request used by registered kernel
// tools. It is deliberately distinct from the backend Request contract.
type ToolRequest struct {
	Tool       string
	Operation  string
	Target     string
	Parameters map[string]string
}

type ToolResult struct {
	Status string
	Output string
}

// Tool is the narrow execution boundary used by FTN's notebook/network kernel.
type Tool interface {
	Name() string
	Operations() []string
	Execute(ctx context.Context, request ToolRequest) (ToolResult, error)
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

func (r *Registry) Execute(ctx context.Context, request ToolRequest) (ToolResult, error) {
	if request.Tool == "" || request.Operation == "" {
		return ToolResult{}, errors.New("tool and operation are required")
	}
	tool, ok := r.tools[request.Tool]
	if !ok {
		return ToolResult{}, fmt.Errorf("kernel tool not registered: %s", request.Tool)
	}
	return tool.Execute(ctx, request)
}

// DoExecute is the controlled equivalent of a kernel do_execute operation.
// It accepts structured tool calls only and never forwards arbitrary code or
// shell text to a host.
func (r *Registry) DoExecute(ctx context.Context, request ToolRequest) (ToolResult, error) {
	return r.Execute(ctx, request)
}
