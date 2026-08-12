package api

// Stable HTTP/WebSocket route names for the first two panels. Feature modules
// added later must register beneath these surfaces instead of creating a new
// top-level control panel.
const (
	RouteControlHealth      = "/api/v1/control/health"
	RouteControlSummary     = "/api/v1/control/summary"
	RouteControlFeatures    = "/api/v1/control/features"
	RouteControlTargets     = "/api/v1/control/targets"
	RouteControlCommands    = "/api/v1/control/commands"
	RouteControlApprovals   = "/api/v1/control/approvals"
	RouteControlAudit       = "/api/v1/control/audit"
	RouteControlStream      = "/api/v1/control/stream"

	RouteMonitoringSummary  = "/api/v1/monitoring/summary"
	RouteMonitoringNodes    = "/api/v1/monitoring/nodes"
	RouteMonitoringServices = "/api/v1/monitoring/services"
	RouteMonitoringAlerts   = "/api/v1/monitoring/alerts"
	RouteMonitoringMetrics  = "/api/v1/monitoring/metrics"
	RouteMonitoringLogs     = "/api/v1/monitoring/logs"
	RouteMonitoringTopology = "/api/v1/monitoring/topology"
	RouteMonitoringStream   = "/api/v1/monitoring/stream"
)
