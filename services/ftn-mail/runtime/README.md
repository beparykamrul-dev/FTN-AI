# FTN Mail Runtime

The runtime is intentionally private-use and lightweight.

## Scope

- `familytimenet.com` only
- Personal, family, and Family Time Network accounts
- No public self-signup
- No open relay
- Authenticated submission only
- TLS required for authenticated client access
- Mail events can be emitted to FTN Account Intelligence

## Runtime contract

Transport adapters must depend on the mailbox/auth interfaces rather than the accounting database. Account Intelligence consumes sanitized mail events and cannot directly control SMTP/IMAP transport.

Production deployment must provide real TLS certificates, DNS, storage paths, secrets, and network policy through environment/secret management; none are committed to this repository.
