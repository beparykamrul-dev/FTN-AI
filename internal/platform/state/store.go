package state

import (
	"context"
	"errors"
)

var ErrNotConfigured = errors.New("state store is not configured")

type Decision struct {
	ServiceID     string
	PeerID        string
	Transport     string
	PolicyVersion string
	Status        string
	Reason        string
}

type Store interface {
	SaveDecision(context.Context, Decision) error
}

// ValidateDecision enforces the persistence boundary before a decision is
// written. Persistence never implies permission to mutate the network.
func ValidateDecision(d Decision) error {
	if d.ServiceID == "" || d.PolicyVersion == "" {
		return errors.New("service_id and policy_version are required")
	}
	if d.Status == "" {
		return errors.New("decision status is required")
	}
	return nil
}
