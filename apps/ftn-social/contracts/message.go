package social

import (
	"context"
	"strings"
	"unicode"
)

type Message struct { ID string `json:"id"`; ThreadID string `json:"thread_id"`; SenderID string `json:"sender_id"`; Body string `json:"body"`; CreatedAt string `json:"created_at"` }
type MessageStore interface { Create(context.Context, Message) error; Get(context.Context,string)(Message,error); ListThread(context.Context,string,int)([]Message,error) }
func validText(v string,max int)bool{v=strings.TrimSpace(v);if v==""||len(v)>max{return false};for _,r:=range v{if unicode.IsControl(r)&&r!='\n'&&r!='\t'{return false}};return true}
func (m Message) Valid() bool { return validText(m.ID,256)&&validText(m.ThreadID,256)&&validText(m.SenderID,256)&&validText(m.Body,1<<20)&&len(m.CreatedAt)<=128 }
func (m Message) Normalize() Message { m.ID=strings.TrimSpace(m.ID);m.ThreadID=strings.TrimSpace(m.ThreadID);m.SenderID=strings.TrimSpace(m.SenderID);m.Body=strings.TrimSpace(m.Body);m.CreatedAt=strings.TrimSpace(m.CreatedAt);return m }
