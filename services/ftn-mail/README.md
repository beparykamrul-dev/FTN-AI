# FTN Mail Server

Private, lightweight, self-hosted mail infrastructure for the owner, family, and Family Time Network business under `familytimenet.com`.

## Service boundary

- SMTP submission/transport boundary
- IMAP mailbox boundary
- Mailbox/account provisioning
- TLS and authentication integration points
- Mail event publishing to FTN Account
- DKIM/SPF/DMARC deployment configuration
- Rate limiting, abuse controls and audit hooks
- No public self-service mailbox signup
- No open relay

The mail transport remains independent from the FTN Account Engine. Financial interpretation is performed by the account/event pipeline, never by SMTP/IMAP transport code. The mail service has no direct financial posting permission.

## Ownership scopes

- `personal` — owner's private mail
- `family` — explicitly provisioned family mailboxes
- `business` — Family Time Network operational mail

## Example service mailboxes

- `billing@familytimenet.com`
- `accounts@familytimenet.com`
- `support@familytimenet.com`
- `noreply@familytimenet.com`

Personal and family addresses are provisioned only after explicit owner selection.

## Deployment

This directory contains FTN-owned configuration and integration code. Production DNS, TLS certificates, DKIM private keys, mailbox passwords, and service tokens must remain in the deployment secret store and must never be committed to Git.
