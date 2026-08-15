package agent

import "context"

// RouteRequest is the normalized decision input for FTN's AI layer router.
type RouteRequest struct {
	Scope    Scope
	Category Category
	Input    string
}

// LayerRuntime executes one registered AI layer.
type LayerRuntime interface {
	Run(ctx context.Context, request RouteRequest) (Response, error)
}

// Router selects the approved layer and keeps the public API independent of the model provider.
type Router struct {
	registry *LayerRegistry
	runtimes map[string]LayerRuntime
}

func NewRouter(registry *LayerRegistry, runtimes map[string]LayerRuntime) *Router {
	return &Router{registry: registry, runtimes: runtimes}
}

func (r *Router) Handle(ctx context.Context, request RouteRequest) (Response, error) {
	layer, err := r.registry.Resolve(request.Category)
	if err != nil { return Response{}, err }
	runtime, ok := r.runtimes[layer.ID]
	if !ok { return Response{}, &layerUnavailableError{ID: layer.ID} }
	return runtime.Run(ctx, request)
}

type layerUnavailableError struct{ ID string }
func (e *layerUnavailableError) Error() string { return "AI layer runtime unavailable: " + e.ID }
