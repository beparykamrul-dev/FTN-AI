# FTN AI Category Fleet

FTN uses one lightweight orchestration layer with specialized AI roles instead of deploying a heavyweight model for every feature.

## Categories

| Agent | Purpose | Primary audience |
|---|---|---|
| Studio AI | Preview, UI review, content and implementation assistance | users + developers |
| Call Center AI | Support conversations, intent detection, ticket summaries and routing | customers + support staff |
| Billing Gateway AI | Billing explanations, payment-flow assistance, anomaly summaries and reconciliation support | customers + billing staff |
| Network/Service AI | Service health, incidents, dependency context and operational summaries | service operators |
| Developer AI | Code understanding, debugging, tests, documentation and patch planning | developers |
| Customer AI | General FTN assistant with user-scoped context | customers |
| Executive Summary AI | Important summaries and drill-down reports | administrators |

## Lightweight-first strategy

1. Rules, deterministic validators and retrieval are attempted first.
2. Small/local models handle classification, extraction, routing and short summaries.
3. A larger approved model is used only when the task needs deeper reasoning.
4. Tool execution remains behind IAM and approval policy.

The category layer does not require one model per agent. Multiple categories can share the same runtime while receiving different tools, prompts, memory namespaces and policies.

## Control-panel presentation

Every category exposes two views:

- **Summary:** important events, decisions, usage, failures and recommended next steps.
- **Details:** evidence, timeline, request/tool traces, affected service/user scope and supporting data.

The control panel should default to Summary and allow authorized operators to drill into Details.

## Studio preview agent

Studio AI receives a preview URL/build artifact and produces UI/UX, content, accessibility, functional and error summaries. It must not publish or deploy changes without explicit authorization.

## Call center agent

Call Center AI classifies intent, gathers approved account context, suggests responses, creates/updates support tickets when authorized, and produces concise call summaries. Sensitive customer data remains tenant-scoped.

## Billing gateway agent

Billing Gateway AI explains invoices/payments, detects suspicious or inconsistent billing states, summarizes gateway failures, and assists reconciliation. It does not independently authorize payments, refunds, or financial state changes.

## Human control

All high-impact actions—financial changes, account changes, infrastructure operations, publication/deployment, or customer-data export—require explicit authorization and are audited.
