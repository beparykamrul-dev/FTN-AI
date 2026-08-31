# FTN ISP Multi-Admin Panel Model

FTN ISP supports multiple independently scoped admin panels over one shared control plane.

## Panels

- Super Admin — global governance, security, tenancy and approvals
- ISP Admin — customers, services, packages and operational configuration
- Billing Admin — invoices, ledger, payments and reconciliation
- NOC Admin — devices, topology, monitoring and incidents
- Engineer Panel — assigned devices, diagnostics and approved operations
- Support Admin — tickets, customer support and AI call-center queues
- Reseller Panel — reseller customers, packages, commissions and reports
- Partner Panel — explicitly entitled partner resources
- Auditor Panel — read-only audit and compliance views
- AI Operations Panel — AI recommendations, incidents and approval queue

## Isolation

Panels share APIs and data services but permissions are enforced server-side. A panel is not a security boundary by itself. Every request is evaluated against identity, role, permission, tenant/account scope and resource ownership.

## Privileged actions

Device configuration, recovery execution, credential operations and other privileged mutations require the appropriate permission and approval policy. UI controls never bypass server-side authorization.

## Suggested hostnames

- admin.familytimenet.com
- billing.familytimenet.com
- noc.familytimenet.com
- engineer.familytimenet.com
- support.familytimenet.com
- reseller.familytimenet.com
- partner.familytimenet.com
- audit.familytimenet.com
- aiops.familytimenet.com

All hostnames should terminate at the authenticated API/control plane rather than exposing database or infrastructure credentials.
