package transport

import (
	"context"
	"errors"
	"strings"
)

type MailboxAuthorizer interface {
	AuthorizeRecipient(ctx context.Context, identityID, recipient string) error
}

type MailboxLookup interface {
	MailboxForIdentity(ctx context.Context, identityID, localPart, domain string) (string, error)
}

type Authorizer struct{ Lookup MailboxLookup }

func (a Authorizer) AuthorizeRecipient(ctx context.Context, identityID, recipient string) error {
	if a.Lookup == nil || identityID == "" { return errors.New("mailbox authorization failed") }
	parts := strings.SplitN(strings.TrimSpace(recipient), "@", 2)
	if len(parts) != 2 || parts[0] == "" || !strings.EqualFold(parts[1], "familytimenet.com") {
		return errors.New("recipient not allowed")
	}
	_, err := a.Lookup.MailboxForIdentity(ctx, identityID, strings.ToLower(parts[0]), "familytimenet.com")
	return err
}
