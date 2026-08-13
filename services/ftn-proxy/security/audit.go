package security

import "time"

// AuditEvent is a metadata-only security event. Sensitive credentials and
// payload contents must never be recorded here.
type AuditEvent struct {
	RequestID string
	Service   string
	Action    string
	Allowed   bool
	Reason    string
	At        time.Time
}

// AuditSink is intentionally small so it can later be backed by FTN Metrics,
// SIEM, immutable storage, or another approved security system.
type AuditSink interface {
	Record(AuditEvent)
}

func NormalizeEvent(e AuditEvent) AuditEvent {
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	return e
}
