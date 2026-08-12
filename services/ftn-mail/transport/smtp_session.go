package transport

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
)

type SessionConfig struct { Hostname string; TLSConfig *tls.Config; RequireTLS bool; MaxMessageSize int64 }
type Authenticator interface { Authenticate(context.Context, string, string) (string, error) }
type Delivery interface { Deliver(context.Context, string, string, []string, []byte) error }

// ServeSession handles the SMTP command boundary for one connection.
// AUTH is delegated to an explicit adapter; unauthenticated submission is rejected.
func ServeSession(ctx context.Context, conn net.Conn, cfg SessionConfig, auth Authenticator, delivery Delivery) error {
	if conn == nil || cfg.MaxMessageSize <= 0 || delivery == nil { return errors.New("invalid SMTP session configuration") }
	defer conn.Close()
	br := bufio.NewReader(conn); bw := bufio.NewWriter(conn)
	write := func(code int, text string) error { if _, err := fmt.Fprintf(bw, "%d %s\r\n", code, text); err != nil { return err }; return bw.Flush() }
	if err := write(220, cfg.Hostname+" FTN Mail Service Ready"); err != nil { return err }

	tlsActive, authenticated := false, false
	identityID, mailFrom := "", ""
	var recipients []string
	for {
		select { case <-ctx.Done(): return ctx.Err(); default: }
		line, err := br.ReadString('\n'); if err != nil { return err }
		cmdLine := strings.TrimSpace(line); parts := strings.Fields(cmdLine)
		if len(parts) == 0 { if err := write(500, "Invalid command"); err != nil { return err }; continue }
		cmd := strings.ToUpper(parts[0])
		switch cmd {
		case "EHLO", "HELO":
			if err := write(250, cfg.Hostname+"\r\n250-AUTH PLAIN LOGIN\r\n250 STARTTLS"); err != nil { return err }
		case "STARTTLS":
			if tlsActive || cfg.TLSConfig == nil { if err := write(454, "TLS unavailable"); err != nil { return err }; continue }
			if err := write(220, "Ready to start TLS"); err != nil { return err }
			tlsConn := tls.Server(conn, cfg.TLSConfig); if err := tlsConn.HandshakeContext(ctx); err != nil { return err }
			conn = tlsConn; br = bufio.NewReader(conn); bw = bufio.NewWriter(conn); tlsActive = true
		case "AUTH":
			if auth == nil { if err := write(454, "Authentication unavailable"); err != nil { return err }; continue }
			// Mechanism-specific AUTH exchange is intentionally implemented by the adapter.
			if err := write(502, "AUTH adapter required"); err != nil { return err }
		case "MAIL":
			if cfg.RequireTLS && !tlsActive { if err := write(530, "Must issue STARTTLS first"); err != nil { return err }; continue }
			if !authenticated { if err := write(530, "Authentication required"); err != nil { return err }; continue }
			arg := strings.TrimSpace(cmdLine[len(parts[0]):]); if len(arg) < 5 || !strings.EqualFold(arg[:5], "FROM:") { if err := write(501, "Invalid sender"); err != nil { return err }; continue }
			mailFrom = strings.TrimSpace(arg[5:]); if mailFrom == "" { if err := write(501, "Invalid sender"); err != nil { return err }; continue }
			recipients = recipients[:0]; if err := write(250, "Sender OK"); err != nil { return err }
		case "RCPT":
			if mailFrom == "" { if err := write(503, "Need MAIL FROM first"); err != nil { return err }; continue }
			arg := strings.TrimSpace(cmdLine[len(parts[0]):]); if len(arg) < 3 || !strings.EqualFold(arg[:3], "TO:") { if err := write(501, "Invalid recipient"); err != nil { return err }; continue }
			rcpt := strings.TrimSpace(arg[3:]); if rcpt == "" { if err := write(501, "Invalid recipient"); err != nil { return err }; continue }
			recipients = append(recipients, rcpt); if err := write(250, "Recipient OK"); err != nil { return err }
		case "DATA":
			if mailFrom == "" || len(recipients) == 0 { if err := write(503, "Need sender and recipient"); err != nil { return err }; continue }
			if err := write(354, "End data with <CRLF>.<CRLF>"); err != nil { return err }
			raw, err := ReadMessage(br, cfg.MaxMessageSize); if err != nil { if err2 := write(552, "Message too large or invalid"); err2 != nil { return err2 }; continue }
			if err := delivery.Deliver(ctx, identityID, mailFrom, recipients, raw); err != nil { if err2 := write(451, "Temporary delivery failure"); err2 != nil { return err2 }; continue }
			if err := write(250, "Message accepted"); err != nil { return err }; mailFrom = ""; recipients = recipients[:0]
		case "RSET": mailFrom = ""; recipients = recipients[:0]; if err := write(250, "Reset OK"); err != nil { return err }
		case "NOOP": if err := write(250, "OK"); err != nil { return err }
		case "QUIT": _ = write(221, "Bye"); return nil
		default: if err := write(502, "Command not implemented"); err != nil { return err }
		}
	}
}
