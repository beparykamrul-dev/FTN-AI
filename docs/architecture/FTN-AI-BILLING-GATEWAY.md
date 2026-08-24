# FTN-AI Billing Gateway & Payment Agent

## Purpose

A provider-neutral billing and payment orchestration layer for FTN ISP, e-commerce, corporate services, subscriptions and approved marketplace services.

## Architecture

```text
Customer / Business / Reseller
            |
      FTN Account / Wallet
            |
      Billing Gateway API
            |
    Payment Orchestrator
      /      |       \
     /       |        \
 Local   Bank/MFS   Card/Other
 adapters  adapters  adapters
            |
     Provider Webhooks
            |
   Signature + Idempotency
            |
      Ledger / Reconciliation
       /       |        \
 PostgreSQL  Redis   ClickHouse
            |
       Invoice / Receipt
            |
       FTN AI Payment Agent
            |
   Explain / Notify / Assist
```

## Core capabilities

- customer account and service billing;
- subscription and recurring-billing state;
- invoices, receipts and credits;
- payment-intent lifecycle;
- provider adapter abstraction;
- webhook verification and idempotency;
- reconciliation and settlement records;
- refunds/voids as policy-controlled workflows;
- reseller/corporate billing boundaries;
- e-commerce order-payment linkage;
- payment status notifications through FTN Notify;
- accounting/export APIs;
- fraud/risk signal adapters;
- audit trail and immutable transaction references.

## Payment-agent boundary

The AI agent is an assistant/orchestrator, not an unrestricted financial actor. It may explain invoices, identify the correct service/order, show available payment options and report payment status. Sensitive actions such as payment execution, refund, beneficiary/account changes or financial commitment require explicit customer confirmation and applicable authorization/policy gates.

The platform must never store raw payment secrets when a compliant provider tokenization flow is available. Provider credentials remain in a secrets manager and are isolated per adapter/tenant.

## Provider-neutral design

FTN does not hard-code one payment provider. Adapters can be added for licensed/authorized banks, mobile financial services, cards, gateways and other legally supported payment rails. Each adapter implements a common contract for intent creation, status lookup, webhook verification, reconciliation and supported refund/void operations.

## Ledger rules

The authoritative transaction ledger is append-oriented and auditable. Payment-provider callbacks are treated as untrusted input until authenticated and validated. Duplicate callbacks must be idempotent. Financial state changes require deterministic transaction IDs and reconciliation support.

## Security

- TLS/mTLS for internal services;
- RBAC and MFA for privileged operations;
- tenant isolation;
- secrets manager integration;
- signed webhook verification;
- idempotency keys;
- encryption at rest;
- immutable audit events;
- rate limiting and abuse controls;
- least-privilege provider credentials;
- backup and restore testing.

## Deployment

The gateway is an independent commerce/billing plane. POPs may provide local UI, cache and notification functions, while authoritative billing/ledger state remains on approved durable backend nodes. The FTN main control plane is not required to be the transit path for ordinary payment or commerce traffic.

## Compliance gate

Before enabling any real payment rail, FTN must verify the applicable Bangladesh regulatory, licensing, KYC/AML, tax, consumer-protection, data-retention and provider-contract requirements. A development/test adapter must never be presented as a production payment service.
