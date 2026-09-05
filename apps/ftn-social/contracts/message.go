package social

import (
	"context"
	"strings"
)

type Message struct { ID string `json:"id"`; ThreadID string `json:"thread_id"`; SenderID string `json:"sender_id"`; Body string `json:"body"`; CreatedAt string `json:"created_at"` }
type MessageStore interface { Create(context.Context, Message) error; Get(context.Context,string)(Message,error); ListThread(context.Context,string,int)([]Message,error) }
func (m Message) Valid() bool { id:=strings.TrimSpace(m.ID);thread:=strings.TrimSpace(m.ThreadID);sender:=strings.TrimSpace(m.SenderID);body:=strings.TrimSpace(m.Body);return id!=""&&thread!=""&&sender!=""&&body!=""&&len(id)<=256&&len(thread)<=256&&len(sender)<=256&&len(body)<=1<<20&&len(m.CreatedAt)<=128 }
func (m Message) Normalize() Message { m.ID=strings.TrimSpace(m.ID);m.ThreadID=strings.TrimSpace(m.ThreadID);m.SenderID=strings.TrimSpace(m.SenderID);m.Body=strings.TrimSpace(m.Body);m.CreatedAt=strings.TrimSpace(m.CreatedAt);return m }
