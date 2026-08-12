package transport

import (
	"context"
	"errors"
	"strings"
)

type DeliveryStore interface {
	Deliver(ctx context.Context, identityID, localPart, domain, sender string, recipients []string, raw []byte) error
}

type InboundDelivery struct {
	Auth  MailboxAuthorizer
	Store DeliveryStore
}

func (d InboundDelivery) Deliver(ctx context.Context, identityID, sender string, recipients []string, raw []byte) error {
	if d.Auth == nil || d.Store == nil || identityID == "" || sender == "" || len(raw) == 0 {
		return errors.New("invalid delivery request")
	}
	for _, recipient := range recipients {
		parts := strings.SplitN(strings.TrimSpace(recipient), "@", 2)
		if len(parts) != 2 || parts[0] == "" || !strings.EqualFold(parts[1], "familytimenet.com") {
			return errors.New("recipient not allowed")
		}
		if err := d.Auth.AuthorizeRecipient(ctx, identityID, recipient); err != nil {
			return err
		}
		if err := d.Store.Deliver(ctx, identityID, strings.ToLower(parts[0]), "familytimenet.com", sender, []string{recipient}, raw); err != nil {
			return err
		}
	}
	return nil
}
