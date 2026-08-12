package ws

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

var (
	ErrUnauthorized = errors.New("websocket unauthorized")
	ErrInvalidToken = errors.New("websocket invalid token")
)

type AuthPolicy struct {
	RequireTLS     bool
	AllowedOrigins []string
	MaxTokenAge    time.Duration
}

type Credentials struct {
	Token     string
	IssuedAt  time.Time
	Subject   string
	NodeID    string
}

// NewSessionID creates a locally generated opaque identifier. It deliberately
// has no dependency on an external identity/session service.
func NewSessionID() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil { return "", err }
	return hex.EncodeToString(b), nil
}

func (p AuthPolicy) ValidateTLS(secure bool) error {
	if p.RequireTLS && !secure { return ErrUnauthorized }
	return nil
}

func (p AuthPolicy) ValidateOrigin(origin string) error {
	if len(p.AllowedOrigins) == 0 { return nil }
	for _, allowed := range p.AllowedOrigins {
		if subtle.ConstantTimeCompare([]byte(origin), []byte(allowed)) == 1 { return nil }
	}
	return ErrUnauthorized
}

func (p AuthPolicy) ValidateCredentials(c Credentials, now time.Time) error {
	if strings.TrimSpace(c.Token) == "" || strings.TrimSpace(c.Subject) == "" { return ErrInvalidToken }
	if p.MaxTokenAge > 0 && (c.IssuedAt.IsZero() || now.Sub(c.IssuedAt) > p.MaxTokenAge) { return ErrInvalidToken }
	return nil
}
