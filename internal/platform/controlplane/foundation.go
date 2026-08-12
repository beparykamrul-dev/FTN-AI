package controlplane

import "time"

// HealthState is intentionally shared by Control and Monitoring panels so
// operational state has one vocabulary across the platform.
type HealthState string

const (
	HealthUnknown  HealthState = "unknown"
	HealthHealthy  HealthState = "healthy"
	HealthDegraded HealthState = "degraded"
	HealthCritical HealthState = "critical"
	HealthOffline  HealthState = "offline"
)

// NodeHealth is the normalized health snapshot reported by an FTN agent.
type NodeHealth struct {
	NodeID       string      `json:"node_id"`
	State        HealthState `json:"state"`
	CPUPercent   float64     `json:"cpu_percent"`
	MemoryPercent float64    `json:"memory_percent"`
	DiskPercent  float64     `json:"disk_percent"`
	NetworkMbps  float64     `json:"network_mbps"`
	LatencyMS    float64     `json:"latency_ms"`
	PacketLoss   float64     `json:"packet_loss_percent"`
	ServicesUp   int         `json:"services_up"`
	ServicesDown int         `json:"services_down"`
	LastSeen     time.Time   `json:"last_seen"`
}

// ServiceStatus makes services portable across nodes. Placement is runtime
// state; service identity and desired state remain independent of a server.
type ServiceStatus struct {
	ServiceID    string      `json:"service_id"`
	Name         string      `json:"name"`
	Desired      int         `json:"desired"`
	Ready        int         `json:"ready"`
	State        HealthState `json:"state"`
	Nodes        []string    `json:"nodes,omitempty"`
	LastReconcile time.Time  `json:"last_reconcile"`
}

// MonitoringSnapshot is the single read model shared by the two panels.
type MonitoringSnapshot struct {
	GeneratedAt   time.Time       `json:"generated_at"`
	ClusterState  HealthState     `json:"cluster_state"`
	Nodes         []NodeHealth    `json:"nodes"`
	Services      []ServiceStatus `json:"services"`
	ActiveAlerts  int             `json:"active_alerts"`
	CriticalAlerts int            `json:"critical_alerts"`
}

// ControlCommand is the normalized action contract. Adapters execute actions;
// the control plane decides policy, target and approval requirements.
type ControlCommand struct {
	ID          string            `json:"id"`
	Action      string            `json:"action"`
	Target      string            `json:"target"`
	Parameters  map[string]string `json:"parameters,omitempty"`
	RequestedBy string            `json:"requested_by"`
	CreatedAt   time.Time         `json:"created_at"`
	Risk        CommandRisk       `json:"risk"`
	Approval    ApprovalState     `json:"approval"`
}

type CommandRisk string

const (
	RiskRead     CommandRisk = "read"
	RiskLow      CommandRisk = "low"
	RiskMedium   CommandRisk = "medium"
	RiskHigh     CommandRisk = "high"
	RiskCritical CommandRisk = "critical"
)

type ApprovalState string

const (
	ApprovalNotRequired ApprovalState = "not_required"
	ApprovalPending     ApprovalState = "pending"
	ApprovalApproved    ApprovalState = "approved"
	ApprovalRejected    ApprovalState = "rejected"
)

// ReconcileResult records the outcome of desired-state reconciliation without
// coupling the control plane to a particular service runtime.
type ReconcileResult struct {
	ResourceID string      `json:"resource_id"`
	DesiredHash string     `json:"desired_hash"`
	ActualHash  string     `json:"actual_hash"`
	Changed     bool       `json:"changed"`
	State       HealthState `json:"state"`
	Message     string     `json:"message,omitempty"`
	At          time.Time  `json:"at"`
}
