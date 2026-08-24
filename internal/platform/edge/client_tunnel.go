package edge

import (
    "context"
    "fmt"
)

type ClientTunnelState struct {
    DeviceID string `json:"device_id"`
    NodeID string `json:"node_id"`
    State string `json:"state"`
    DNSMode string `json:"dns_mode"`
    PolicyVersion string `json:"policy_version"`
}

type ClientTunnel interface {
    Connect(context.Context, string) error
    Disconnect(context.Context, string) error
    Status(context.Context, string) (ClientTunnelState, error)
    Sync(context.Context, string) error
}

func ValidateTunnelState(s ClientTunnelState) error {
    if s.DeviceID == "" || s.NodeID == "" { return fmt.Errorf("tunnel identity is incomplete") }
    if s.State == "" { return fmt.Errorf("tunnel state is required") }
    return nil
}
