-- FTN SMS queue claim primitives.
-- Requires services/ftn-sms/storage/schema.sql to be applied first.

-- Claim one eligible message atomically. The caller must run this statement
-- inside a transaction and keep the row lock only for the claim operation.
WITH candidate AS (
    SELECT id
    FROM sms_messages
    WHERE status = 'QUEUED'
      AND scheduled_at <= NOW()
      AND attempts < 8
    ORDER BY priority DESC, scheduled_at ASC, created_at ASC
    FOR UPDATE SKIP LOCKED
    LIMIT 1
)
UPDATE sms_messages m
SET status = 'PROCESSING',
    attempts = attempts + 1,
    updated_at = NOW()
FROM candidate c
WHERE m.id = c.id
RETURNING m.*;

-- Successful transport acceptance.
UPDATE sms_messages
SET status = 'SENT',
    sent_at = COALESCE(sent_at, NOW()),
    updated_at = NOW()
WHERE id = $1
  AND status = 'PROCESSING';

-- Delivery confirmation.
UPDATE sms_messages
SET status = 'DELIVERED',
    delivered_at = COALESCE(delivered_at, NOW()),
    updated_at = NOW()
WHERE id = $1
  AND status IN ('SENT', 'PROCESSING');

-- Temporary failure: return to queue with a caller-computed bounded retry time.
UPDATE sms_messages
SET status = 'QUEUED',
    scheduled_at = $2,
    last_error_class = 'TEMPORARY',
    last_error = $3,
    updated_at = NOW()
WHERE id = $1
  AND status = 'PROCESSING'
  AND attempts < 8;

-- Permanent failure or exhausted retries.
UPDATE sms_messages
SET status = 'FAILED',
    last_error_class = $2,
    last_error = $3,
    failed_at = COALESCE(failed_at, NOW()),
    updated_at = NOW()
WHERE id = $1
  AND status = 'PROCESSING';
