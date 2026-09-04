package main

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
)

type FTNBFDObservation struct {
	Peer string
	State FTNBFDState
}

type FTNBFDSource interface {
	ListBFD(context.Context) ([]FTNBFDObservation, error)
}

// FTNBFDReconciler publishes only normalized BFD state transitions. It does
// not create, delete, or reconfigure BFD sessions.
type FTNBFDReconciler struct {
	source   FTNBFDSource
	bridge   *FTNRoutedEventBridge
	previous map[string]FTNBFDState
}

func NewFTNBFDReconciler(source FTNBFDSource, bridge *FTNRoutedEventBridge) *FTNBFDReconciler {
	return &FTNBFDReconciler{source: source, bridge: bridge, previous: make(map[string]FTNBFDState)}
}

func (r *FTNBFDReconciler) Reconcile(ctx context.Context) error {
	if r == nil || r.source == nil || r.bridge == nil {
		return fmt.Errorf("BFD source and event bridge are required")
	}
	if ctx == nil {
		return fmt.Errorf("context is required")
	}
	observations, err := r.source.ListBFD(ctx)
	if err != nil {
		return err
	}
	seen := make(map[string]struct{}, len(observations))
	for _, observation := range observations {
		peer := strings.TrimSpace(observation.Peer)
		if _, err := netip.ParseAddr(peer); err != nil {
			continue
		}
		switch observation.State {
		case FTNBFDUp, FTNBFDDown, FTNBFDUnknown:
		default:
			continue
		}
		seen[peer] = struct{}{}
		if err := r.bridge.BFDState(ctx, peer, observation.State); err != nil {
			return err
		}
		r.previous[peer] = observation.State
	}
	for peer := range r.previous {
		if _, ok := seen[peer]; ok {
			continue
		}
		if err := r.bridge.BFDState(ctx, peer, FTNBFDUnknown); err != nil {
			return err
		}
		delete(r.previous, peer)
	}
	return nil
}
