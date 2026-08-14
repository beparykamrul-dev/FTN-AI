package alert

import "context"

type SMSRoute struct {
	Provider string
	Endpoint string
	Encrypted bool
	FTNOwned bool
	ExternalFallbackAllowed bool
}

type MaskedSMS struct {
	IncidentID string
	RecipientRef string
	Message string
	Route SMSRoute
}

// PrivateSMSService keeps alert delivery on FTN-controlled infrastructure.
// External fallback is disabled by policy; outbound delivery is not assumed.
type PrivateSMSService interface {
	Mask(context.Context, string) (string, error)
	Send(context.Context, MaskedSMS) error
	Status(context.Context, string) (string, error)
}
