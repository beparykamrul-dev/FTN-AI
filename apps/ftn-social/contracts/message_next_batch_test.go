package social

import "testing"

func TestMessageRejectsControlCharacters(t *testing.T) {
	m := Message{ID:"m", ThreadID:"t", SenderID:"s", Body:"hello\x00world", CreatedAt:"now"}
	if m.Valid() { t.Fatal("NUL-containing message must be invalid") }
	m.Body = " hello "
	if !m.Normalize().Valid() { t.Fatal("normalized valid message rejected") }
}
