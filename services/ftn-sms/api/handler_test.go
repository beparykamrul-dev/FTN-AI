package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type denyAuthorizer struct{}

func (denyAuthorizer) Authorize(context.Context, string) error { return context.Canceled }

type testQueue struct{ called bool }

func (q *testQueue) Enqueue(context.Context, string, string, string, any, int) (string, error) {
	q.called = true
	return "unexpected", nil
}

func TestSendRejectsWithoutCapability(t *testing.T) {
	q := &testQueue{}
	h := Handler{Auth: denyAuthorizer{}, Queue: q}

	r := httptest.NewRequest(http.MethodPost, "/v1/sms/messages", strings.NewReader(`{"sender_id":"FTN","recipient":"+8801700000000","body":"test"}`))
	w := httptest.NewRecorder()

	h.Send(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	if q.called {
		t.Fatal("queue must not be called when authorization fails")
	}
}
