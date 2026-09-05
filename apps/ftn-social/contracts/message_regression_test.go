package social

import (
	"strings"
	"testing"
)

func TestMessageRejectsOversizedBody(t *testing.T) {
	m := Message{ID:"m", ThreadID:"t", SenderID:"s", Body:strings.Repeat("x", (1<<20)+1)}
	if m.Valid() { t.Fatal("oversized message body must be rejected") }
}

func TestMessageNormalizeTrimsFields(t *testing.T) {
	m := (Message{ID:" m ",ThreadID:" t ",SenderID:" s ",Body:" body "}).Normalize()
	if m.ID!="m" || m.ThreadID!="t" || m.SenderID!="s" || m.Body!="body" { t.Fatalf("unexpected normalization: %#v", m) }
}
