# FTN SMS Transport V1

Common transport boundary for self-hosted SMS delivery.

## Interface contract

The queue worker submits an authorized message to a transport adapter and receives a normalized result:

- `SENT` — accepted by the transport
- `TEMPORARY_FAILURE` — safe to retry
- `PERMANENT_FAILURE` — do not retry

Adapters must also expose readiness/health state and a stable transport identifier for metrics.

## Adapters

### GSM modem

FTN-controlled modem/SIM equipment. The adapter owns serial/USB communication, modem state, reconnect handling, and normalized delivery submission. It must not contain IAM or business-policy logic.

### Authorized SMPP

Direct connection to an operator-authorized SMPP endpoint. Credentials are loaded from the secret provider/environment and never committed to the repository. The adapter owns bind/rebind, submit_sm response normalization, and transport-level errors.

## Security

- Only approved Sender IDs reach the adapter.
- No sender spoofing or operator-control bypass.
- Secrets are never persisted in source control or logged.
- Recipient/body data must not be emitted as metric labels.
- The adapter cannot bypass queue policy, IAM, approval, or rate limits.

## Privacy

No third-party cloud SMS API is required. External traffic is limited to the configured GSM operator path or an authorized SMPP operator connection.
