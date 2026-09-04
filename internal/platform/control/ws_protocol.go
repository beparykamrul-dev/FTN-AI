package control

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type WSMessage struct {
	Type      string          `json:"type"`
	RequestID string          `json:"request_id,omitempty"`
	ServerID  string          `json:"server_id,omitempty"`
	Operation Operation       `json:"operation,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

const maxWSMessageBytes = 256 << 10

func DecodeMessage(data []byte) (WSMessage, error) {
	if len(data) == 0 { return WSMessage{}, errors.New("message is empty") }
	if len(data) > maxWSMessageBytes { return WSMessage{}, errors.New("message too large") }
	var m WSMessage
	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&m); err != nil { return m, err }
	var extra any
	if err := dec.Decode(&extra); err != nil {
		if !errors.Is(err, json.EOF) { return m, fmt.Errorf("invalid trailing message data: %w", err) }
	} else { return m, errors.New("multiple JSON values") }
	if strings.TrimSpace(m.Type) == "" { return m, errors.New("message type is required") }
	if len(m.RequestID) > 256 || len(m.ServerID) > 256 { return m, errors.New("message identifier too long") }
	if len(m.Payload) > maxWSMessageBytes { return m, errors.New("message payload too large") }
	return m, nil
}

func EncodeMessage(m WSMessage) ([]byte, error) { return json.Marshal(m) }
