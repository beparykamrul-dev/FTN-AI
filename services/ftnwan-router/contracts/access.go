package router

import "context"

type AccessDevice struct {
	DeviceID      string
	Surface       AccessSurface
	MAC           string
	Serial        string
	Authenticated bool
	Authorized    bool
	Online        bool
}

type TunnelProfile struct {
	ID          string
	Protocol    string
	Endpoint    string
	AllowedCIDR []string
	DNS         []string
	Enabled     bool
}

// AccessController is the provider-neutral enrollment/tunnel boundary.
type AccessController interface {
	Discover(context.Context, string) (AccessDevice, error)
	Enroll(context.Context, AccessDevice) (string, error)
	ApplyProfile(context.Context, string, TunnelProfile) error
	Revoke(context.Context, string) error
	Status(context.Context, string) (AccessDevice, error)
}
