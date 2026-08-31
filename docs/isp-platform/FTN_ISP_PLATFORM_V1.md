# FTN ISP Platform V1

## Scope

Separate ISP billing, monitoring, customer application, support and AI platform integrated with FTN-AI through authenticated APIs.

## Service domains

- isp.familytimenet.com — ISP portal
- app.familytimenet.com — application gateway
- noc.familytimenet.com — NOC dashboard
- billing.familytimenet.com — billing operations
- ai.familytimenet.com — AI services
- support.familytimenet.com — support and call-center services
- api.familytimenet.com — versioned API gateway

## Identity and access

Roles: customer, employee, engineer, reseller, partner, support-agent, billing-agent, noc-operator, admin, super-admin.

All privileged operations require authorization and are audited. Secrets are environment-managed and are never stored in repository configuration.

## Customer Android application

- account registration and login
- OTP and session management
- package and subscription status
- usage and service status
- speed test
- invoices, payments, history and receipts
- connection information
- PPPoE/service status
- support tickets and complaints
- AI chat and AI call-center entry point
- package/add-on requests
- notifications
- referral and profile management
- FTN Connect / FTNWAN integration points
- IPTV/FTN TV integration point
- CCTV Cloud integration point

## Billing

Customer -> package -> subscription -> invoice -> payment -> ledger -> receipt.

Payment providers are isolated behind provider interfaces. Webhooks are authenticated, idempotent and audited.

## Monitoring and discovery

Discovery adapters cover MikroTik, OLT, ONU, switches, servers, SNMP, PPPoE, VLAN, IP pools and NetFlow. Discovery produces inventory and topology relationships; it does not automatically perform privileged configuration changes.

Monitoring covers device health, interfaces, traffic, latency, packet loss, PPPoE sessions, ONU/OLT state, optical telemetry when exposed by the vendor, VLAN/IP-pool state, NetFlow, alerts and historical metrics.

## AI

AI consumes authorized customer, billing, support and network context. It may classify, explain, correlate and recommend. Network/device mutations remain behind an explicit approval boundary.

## AI call center

Customer -> identity verification -> account/service context -> AI diagnosis -> resolution or human/engineer escalation.

## API principles

- versioned routes
- request IDs
- structured logs
- centralized errors
- rate limiting
- security headers
- input validation
- audit events
- least-privilege authorization
- fail-closed privileged operations

## Recovery

Health detection -> incident record -> AI diagnosis -> proposed recovery -> approval -> execution -> verification -> rollback on failed verification -> audit.

## Deployment

The ISP platform is separately deployable from FTN-AI while using controlled authenticated integration points. Android and web clients never receive infrastructure credentials.
