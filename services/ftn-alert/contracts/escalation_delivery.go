package alert

import "context"

type DeliveryChannel string
const (
	ChannelFTNSocial DeliveryChannel = "ftn-social"
	ChannelSMS DeliveryChannel = "sms"
	ChannelAICall DeliveryChannel = "ai-call"
	ChannelOSBackground DeliveryChannel = "os-background"
)

type DeliveryAttempt struct {
	IncidentID string
	Channel DeliveryChannel
	AttemptedAt string
	Acknowledged bool
	FailureReason string
}

type DeliveryPolicy struct {
	Channels []DeliveryChannel
	ResponseTimeoutSeconds int64
	CriticalOnlyBackground bool
	MaxCallAttempts int
}

type DeliveryService interface {
	Send(context.Context, string, DeliveryChannel) (DeliveryAttempt, error)
	Escalate(context.Context, string, DeliveryPolicy) error
	Acknowledge(context.Context, string) error
}
