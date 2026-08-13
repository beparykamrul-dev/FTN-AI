# FTN SMS Worker Contract V1

The worker consumes only committed `QUEUED` records. Claiming must use the atomic PostgreSQL queue primitive and must tolerate multiple workers.

## Processing

1. Claim one eligible message.
2. Re-check message state and delivery policy.
3. Resolve the configured authorized transport.
4. Submit through the transport interface.
5. Persist the normalized result in the same service storage model.
6. Record a delivery event.

## Result handling

- `SENT`: transition to `SENT`.
- `DELIVERED`: transition to `DELIVERED` when a trusted delivery event is received.
- `TEMPORARY_FAILURE`: increment attempts and requeue with bounded exponential backoff and jitter when retryable.
- `PERMANENT_FAILURE`: transition to `FAILED`.
- retry budget exhausted: transition to `FAILED`.

## Concurrency

Workers must never process the same claimed message concurrently. Claim ownership and lease/recovery semantics must be persisted; a worker restart must not permanently strand a message in `PROCESSING`.

## Safety boundaries

The worker cannot bypass IAM, Sender ID approval, rate limits, AI approval, or transport authorization. SMS body and recipient values must not appear in logs or metric labels.

## Shutdown

On graceful shutdown, stop claiming new work, finish bounded in-flight operations, and leave recoverable work in a persisted state for the next worker.
