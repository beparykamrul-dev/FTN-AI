package transport

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"
)

var ErrMessageTooLarge = errors.New("message too large")

type MessageSink interface {
	Save(ctx context.Context, sender string, recipients []string, raw []byte) error
}

// ReadMessage reads an already authenticated SMTP DATA stream with a hard size limit.
// Network/session handling and SMTP command parsing remain in the transport adapter.
func ReadMessage(r io.Reader, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 { return nil, ErrMessageTooLarge }
	lr := io.LimitReader(r, maxBytes+1)
	data, err := io.ReadAll(bufio.NewReader(lr))
	if err != nil { return nil, err }
	if int64(len(data)) > maxBytes { return nil, ErrMessageTooLarge }
	return []byte(strings.ReplaceAll(string(data), "\r\n", "\n")), nil
}
