package social

import "context"

type Message struct {
	ID        string `json:"id"`
	ThreadID  string `json:"thread_id"`
	SenderID  string `json:"sender_id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type MessageStore interface {
	Create(context.Context, Message) error
	Get(context.Context, string) (Message, error)
	ListThread(context.Context, string, int) ([]Message, error)
}
