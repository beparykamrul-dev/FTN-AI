package apierror

import "encoding/json"

type Error struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	Details   any    `json:"details,omitempty"`
}

func (e Error) MarshalJSON() ([]byte, error) {
	type alias Error
	return json.Marshal(alias(e))
}

func New(code, message string) Error { return Error{Code: code, Message: message} }

func WithRequestID(e Error, requestID string) Error {
	e.RequestID = requestID
	return e
}
