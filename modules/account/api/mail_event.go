package api

import (
	"encoding/json"
	"net/http"
)

type MailEvent struct {
	ID         string          `json:"id"`
	Source     string          `json:"source"`
	SourceID   string          `json:"source_id"`
	EventType  string          `json:"event_type"`
	Payload    json.RawMessage `json:"payload"`
	Confidence float64         `json:"confidence"`
	Status     string          `json:"status"`
}

type MailEventProcessor interface {
	Process(MailEvent) error
}

type MailEventAPI struct {
	Processor MailEventProcessor
}

func (a MailEventAPI) Process(w http.ResponseWriter, r *http.Request) {
	var event MailEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid mail event", http.StatusBadRequest)
		return
	}
	if event.ID == "" || event.Source == "" || event.EventType == "" {
		http.Error(w, "id, source and event_type are required", http.StatusBadRequest)
		return
	}
	if err := a.Processor.Process(event); err != nil {
		http.Error(w, "mail event processing failed", http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
