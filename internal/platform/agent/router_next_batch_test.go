package agent

import (
	"context"
	"testing"
)

func TestRouterHandleFailsClosed(t *testing.T) {
	r := NewRouter(nil, nil)
	if _, err := r.Handle(context.Background(), RouteRequest{Input:"status"}); err == nil { t.Fatal("missing registry must fail") }
	if _, err := r.Handle(nil, RouteRequest{Input:"status"}); err == nil { t.Fatal("nil context must fail") }
}
