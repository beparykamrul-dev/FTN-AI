# FTN SMS Observability V1

## Health

- `/health/live` reports process liveness only.
- `/health/ready` reports whether the service can safely accept traffic.
- `/health/deps` reports normalized dependency state and latency.

## Prometheus metrics

```text
ftn_sms_service_up
ftn_sms_service_ready
ftn_sms_queue_depth
ftn_sms_queue_oldest_age_seconds
ftn_sms_messages_total{status}
ftn_sms_delivery_events_total{event_type}
ftn_sms_transport_up{transport}
ftn_sms_transport_latency_ms{transport}
ftn_sms_transport_errors_total{transport,error_class}
ftn_sms_retry_total{reason}
ftn_sms_rate_limit_total{scope}
ftn_sms_approval_pending
```

## Label policy

Do not use recipient, message body, identity ID, account ID, phone number, message ID, or Sender ID as Prometheus labels. These values remain in application storage/audit according to retention policy.

## NOC mapping

- service up + ready + transport available: `HEALTHY`
- service up but queue/transport degraded: `DEGRADED`
- service unavailable or no usable transport: `CRITICAL`

## Alerts

Critical conditions include service not ready for a sustained window and all configured transports unavailable. Warning conditions include sustained queue growth, retry spikes, elevated transport latency, and delivery failure rate above configured policy.

No automatic AI remediation is implied by an alert. AI may analyze and recommend; privileged remediation remains subject to FTN IAM/approval policy.
