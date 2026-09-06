package agent

import "testing"

func TestRouteRequestRequiresScopeIdentity(t *testing.T) {
	if _, err := NewRouter(nil); err == nil { t.Fatal("nil registry must fail") }
}
