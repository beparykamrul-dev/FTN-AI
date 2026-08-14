package social

import "context"

type Conversation struct {
	ID string
	Title string
	CreatedBy string
	CreatedAt string
}

type Message struct {
	ID string
	ConversationID string
	SenderID string
	Body string
	CreatedAt string
}

type Notification struct {
	ID string
	RecipientID string
	Title string
	Body string
	Severity string
	Source string
	CreatedAt string
	Read bool
}

// SocialStore is provider-neutral and supports a native FTN social system.
type SocialStore interface {
	CreateConversation(context.Context, Conversation) error
	AppendMessage(context.Context, Message) error
	CreateNotification(context.Context, Notification) error
	MarkRead(context.Context, string) error
}

// DeliveryProvider is an optional adapter. FTN Social remains functional
// without Telegram or another external messenger.
type DeliveryProvider interface {
	Name() string
	Deliver(context.Context, Notification) error
}

// Assistant is policy-bound; it cannot independently execute privileged
// network actions from a social conversation.
type Assistant interface {
	Respond(context.Context, string, string) (string, error)
}
