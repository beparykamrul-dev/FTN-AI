package proxy

import "testing"

func TestDefaultHealthCheckIsValid(t *testing.T) {
	if !DefaultHealthCheck().Valid() { t.Fatal("default health policy must be valid") }
}

func TestHealthCheckRejectsTimeoutAboveInterval(t *testing.T) {
	h := DefaultHealthCheck(); h.Timeout = h.Interval + 1
	if h.Valid() { t.Fatal("timeout above interval must be invalid") }
}
