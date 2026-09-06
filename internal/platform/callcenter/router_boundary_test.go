package callcenter

import (
	"errors"
	"testing"
)

func TestRouterDeduplicatesHandlers(t *testing.T) {
	r := NewRouter()
	count := 0
	h := func(CallEvent) error { count++; return nil }
	r.Register("connected", h)
	r.Register("connected", h)
	if errs := r.Dispatch(CallEvent{Type: "connected"}); len(errs) != 0 || count != 1 {
		t.Fatalf("handler was not deduplicated: errs=%v count=%d", errs, count)
	}
}

func TestRouterReturnsHandlerErrors(t *testing.T) {
	r := NewRouter()
	r.Register("failed", func(CallEvent) error { return errors.New("boom") })
	if errs := r.Dispatch(CallEvent{Type: "failed"}); len(errs) != 1 || errs[0].Error() != "boom" {
		t.Fatalf("unexpected dispatch errors: %v", errs)
	}
}
