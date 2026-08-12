package transport

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
)

var ErrAuthExchange = errors.New("SMTP AUTH exchange failed")

// ParsePLAIN validates an RFC 4616-style AUTH PLAIN initial response.
// It returns the login and password without storing them.
func ParsePLAIN(encoded string) (login, password string, err error) {
	if encoded == "" { return "", "", ErrAuthExchange }
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil { return "", "", ErrAuthExchange }
	parts := strings.Split(string(decoded), "\x00")
	if len(parts) != 3 || parts[1] == "" || parts[2] == "" { return "", "", ErrAuthExchange }
	return parts[1], parts[2], nil
}

type AuthExchange struct { Auth Authenticator }

func (a AuthExchange) AuthenticatePLAIN(ctx context.Context, encoded string) (string, error) {
	if a.Auth == nil { return "", ErrAuthExchange }
	login, password, err := ParsePLAIN(encoded)
	if err != nil { return "", err }
	return a.Auth.Authenticate(ctx, login, password)
}
