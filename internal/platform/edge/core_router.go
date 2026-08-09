package edge

import "context"

type RouterState struct {
	ID string `json:"id"`
	Hostname string `json:"hostname"`
	ManagementIP string `json:"management_ip"`
	Status string `json:"status"`
	Interfaces []InterfaceState `json:"interfaces,omitempty"`
}

type InterfaceState struct { Name string `json:"name"`; Address string `json:"address,omitempty"`; State string `json:"state"`; RXBPS float64 `json:"rx_bps,omitempty"`; TXBPS float64 `json:"tx_bps,omitempty"` }

type CoreRouter interface {
	Identity(ctx context.Context) (RouterState, error)
	Interfaces(ctx context.Context) ([]InterfaceState, error)
	Apply(ctx context.Context, request ChangeRequest) error
}

type ChangeRequest struct { ID string `json:"id"`; Action string `json:"action"`; Target string `json:"target"`; Reason string `json:"reason"`; ApprovedBy string `json:"approved_by"` }
