package ftnmap

import "time"

// StaleDetector identifies topology entries that have not been observed
// within the configured freshness window.
type StaleDetector struct{ MaxAge time.Duration }

func (d StaleDetector) Stale(e TopologyEvent, now time.Time) bool {
	if d.MaxAge <= 0 || e.ObservedAt.IsZero() { return true }
	return now.UTC().Sub(e.ObservedAt) > d.MaxAge
}
