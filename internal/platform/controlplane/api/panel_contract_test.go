package api

import (
	"testing"
	"time"

	"github.com/beparykamrul-dev/FTN-AI/internal/platform/controlplane"
)

func TestBuildPanelSummary(t *testing.T) {
	now := time.Unix(10, 0)
	snapshot := controlplane.MonitoringSnapshot{
		GeneratedAt: now,
		ClusterState: controlplane.HealthDegraded,
		Nodes: []controlplane.NodeHealth{
			{NodeID: "n1", State: controlplane.HealthHealthy},
			{NodeID: "n2", State: controlplane.HealthDegraded},
			{NodeID: "n3", State: controlplane.HealthOffline},
		},
		Services: []controlplane.ServiceStatus{
			{ServiceID: "dns", State: controlplane.HealthHealthy},
			{ServiceID: "monitoring", State: controlplane.HealthDegraded},
		},
		ActiveAlerts: 4,
		CriticalAlerts: 1,
	}

	got := BuildPanelSummary(snapshot)
	if got.ContractVersion != ContractVersion {
		t.Fatalf("contract version = %q, want %q", got.ContractVersion, ContractVersion)
	}
	if got.NodeCount != 3 || got.HealthyNodes != 1 || got.DegradedNodes != 1 || got.OfflineNodes != 1 {
		t.Fatalf("unexpected node summary: %+v", got)
	}
	if got.ServiceCount != 2 || got.HealthyServices != 1 {
		t.Fatalf("unexpected service summary: %+v", got)
	}
	if got.ActiveAlerts != 4 || got.CriticalAlerts != 1 {
		t.Fatalf("unexpected alert summary: %+v", got)
	}
}
