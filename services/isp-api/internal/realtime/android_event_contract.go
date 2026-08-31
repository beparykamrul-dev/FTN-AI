package realtime

type AndroidEventType string

const (
	EventAccountUpdated   AndroidEventType = "account.updated"
	EventBillingUpdated   AndroidEventType = "billing.updated"
	EventPaymentUpdated   AndroidEventType = "payment.updated"
	EventTicketUpdated    AndroidEventType = "ticket.updated"
	EventIncidentCreated  AndroidEventType = "incident.created"
	EventDeviceStatus     AndroidEventType = "device.status"
	EventDiscoveryUpdated AndroidEventType = "discovery.updated"
	EventRecoveryUpdated  AndroidEventType = "recovery.updated"
	EventNotification     AndroidEventType = "notification.created"
)

type AndroidEvent struct {
	EventID   string           `json:"event_id"`
	Type      AndroidEventType `json:"type"`
	Channel   string           `json:"channel"`
	Sequence  uint64           `json:"sequence"`
	Timestamp string           `json:"timestamp"`
	Resource  string           `json:"resource"`
	Payload   map[string]any   `json:"payload,omitempty"`
}
