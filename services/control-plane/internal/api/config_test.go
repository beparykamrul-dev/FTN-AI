package api

import "testing"

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("FTN_CONTROL_PLANE_ADDR", "")
	t.Setenv("FTN_REQUEST_ID_HEADER", "")
	cfg := LoadConfig()
	if cfg.Addr != ":8080" { t.Fatalf("unexpected addr: %s", cfg.Addr) }
	if cfg.RequestIDHeader != "X-Request-ID" { t.Fatalf("unexpected request header: %s", cfg.RequestIDHeader) }
}
