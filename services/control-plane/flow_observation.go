package main

import (
	"errors"
	"net"
	"strings"
	"time"
)

// FlowObservation is the normalized telemetry envelope used between adapters
// and durable stores. It keeps the flow counters immutable while carrying the
// source/export timestamp independently from request arrival time.
type FlowObservation struct {
	FlowRecord
	ObservedAt time.Time `json:"observed_at"`
	FirstSeen  time.Time `json:"first_seen,omitempty"`
	LastSeen   time.Time `json:"last_seen,omitempty"`
}

func (o FlowObservation) Validate() error {
	if strings.TrimSpace(o.ExporterID) == "" {
		return errors.New("exporter_required")
	}
	if net.ParseIP(strings.TrimSpace(o.SourceIP)) == nil || net.ParseIP(strings.TrimSpace(o.DestinationIP)) == nil {
		return errors.New("flow_source_destination_required")
	}
	if o.Version != 5 && o.Version != 9 && o.Version != 10 {
		return errors.New("unsupported_flow_version")
	}
	if o.SamplingRate == 0 {
		return errors.New("invalid_sampling_rate")
	}
	if o.ObservedAt.IsZero() {
		return errors.New("observed_at_required")
	}
	if !o.FirstSeen.IsZero() && !o.LastSeen.IsZero() && o.LastSeen.Before(o.FirstSeen) {
		return errors.New("invalid_flow_time_range")
	}
	return nil
}

func normalizeFlowObservation(f FlowRecord, observedAt, firstSeen, lastSeen time.Time) (FlowObservation, error) {
	if observedAt.IsZero() {
		observedAt = time.Now().UTC()
	}
	observedAt = observedAt.UTC()
	if !firstSeen.IsZero() {
		firstSeen = firstSeen.UTC()
	}
	if !lastSeen.IsZero() {
		lastSeen = lastSeen.UTC()
	}
	o := FlowObservation{FlowRecord: f, ObservedAt: observedAt, FirstSeen: firstSeen, LastSeen: lastSeen}
	if err := o.Validate(); err != nil {
		return FlowObservation{}, err
	}
	return o, nil
}
