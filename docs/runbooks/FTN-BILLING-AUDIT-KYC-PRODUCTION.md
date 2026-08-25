# FTN Billing Audit + KYC Production Runbook

## Source of truth

PostgreSQL is authoritative for invoices, payments, KYC status, employee identity, audit events, and evidence. AI is advisory and must never mutate the financial ledger.

## Audit questions

1. Which customers have outstanding balances?
2. Which invoices have verified payments, pending verification, or no recorded payment?
3. Which payment records have no entry timestamp?
4. Which employee collected a payment and which employee entered it?
5. Which collections have collector/entry mismatches?
6. What evidence supports a missing-entry cause?
7. Which employees have historical collection activity, and what role/status applied at collection time?
8. Which KYC records require manual review?
9. Which anomalies require management review?

## Missing-entry classification

The system distinguishes recorded facts from hypotheses. A payment with evidence of a customer payment but no ledger entry is classified as `customer_payment_not_entered`. Evidence of an employee collection task without an entry is `employee_collection_not_entered`. Queue/system evidence produces `system_failure`. Duplicate/validation evidence produces `duplicate_or_invalid`. Otherwise the case remains `cause_unknown_requires_review`.

Absence of a ledger payment is never treated as proof that a customer did not pay; cash, gateway, employee-session, device-session, and other evidence must be reconciled first.

## Employee accountability

For each employee, the database-backed audit surface reports collection count, collected amount, verified collected amount, unentered amount/count, and collector-vs-entry mismatch count. Role history is evaluated at the transaction timestamp rather than using the employee's current role.

## KYC + AI

KYC uses database records, provenance, evidence references, duplicate-identity checks, deterministic rules, and AI advisory review. AI returns a finding, confidence, explanation, evidence references, and recommended next action. Automated approval/rejection is only allowed when an explicit policy permits it; otherwise manual review is required.

## Evidence / legal review

Potential misconduct is represented as evidence-backed audit findings and relationship signals, not automatic criminal or syndicate conclusions. Authorized reviewers can export case packages containing the finding, source references, timestamps, hashes, provenance, and legal-hold state. Original evidence is preserved.

## Recommended operational workflow

`collect -> verify -> reconcile -> explain -> review -> authorize -> export/resolve`

Do not overwrite original payment or audit evidence to repair a finding. Corrections must be append-oriented, authorized, reason-coded, and auditable.

## UI

Management receives balance/aging, missing-entry, employee accountability, reconciliation, KYC, case, and evidence views. Employees receive only their authorized collection/entry workload. Customers receive their own balance, invoice history, payment history, and verification status.

## Privacy and access

Use RBAC and least privilege. Sensitive identity fields are masked in ordinary views. Legal/compliance exports require an authorized role. Private social-account discovery and unauthorized tracking are prohibited.
