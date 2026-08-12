package provisioning

import (
	"context"
	"errors"
)

type MailboxProvisioner interface {
	ProvisionMailbox(ctx context.Context, identityID, localPart, domain string) error
}

type ServiceAssignmentStore interface {
	AssignService(ctx context.Context, identityID, serviceID, provisionedBy string) error
}

var ErrInvalidMailboxRequest = errors.New("invalid mailbox provisioning request")

// ProvisionFTNMail assigns the FTN Mail service and provisions the mailbox.
// The Control Panel is the only caller allowed to perform this operation.
func ProvisionFTNMail(ctx context.Context, assignments ServiceAssignmentStore, mail MailboxProvisioner, identityID, localPart, provisionedBy string) error {
	if identityID == "" || localPart == "" || provisionedBy == "" {
		return ErrInvalidMailboxRequest
	}
	if err := assignments.AssignService(ctx, identityID, "mail", provisionedBy); err != nil {
		return err
	}
	return mail.ProvisionMailbox(ctx, identityID, localPart, "familytimenet.com")
}
