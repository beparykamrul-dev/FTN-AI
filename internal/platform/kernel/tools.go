package kernel

import "context"

// NamedTool provides the registration metadata shared by network automation
// integrations. Implementations should delegate to approved libraries/APIs.
type NamedTool struct {
	ToolName string
	Ops      []string
	Run      func(context.Context, Request) (Result, error)
}

func (t NamedTool) Name() string { return t.ToolName }
func (t NamedTool) Operations() []string { return append([]string(nil), t.Ops...) }
func (t NamedTool) Execute(ctx context.Context, request Request) (Result, error) {
	if t.Run == nil {
		return Result{}, ErrToolNotImplemented
	}
	return t.Run(ctx, request)
}

var ErrToolNotImplemented = errNotImplemented{}
type errNotImplemented struct{}
func (errNotImplemented) Error() string { return "kernel tool implementation is not configured" }

// StandardTools describes the optional FTN integration points. The repository
// does not vendor Python runtimes or credentials; adapters can be enabled by
// the deployment that owns those dependencies.
func StandardTools() []NamedTool {
	return []NamedTool{
		{ToolName: "cert-netsa", Ops: []string{"query", "inspect"}},
		{ToolName: "netmiko", Ops: []string{"connect", "collect"}},
		{ToolName: "paramiko", Ops: []string{"connect", "collect"}},
		{ToolName: "napalm", Ops: []string{"get", "collect"}},
		{ToolName: "ipykernel", Ops: []string{"kernel-info", "execute-tool"}},
	}
}
