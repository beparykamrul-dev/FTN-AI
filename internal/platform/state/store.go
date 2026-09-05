package state

import (
	"context"
	"errors"
	"strings"
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

func ValidateDecision(d Decision) error {
	if strings.TrimSpace(d.ServiceID) == "" || strings.TrimSpace(d.PolicyVersion) == "" {
		return errors.New("service_id and policy_version are required")
	}
	if len(strings.TrimSpace(d.ServiceID)) > 256 || len(strings.TrimSpace(d.PeerID)) > 256 || len(strings.TrimSpace(d.PolicyVersion)) > 256 || len(strings.TrimSpace(d.Reason)) > 4096 {
		return errors.New("decision field is too large")
	}
	if strings.TrimSpace(d.Status) == "" {
		return errors.New("decision status is required")
	}
	// Transport is optional for compatibility with decisions that describe
	// service selection before a concrete transport has been negotiated.
	return nil
}
