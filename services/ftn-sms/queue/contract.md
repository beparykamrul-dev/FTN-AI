# FTN SMS Queue V1

Production queue contract for the self-hosted FTN SMS service.

## Lifecycle

`QUEUED -> PROCESSING -> SENT -> DELIVERED`

Temporary transport failures return to `QUEUED` with bounded exponential backoff. Permanent failures become `FAILED` and are not retried indefinitely.

## Worker invariants

- Claim messages atomically; one active worker owns a message at a time.
- Respect `scheduled_at`, priority, and per-sender/per-recipient rate limits.
- Validate Account/IAM authorization before enqueueing and before privileged resend operations.
- Only approved/active Sender IDs may be submitted to a gateway adapter.
- Persist every state transition and delivery event.
- Do not log SMS body, authentication secrets, or full recipient identifiers.
- No third-party cloud SMS provider is required.
- Transport adapters are limited to FTN-controlled GSM modem or authorized operator SMPP.

## Retry

Use bounded exponential backoff with jitter. Retry only transport-class temporary failures. Authentication, policy, invalid recipient, invalid Sender ID, and other permanent failures are terminal.

## Adapter boundary

The queue worker depends on a transport interface; GSM modem and authorized SMPP implementations remain separate adapters. Business logic must not contain modem/SMPP-specific code.
