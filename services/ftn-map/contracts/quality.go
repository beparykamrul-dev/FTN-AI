package fiber

import "context"

type QualityIssueType string

const (
	DuplicateDevice      QualityIssueType = "duplicate-device"
	InvalidTopology      QualityIssueType = "invalid-topology"
	InvalidCoordinate    QualityIssueType = "invalid-coordinate"
	DistanceMismatch     QualityIssueType = "distance-mismatch"
	MissingRelationship  QualityIssueType = "missing-relationship"
	StaleTelemetry       QualityIssueType = "stale-telemetry"
)

type QualityIssue struct {
	ID string
	Type QualityIssueType
	EntityID string
	Severity string
	Message string
	DetectedAt string
}

type TopologyQuality struct {
	Score float64
	Issues []QualityIssue
	CheckedAt string
}

type QualityAnalyzer interface {
	Validate(context.Context, FiberTopology) (TopologyQuality, error)
}
