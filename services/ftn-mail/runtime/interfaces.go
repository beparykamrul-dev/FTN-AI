package runtime

import "context"

type Mailbox struct {
	ID      string
	Address string
	Enabled bool
}

type Authenticator interface {
	Authenticate(ctx context.Context, address, secret string) (Mailbox, error)
}

type MailStore interface {
	Append(ctx context.Context, mailboxID string, message []byte) (string, error)
	List(ctx context.Context, mailboxID string, limit, offset int) ([]string, error)
	Get(ctx context.Context, mailboxID, messageID string) ([]byte, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, eventType, sourceID string, payload []byte) error
}
