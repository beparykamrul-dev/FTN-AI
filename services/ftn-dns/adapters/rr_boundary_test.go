package adapters

import "testing"

func TestParseAnswersRejectsTruncatedMessage(t *testing.T) {
	if _, err := ParseAnswers([]byte{0, 1, 0}); err == nil {
		t.Fatal("truncated DNS message must be rejected")
	}
}

func TestParseAnswersAcceptsHeaderOnlyWithoutAnswers(t *testing.T) {
	msg := make([]byte, 12)
	if _, err := ParseAnswers(msg); err != nil {
		t.Fatalf("header-only message should parse: %v", err)
	}
}
