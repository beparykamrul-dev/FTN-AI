package mesh

// DeviceProfile describes the network capabilities and enrollment policy that
// FTN can assign to a managed device. Secrets and private keys are deliberately
// excluded; those are issued by the authenticated enrollment service.
type DeviceProfile struct {
	ID string `json:"id"`
	DeviceKind string `json:"device_kind"`
	MeshEnabled bool `json:"mesh_enabled"`
	WireGuardEnabled bool `json:"wireguard_enabled"`
	NetBirdEnabled bool `json:"netbird_enabled"`
	KeepaliveSeconds uint32 `json:"keepalive_seconds"`
	AllowedGroups []string `json:"allowed_groups,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

func DefaultDeviceProfile(kind string) DeviceProfile {
	return DeviceProfile{
		ID: "default-" + kind,
		DeviceKind: kind,
		MeshEnabled: true,
		WireGuardEnabled: true,
		NetBirdEnabled: true,
		KeepaliveSeconds: 15,
	}
}

// CanJoin reports whether a profile permits mesh enrollment. Enrollment still
// requires authenticated identity and an authorized enrollment token.
func (p DeviceProfile) CanJoin() bool {
	return p.MeshEnabled && (p.WireGuardEnabled || p.NetBirdEnabled)
}
