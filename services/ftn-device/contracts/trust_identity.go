package device

import "context"

type TrustState string

const (
	TrustUnknown TrustState = "unknown"
	TrustPending TrustState = "pending"
	TrustVerified TrustState = "verified"
	TrustRevoked TrustState = "revoked"
)

type DeviceIdentity struct {
	DeviceID string
	MAC string
	SerialNumber string
	CertificateFingerprint string
	Vendor string
	Model string
	Trust TrustState
}

// TrustIdentityService is the canonical gate between discovery and privileged
// provisioning. Verification alone never grants a device unrestricted access.
type TrustIdentityService interface {
	Register(context.Context, DeviceIdentity) error
	Verify(context.Context, string) (TrustState, error)
	Revoke(context.Context, string, string) error
	Get(context.Context, string) (DeviceIdentity, error)
}
