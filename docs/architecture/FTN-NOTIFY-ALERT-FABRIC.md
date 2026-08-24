# FTN Unified Notification & Alert Fabric

## Goal
A provider-neutral, self-hosted notification layer for FTN, local users, ISPs, and corporate customers.

## Channels
- FTN Bot (native primary channel)
- Web/PWA push
- Android push
- Email
- SMS adapter
- Telegram adapter
- WhatsApp adapter
- IMO adapter
- Webhook/API
- WebSocket/SSE realtime alerts

## FTN Bot
FTN Bot is the preferred native experience rather than a clone of Telegram/WhatsApp.

Capabilities:
- one-to-one and group alerts
- service status cards
- incident timeline
- acknowledgement and escalation
- alert severity and deduplication
- maintenance notices
- customer-specific service alerts
- ISP/corporate tenant alert rooms
- local-language notifications
- offline queue and retry
- device/session management
- signed event metadata

## Alert Pipeline
Event -> normalize -> deduplicate -> correlate -> severity -> tenant policy -> channel routing -> delivery -> acknowledgement -> escalation -> audit.

## Monitoring integration
Network, DNS, BGP, flow, eBPF, security, media, game, provider, billing and infrastructure events can emit into the same alert fabric without making notification a control-plane dependency.

## Local/free service model
FTN can expose a free local alert service for FTN users and optional tenant alerting for ISP/corporate customers. Delivery quotas and resource policies are isolated per tenant so customer alert traffic cannot starve FTN NOC alerts.

## Reliability
- local durable queue
- retry with exponential backoff
- dead-letter queue
- delivery receipts
- channel health scoring
- fallback channel selection
- rate limiting
- per-tenant quotas
- outage-aware suppression
- idempotency keys
- audit trail

## Security
- mTLS service-to-service
- signed webhooks
- short-lived credentials
- tenant isolation
- encrypted secrets
- minimal notification payloads
- no provider is required for core FTN alerts

## Architecture
FTN Monitoring -> Alert/Event Bus -> FTN Alert Engine -> FTN Bot Gateway -> Local delivery / optional external adapters.

External messaging providers remain replaceable adapters; FTN Bot remains the native primary interface.