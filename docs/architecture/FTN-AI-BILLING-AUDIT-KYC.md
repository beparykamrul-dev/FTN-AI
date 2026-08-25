# FTN AI Billing Audit + Customer Verification

## Purpose

Production control layer for customer billing completeness, collection attribution, employee accountability, KYC/verification status, and reconciliation.

## Required outcomes

- Determine the expected bill for every active customer and billing period.
- Determine whether an invoice/bill entry exists.
- Determine whether payment exists and how much was received.
- Attribute every collection to a verified employee/reseller/channel and collection method.
- Detect bills that are missing, late, duplicated, edited, reversed, or inconsistent.
- Distinguish evidence-backed causes of a missing bill from unknown causes.
- Prevent former employees or unauthorized collectors from creating new collection records.
- Preserve historical records; corrections are append-only adjustments, never silent overwrites.

## State model

`EXPECTED -> INVOICED -> COLLECTED -> RECONCILED`

Exception states:

`MISSING_INVOICE`, `OVERDUE`, `PARTIAL`, `DUPLICATE`, `UNMATCHED_PAYMENT`, `REVERSED`, `DISPUTED`, `CLOSED`

## Missing-entry investigation

The system must not guess why a bill is missing. It evaluates evidence:

1. Customer-side evidence: payment intent, receipt, submitted payment reference, gateway callback, cash/collection confirmation.
2. Employee-side evidence: assigned collection task, collection session, device/session audit, attempted entry, submitted receipt, correction/rejection event.
3. System evidence: invoice generation job, API request, validation error, outage, queue failure, duplicate detection.
4. If evidence is insufficient: `UNKNOWN_REQUIRES_REVIEW`.

AI may summarize the evidence and rank likely causes, but the authoritative classification remains deterministic and auditable.

## Employee collection accountability

Every collection record requires:

- collector_id;
- role at collection time;
- customer_id;
- invoice_id/payment_intent_id;
- amount and currency;
- timestamp;
- collection method;
- reference/receipt;
- source device/session;
- approval/reconciliation state;
- immutable audit event.

Employment status is evaluated at transaction time. A former employee may retain historical visibility according to policy, but cannot create or modify new collection records unless explicitly re-authorized.

## KYC / verification

KYC is service-specific and follows the existing FTN requirements engine. Existing verified attributes may be reused only where lawful and authorized. Verification adapters must expose evidence, status, timestamp, verification level and correction/review path.

AI verification is advisory unless a deterministic verification policy explicitly permits an automated status transition. Sensitive identity data is minimized, encrypted, access-controlled and excluded from telemetry.

## AI audit functions

AI can:

- explain outstanding balances;
- detect unusual collection patterns;
- identify missing-entry clusters;
- compare employee/customer collection patterns;
- summarize reconciliation exceptions;
- prioritize cases for human review;
- explain why a record is flagged using available evidence.

AI cannot silently:

- mark a payment as received;
- erase or rewrite ledger history;
- accuse an employee or customer of wrongdoing;
- change employment permissions;
- issue refunds or financial commitments.

## Dashboard

### Customer view

- current due;
- overdue amount;
- invoice history;
- payment history;
- verification status;
- disputed/unmatched payments;
- receipts.

### Employee view

- assigned collection amount;
- collected amount;
- pending entry amount;
- reconciled amount;
- exceptions;
- correction/reversal history;
- active/inactive status.

### Management view

- total expected;
- invoiced;
- collected;
- outstanding;
- missing invoice entries;
- unmatched payments;
- employee/reseller collection attribution;
- reconciliation gap;
- aging;
- AI-prioritized exceptions.

## Financial integrity

The authoritative ledger is append-oriented. Every payment callback is authenticated, idempotent and reconciled. Manual adjustments require reason, actor, authorization and audit reference.

This module integrates with the existing FTN Billing Gateway, KYC/requirements engine, Control Plane, Audit service and notification layer. fileciteturn723file0 fileciteturn724file0
