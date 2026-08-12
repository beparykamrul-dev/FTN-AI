# FTN Mail Server

Lightweight self-hosted mail service for `familytimenet.com`.

## Service boundary

- SMTP submission/transport boundary
- IMAP mailbox boundary
- Mailbox/account provisioning
- TLS and authentication integration points
- Mail event publishing to FTN Account
- DKIM/SPF/DMARC deployment configuration
- Rate limiting, abuse controls and audit hooks

The mail transport remains independent from the FTN Account Engine. Financial interpretation is performed by the account/event pipeline, never by SMTP/IMAP transport code.

## Planned mail identities

- `@familytimenet.com`
- configurable service mailboxes such as `billing@familytimenet.com`, `support@familytimenet.com`, and `accounts@familytimenet.com`

## Deployment

This directory contains FTN-owned configuration and integration code. Actual DNS records and TLS certificates must be provisioned for the production host before public mail delivery is enabled.
