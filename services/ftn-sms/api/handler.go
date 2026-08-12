package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Authorizer is the existing FTN Control/IAM boundary. SMS does not own identities.
type Authorizer interface {
	Authorize(ctx context.Context, capability string) error
}

type Enqueuer interface {
	Enqueue(ctx context.Context, senderID, recipient, body string, scheduledAt any, priority int) (string, error)
}

type SendRequest struct {
	SenderID    string `json:"sender_id"`
	Recipient   string `json:"recipient"`
	Body        string `json:"body"`
	ScheduledAt any    `json:"scheduled_at,omitempty"`
	Priority    int    `json:"priority,omitempty"`
}

type Handler struct {
	Auth  Authorizer
	Queue Enqueuer
}

func (h Handler) Send(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil || h.Queue == nil {
		writeError(w, http.StatusServiceUnavailable, "SMS_SERVICE_UNAVAILABLE")
		return
	}
	if err := h.Auth.Authorize(r.Context(), "sms.send"); err != nil {
		writeError(w, http.StatusForbidden, "SMS_PERMISSION_DENIED")
		return
	}
	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON")
		return
	}
	req.SenderID = strings.TrimSpace(req.SenderID)
	req.Recipient = strings.TrimSpace(req.Recipient)
	if req.SenderID == "" || req.Recipient == "" || strings.TrimSpace(req.Body) == "" || len(req.Body) > 4096 || len(req.SenderID) > 32 || len(req.Recipient) > 32 || req.Priority < 0 || req.Priority > 100 {
		writeError(w, http.StatusBadRequest, "INVALID_SMS_REQUEST")
		return
	}
	id, err := h.Queue.Enqueue(r.Context(), req.SenderID, req.Recipient, req.Body, req.ScheduledAt, req.Priority)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			writeError(w, http.StatusServiceUnavailable, "SMS_QUEUE_UNAVAILABLE")
			return
		}
		writeError(w, http.StatusServiceUnavailable, "SMS_ENQUEUE_FAILED")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]any{"message_id": id, "status": "QUEUED"})
}

func writeError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"code": code, "message": code})
}
