package alert

import "context"

type EscalationStage string
const (
	StageNOC EscalationStage = "noc"
	StageAICall EscalationStage = "ai-call"
	StageBackgroundAlert EscalationStage = "background-alert"
	StageEngineer EscalationStage = "engineer"
)

type EscalationPolicy struct {
	NoResponseAfterSeconds int64
	CallAfterSeconds int64
	BackgroundAlert bool
	RequireAcknowledgement bool
}

type Incident struct {
	ID string
	Severity string
	Stage EscalationStage
	Acknowledged bool
	NoResponseSeconds int64
}

// EscalationService coordinates FTN alerts. It does not bypass OS notification
// controls or silently defeat user privacy settings.
type EscalationService interface {
	Evaluate(context.Context, Incident, EscalationPolicy) (EscalationStage, error)
	Acknowledge(context.Context, string) error
	Escalate(context.Context, string, EscalationStage) error
}
