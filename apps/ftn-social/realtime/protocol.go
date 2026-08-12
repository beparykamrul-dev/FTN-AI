package realtime

import (
	"encoding/json"
	"fmt"
)

type EventType string

const (
	EventMessage EventType = "message"
	EventAck     EventType = "ack"
	EventTyping  EventType = "typing"
	EventPresence EventType = "presence"
)

type Event struct {
	Type      EventType       `json:"type"`
	MessageID string          `json:"message_id,omitempty"`
	RoomID    string          `json:"room_id,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

func EncodeEvent(event Event) ([]byte, error) {
	if event.Type == "" {
		return nil, fmt.Errorf("event type is required")
	}
	return json.Marshal(event)
}

func DecodeEvent(data []byte) (Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return Event{}, fmt.Errorf("decode realtime event: %w", err)
	}
	if event.Type == "" {
		return Event{}, fmt.Errorf("event type is required")
	}
	return event, nil
}
