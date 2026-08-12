package transport

import "context"

type Result string

const (
	Sent             Result = "SENT"
	Delivered        Result = "DELIVERED"
	TemporaryFailure Result = "TEMPORARY_FAILURE"
	PermanentFailure Result = "PERMANENT_FAILURE"
)

type Message struct {
	ID        string
	Sender    string
	Recipient string
	Body      string
}

// Adapter is implemented by FTN-controlled GSM modem and authorized operator SMPP transports.
type Adapter interface {
	Name() string
	Ready(ctx context.Context) bool
	Send(ctx context.Context, message Message) (Result, error)
}
