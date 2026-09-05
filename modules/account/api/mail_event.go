package api

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strings"
)

type MailEvent struct { ID string `json:"id"`; Source string `json:"source"`; SourceID string `json:"source_id"`; EventType string `json:"event_type"`; Payload json.RawMessage `json:"payload"`; Confidence float64 `json:"confidence"`; Status string `json:"status"` }
type MailEventProcessor interface { Process(MailEvent) error }
type MailEventAPI struct { Processor MailEventProcessor }

func (a MailEventAPI) Process(w http.ResponseWriter, r *http.Request) {
	if w==nil||r==nil{return}
	if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	if a.Processor == nil { http.Error(w, "mail event processor unavailable", http.StatusServiceUnavailable); return }
	if r.Body == nil { http.Error(w, "invalid mail event", http.StatusBadRequest); return }
	var event MailEvent
	dec:=json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	if err := dec.Decode(&event); err != nil { http.Error(w, "invalid mail event", http.StatusBadRequest); return }
	var extra any
	if err:=dec.Decode(&extra); err!=io.EOF { http.Error(w,"multiple json values",http.StatusBadRequest); return }
	event.ID, event.Source, event.SourceID, event.EventType, event.Status = strings.TrimSpace(event.ID), strings.TrimSpace(event.Source), strings.TrimSpace(event.SourceID), strings.TrimSpace(event.EventType), strings.TrimSpace(event.Status)
	if len(event.ID)>256||len(event.Source)>128||len(event.SourceID)>256||len(event.EventType)>128||len(event.Status)>64 { http.Error(w,"mail event field too long",http.StatusBadRequest); return }
	if event.ID == "" || event.Source == "" || event.EventType == "" { http.Error(w, "id, source and event_type are required", http.StatusBadRequest); return }
	if math.IsNaN(event.Confidence) || math.IsInf(event.Confidence, 0) || event.Confidence < 0 || event.Confidence > 1 { http.Error(w, "confidence must be between 0 and 1", http.StatusBadRequest); return }
	if len(event.Payload) > 0 && !json.Valid(event.Payload) { http.Error(w, "payload must be valid JSON", http.StatusBadRequest); return }
	if err := a.Processor.Process(event); err != nil { http.Error(w, "mail event processing failed", http.StatusUnprocessableEntity); return }
	w.WriteHeader(http.StatusAccepted)
}
