package edge

import "context"

type AndroidNetworkProfile struct { DeviceID string `json:"device_id"`; UserID string `json:"user_id"`; FTNNode string `json:"ftn_node"`; TunnelState string `json:"tunnel_state"`; DNSMode string `json:"dns_mode"`; PolicyVersion string `json:"policy_version"` }

type AndroidNetwork interface {
	Register(ctx context.Context, profile AndroidNetworkProfile) error
	Status(ctx context.Context, deviceID string) (AndroidNetworkProfile, error)
	SyncPolicy(ctx context.Context, deviceID string) error
}
