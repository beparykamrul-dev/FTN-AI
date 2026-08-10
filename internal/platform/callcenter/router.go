package callcenter

import "sync"

type Handler func(CallEvent) error

type Router struct {
	mu sync.RWMutex
	handlers map[string][]Handler
}

func NewRouter() *Router { return &Router{handlers: make(map[string][]Handler)} }

func (r *Router) Register(eventType string, h Handler) {
	if eventType == "" || h == nil { return }
	r.mu.Lock()
	r.handlers[eventType] = append(r.handlers[eventType], h)
	r.mu.Unlock()
}

// Dispatch fans a call-center event to registered handlers. One handler
// failure does not prevent the remaining handlers from receiving the event.
func (r *Router) Dispatch(e CallEvent) []error {
	r.mu.RLock()
	handlers := append([]Handler(nil), r.handlers[e.Type]...)
	r.mu.RUnlock()
	errs := make([]error, 0)
	for _, h := range handlers {
		if err := h(e); err != nil { errs = append(errs, err) }
	}
	return errs
}
