# FTN SMS Queue Transaction Boundary V1

This boundary defines the production transaction around SMS enqueueing. It is an integration contract, not a mock implementation.

## Required atomic operation

The API layer must authorize the request before creating a queue record. The database transaction must then persist the SMS message and its initial queue state atomically.

Conceptually:

1. Begin PostgreSQL transaction.
2. Re-check the authorization-sensitive Sender ID state inside the service boundary.
3. Insert `sms_messages` with `QUEUED` or `PENDING_APPROVAL`.
4. Insert the corresponding audit/delivery event when required by the persistence model.
5. Commit.
6. Return the server-generated message ID only after commit succeeds.

If the transaction rolls back, the API must not report the message as queued.

## Worker ownership

Workers claim only committed queue records. Claiming must prevent two workers from concurrently processing the same message. State transitions are persisted before the transport result is acknowledged to higher layers.

## Failure semantics

- Database failure: do not submit to transport.
- Authorization failure: do not create a queue record.
- Policy/Sender ID failure: do not create a queue record.
- Temporary transport failure: persist retryable state and retry according to queue policy.
- Permanent transport failure: persist terminal failure.

## Security

Transport adapters never receive unauthenticated API requests directly. The API, IAM, policy, queue, and audit boundaries remain authoritative.
