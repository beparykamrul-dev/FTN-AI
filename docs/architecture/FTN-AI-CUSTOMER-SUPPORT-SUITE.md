# FTN AI Customer & Organization Support Suite

FTN AI is organized as one branded assistant platform with specialized, scoped agents.

## Account and customer assistant

Every authenticated customer can have an FTN AI account linked to their tenant/customer identity. The customer-facing assistant can help with:

- service discovery and service usage
- service requests and support
- billing and payment explanations
- account and entitlement information
- important notices and alerts
- safe troubleshooting guidance

The assistant is displayed as a consistent chat/assistant box across eligible customer service pages.

## FTN organization assistant

An FTN organization workspace can use specialized assistants for operational records such as:

- inventory/material records
- technology/equipment records
- service/site records
- work/field tracking
- safety and incident records
- maintenance planning
- operational summaries

The AI may summarize and analyze authorized records. It must not invent stock, work, safety, financial, or operational facts.

## Call Center AI

Call Center AI supports agents and supervisors with:

- caller intent classification
- customer context retrieval within authorization
- live agent assistance
- call summaries
- ticket creation suggestions
- escalation suggestions
- follow-up summaries
- quality/operations summaries

Automated customer-facing actions remain behind the same IAM, policy, and approval boundaries.

## FTN Chat

FTN Chat is the general conversational entry point for web and mobile clients. It routes requests to the appropriate specialized agent and preserves scope and entitlement boundaries.

## Alert Bot

FTN Alert Bot delivers authorized service/customer/organization alerts through configured channels. It supports severity, deduplication, acknowledgement, escalation, and audit records.

## Lightweight-first

All assistants share a lightweight-first architecture:

```text
Rules / intent / retrieval
        -> small local runtime
        -> larger approved runtime only when needed
```

The category agent owns tools, context and policy; the model runtime is shared.

## Customer chart box

Service pages may expose an assistant/chart box containing:

1. important summary
2. current service status
3. usage/entitlement
4. recommended next step
5. expandable details
6. support/service-request action

Sensitive records are only shown after authorization.
