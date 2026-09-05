package realtime

import (
	"encoding/json"
	"fmt"
)

type EventType string

const (
	EventMessage  EventType = "message"
	EventAck      EventType = "ack"
	EventTyping   EventType = "typing"
	EventPresence EventType = "presence"
)

const MaxEventBytes = 1 << 20

type Event struct {
	Type      EventType       `json:"type"`
	MessageID string          `json:"message_id,omitempty"`
	RoomID    string          `json:"room_id,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

func validEventType(t EventType) bool {
	switch t {
	case EventMessage, EventAck, EventTyping, EventPresence:
		return true
	}
	return false
}

func EncodeEvent(event Event) ([]byte, error) {
	if !validEventType(event.Type) {
		return nil, fmt.Errorf("invalid event type")
	}
	if len(event.MessageID) > 256 || len(event.RoomID) > 256 || len(event.Payload) > MaxEventBytes {
		return nil, fmt.Errorf("event exceeds limits")
	}
	b, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	if len(b) > MaxEventBytes {
		return nil, fmt.Errorf("event exceeds %d bytes", MaxEventBytes)
	}
	return b, nil
}

func DecodeEvent(data []byte) (Event, error) {
	if len(data) == 0 || len(data) > MaxEventBytes {
		return Event{}, fmt.Errorf("event size out of bounds")
	}
	var e Event
	if err := json.Unmarshal(data, &e); err != nil {
		return Event{}, fmt.Errorf("decode realtime event: %w", err)
	}
	if !validEventType(e.Type) {
		return Event{}, fmt.Errorf("invalid event type")
	}
	if len(e.MessageID) > 256 || len(e.RoomID) > 256 || len(e.Payload) > MaxEventBytes {
		return Event{}, fmt.Errorf("event field exceeds limits")
	}
	return e, nil
}
