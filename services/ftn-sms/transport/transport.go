package transport

import "context"

// Result is the normalized outcome returned by an SMS transport adapter.
type Result string

const (
	ResultSent             Result = "SENT"
	ResultDelivered        Result = "DELIVERED"
	ResultTemporaryFailure Result = "TEMPORARY_FAILURE"
	ResultPermanentFailure Result = "PERMANENT_FAILURE"
)

// Message contains only data required by an already-authorized transport submission.
type Message struct {
	ID        string
	Sender    string
	Recipient string
	Body      string
}

// Adapter is the common boundary for FTN-controlled GSM modem and authorized SMPP transports.
// IAM, Sender-ID approval, rate limiting, and approval policy remain outside this interface.
type Adapter interface {
	Name() string
	Ready(ctx context.Context) bool
	Send(ctx context.Context, message Message) (Result, error)
}
