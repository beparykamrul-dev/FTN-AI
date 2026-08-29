package payment

import (
	"errors"
	"strings"
)

var ErrInvalidIdempotencyKey = errors.New("invalid idempotency key")

func ValidateIdempotencyKey(key string) error {
	key = strings.TrimSpace(key)
	if len(key) < 16 || len(key) > 128 {
		return ErrInvalidIdempotencyKey
	}
	return nil
}
