package agent

import (
	"context"
	"errors"
	"strings"
)

type RouteRequest struct { Scope Scope; Category Category; Input string }
type LayerRuntime interface { Run(ctx context.Context, request RouteRequest) (Response, error) }
type Router struct { registry *LayerRegistry; runtimes map[string]LayerRuntime }
func NewRouter(registry *LayerRegistry, runtimes map[string]LayerRuntime) *Router {
	copyRuntimes := make(map[string]LayerRuntime, len(runtimes))
	for id, runtime := range runtimes { id = strings.TrimSpace(id); if id != "" && runtime != nil { copyRuntimes[id] = runtime } }
	return &Router{registry: registry, runtimes: copyRuntimes}
}
func (r *Router) Handle(ctx context.Context, request RouteRequest) (Response, error) {
	if r == nil || r.registry == nil { return Response{}, errors.New("AI layer registry unavailable") }
	if ctx == nil { return Response{}, errors.New("context is required") }
	select { case <-ctx.Done(): return Response{}, ctx.Err(); default: }
	request.Input = strings.TrimSpace(request.Input)
	if request.Input == "" { return Response{}, errors.New("input is required") }
	layer, err := r.registry.Resolve(request.Category)
	if err != nil { return Response{}, err }
	if layer == nil { return Response{}, errors.New("AI layer unavailable") }
	runtime, ok := r.runtimes[strings.TrimSpace(layer.ID)]
	if !ok || runtime == nil { return Response{}, &layerUnavailableError{ID: strings.TrimSpace(layer.ID)} }
	return runtime.Run(ctx, request)
}
type layerUnavailableError struct{ ID string }
func (e *layerUnavailableError) Error() string { return "AI layer runtime unavailable: " + e.ID }
