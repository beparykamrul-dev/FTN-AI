package ws

import "time"

const (
	MessageHello       = "hello"
	MessageWelcome     = "welcome"
	MessageSubscribe   = "subscribe"
	MessageUnsubscribe = "unsubscribe"
	MessageEvent       = "event"
	MessageAck         = "ack"
	MessageResume      = "resume"
	MessageResumed     = "resumed"
	MessageResync      = "resync_required"
	MessageError       = "error"
)

type ClientHello struct {
	Protocol string `json:"protocol"`
	ClientID string `json:"client_id"`
	Token    string `json:"token,omitempty"`
	LastID  uint64 `json:"last_id,omitempty"`
}

type ServerWelcome struct {
	Protocol    string    `json:"protocol"`
	SessionID   string    `json:"session_id"`
	Heartbeat   time.Duration `json:"heartbeat_ms"`
	ServerTime  time.Time `json:"server_time"`
	ResumeUntil time.Time `json:"resume_until"`
}

type Ack struct {
	EventID uint64 `json:"event_id"`
}

type ResumeResult struct {
	FromID uint64 `json:"from_id"`
	ToID   uint64 `json:"to_id"`
	Count  int    `json:"count"`
}
