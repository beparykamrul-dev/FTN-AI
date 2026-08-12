package api

import (
	"time"

	"github.com/beparykamrul-dev/FTN-AI/internal/platform/controlplane"
)

const ContractVersion = "v1"

// PanelSummary is deliberately compact so the web UI can render the primary
// operational state without coupling itself to individual data collectors.
type PanelSummary struct {
	ContractVersion string                        `json:"contract_version"`
	GeneratedAt     time.Time                     `json:"generated_at"`
	ClusterState    controlplane.HealthState      `json:"cluster_state"`
	NodeCount       int                           `json:"node_count"`
	HealthyNodes    int                           `json:"healthy_nodes"`
	DegradedNodes   int                           `json:"degraded_nodes"`
	OfflineNodes    int                           `json:"offline_nodes"`
	ServiceCount    int                           `json:"service_count"`
	HealthyServices int                           `json:"healthy_services"`
	ActiveAlerts    int                           `json:"active_alerts"`
	CriticalAlerts  int                           `json:"critical_alerts"`
}

// MonitoringPanelResponse is the read contract for the Monitoring Control
// Panel. Collectors can evolve independently as long as this contract stays
// stable.
type MonitoringPanelResponse struct {
	Summary   PanelSummary                  `json:"summary"`
	Snapshot  controlplane.MonitoringSnapshot `json:"snapshot"`
	UpdatedAt time.Time                     `json:"updated_at"`
}

// ControlPanelResponse is the read contract for the Control Panel. The panel
// exposes capabilities and current state; execution remains behind the policy
// and approval boundaries already defined by the control plane.
type ControlPanelResponse struct {
	ContractVersion string                       `json:"contract_version"`
	GeneratedAt     time.Time                    `json:"generated_at"`
	Features        []controlplane.Feature       `json:"features"`
	Commands       []controlplane.ControlCommand `json:"commands,omitempty"`
	Targets        []controlplane.TargetServer   `json:"targets,omitempty"`
}
