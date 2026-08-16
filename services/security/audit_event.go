package security

// AuditEvent records a security-control decision without storing secrets or source content.
type AuditEvent struct {
	EventID string
	Action string
	Actor string
	Target string
	Allowed bool
	Reason string
}

func (e AuditEvent) Valid() bool {
	return e.EventID != "" && e.Action != "" && e.Actor != "" && e.Target != ""
}
