package control

import (
	"encoding/json"
	"errors"
)

type WSMessage struct {
	Type      string          `json:"type"`
	RequestID string          `json:"request_id,omitempty"`
	ServerID  string          `json:"server_id,omitempty"`
	Operation Operation       `json:"operation,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

func DecodeMessage(data []byte) (WSMessage, error) {
	var m WSMessage
	if err := json.Unmarshal(data, &m); err != nil { return m, err }
	if m.Type == "" { return m, errors.New("message type is required") }
	return m, nil
}

func EncodeMessage(m WSMessage) ([]byte, error) { return json.Marshal(m) }
