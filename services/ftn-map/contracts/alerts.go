package fiber

import "context"

type AlertLevel string
const (
	AlertInfo AlertLevel = "info"
	AlertWarning AlertLevel = "warning"
	AlertCritical AlertLevel = "critical"
)

type FiberAlert struct {
	ID string
	CoreID string
	PathID string
	Level AlertLevel
	Reason string
	Confidence float64
	CreatedAt string
	Acknowledged bool
}

type AlertRepository interface {
	Emit(context.Context, FiberAlert) error
	Acknowledge(context.Context, string, string) error
	Active(context.Context, string) ([]FiberAlert, error)
}

// AlertRouter fans out local alerts to authorized FTN notification channels.
type AlertRouter interface {
	Route(context.Context, FiberAlert) error
}
