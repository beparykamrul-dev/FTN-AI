package router

import "testing"

func TestValidateRouteRejectsInvalidPrefix(t *testing.T) {
	if err := validateRoute(Route{Prefix: "not-a-prefix"}); err == nil {
		t.Fatal("invalid prefix must be rejected")
	}
}

func TestValidateRouteAcceptsValidPrefix(t *testing.T) {
	if err := validateRoute(Route{Prefix: "10.0.0.0/24", NextHop: "10.0.0.1"}); err != nil {
		t.Fatalf("valid route rejected: %v", err)
	}
}
