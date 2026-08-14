package fiber

import "context"

type DeviceType string
const (
	DeviceRouter DeviceType = "router"
	DeviceOLT DeviceType = "olt"
	DeviceONU DeviceType = "onu"
	DeviceSwitch DeviceType = "switch"
	DevicePOP DeviceType = "pop"
	DeviceClientNode DeviceType = "client-node"
)

type DiscoveredDevice struct {
	ID string
	Type DeviceType
	MAC string
	Serial string
	Address string
	Manufacturer string
	Model string
	Firmware string
	TLSReady bool
	Trusted bool
	Status string
}

type CertificateState struct {
	DeviceID string
	Issuer string
	Subject string
	NotBefore string
	NotAfter string
	Fingerprint string
	Installed bool
	Trusted bool
}

// Discovery discovers candidates; it does not grant trust.
type Discovery interface {
	Scan(context.Context, string) ([]DiscoveredDevice, error)
	Inspect(context.Context, DiscoveredDevice) (DiscoveredDevice, error)
}

// CertificateManager handles FTN-approved device certificate enrollment and
// installation. Trust is granted only after identity/policy validation.
type CertificateManager interface {
	Issue(context.Context, DiscoveredDevice) (CertificateState, error)
	Install(context.Context, DiscoveredDevice, CertificateState) error
	Renew(context.Context, string) (CertificateState, error)
	Revoke(context.Context, string) error
	Status(context.Context, string) (CertificateState, error)
}
