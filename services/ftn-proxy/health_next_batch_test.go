package proxy

import (
	"testing"
	"time"
)

func TestHealthCheckRejectsTimeoutAfterInterval(t *testing.T) {
	h := DefaultHealthCheck()
	h.Timeout = h.Interval + time.Second
	if h.Valid() { t.Fatal("timeout greater than interval must be invalid") }
}

func TestHealthTrackerTransitionsAfterThreshold(t *testing.T) {
	p := HealthCheck{Interval:time.Second, Timeout:time.Second, FailureLimit:2, SuccessLimit:2}
	var h HealthTracker
	if h.Observe(false, p) != HealthUnknown { t.Fatal("first failure should not mark unhealthy") }
	if h.Observe(false, p) != HealthUnhealthy { t.Fatal("failure threshold not reached") }
	if h.Observe(true, p) != HealthUnhealthy { t.Fatal("state should remain unhealthy until success threshold") }
	if h.Observe(true, p) != HealthHealthy { t.Fatal("success threshold not reached") }
}
