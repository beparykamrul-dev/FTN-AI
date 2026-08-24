# FTN E-Commerce Fabric

## Purpose

Add an optional, tenant-isolated e-commerce/service-commerce plane to FTN without coupling commerce workloads to the ISP control plane.

## Architecture

`Customer/App -> FTN Commerce Gateway -> Catalog -> Cart -> Order -> Payment Adapter -> Fulfillment -> Notification`

Supporting services:

- Product/catalog management
- Vendor and reseller marketplace
- Inventory and stock reservation
- Cart and checkout
- Order lifecycle and returns
- Coupon/promotion engine
- Customer accounts and addresses
- Local delivery/service-zone integration
- Payment-provider adapters
- Invoice/tax hooks
- Search and recommendations
- Event-driven order processing
- Fraud/risk signals
- Audit and dispute records
- Seller/developer service marketplace

## FTN integration

Commerce remains a separate domain. It can use FTN identity, PKI, notification, observability, billing, database broker, security and local POP/edge capabilities through explicit APIs.

Local-first behavior is preferred when policy and availability permit:

`Customer -> nearest POP -> Commerce API/cache -> appropriate service/database`

The FTN main server is not required to carry ordinary commerce traffic when a healthy POP/edge path is available.

## Payments

Payment integrations are adapter-based. No payment provider is assumed to be the only provider. Sensitive payment credentials and payment data remain isolated from ordinary customer/application data. Production deployment must follow applicable Bangladesh payment, consumer-protection, tax and data-handling requirements and the contractual terms of each provider.

## Database placement

- PostgreSQL: orders, customers, products, inventory and transactional state.
- Redis: short-lived cart/session/cache/locks where appropriate.
- ClickHouse: analytics, events and operational reporting.
- Object storage: product media, documents and exports.

The database broker may select a suitable backend using workload compatibility, health, latency, quota, capacity, residency and policy—not merely free capacity.

## Security

- mTLS between trusted services
- RBAC and tenant isolation
- MFA for privileged commerce administration
- signed webhook verification
- idempotency keys for payment/order operations
- audit trail for price, stock, payment and order changes
- rate limiting and abuse protection
- secrets rotation
- backup/restore testing
- fraud/risk signals without autonomous accusation

## Event-driven model

`OrderCreated -> StockReserved -> PaymentAuthorized -> FulfillmentStarted -> Delivered -> Completed`

Failures use retry, idempotency, dead-letter handling and compensating actions. External side effects require verified provider responses.

## Marketplace

FTN can later expose a local marketplace where approved businesses, resellers and developers publish:

- physical products
- digital products
- ISP-related services
- hosting/cloud services
- local IT services
- software/support packages

Every seller is a separate tenant with explicit permissions and settlement/accounting boundaries.

## AI/AIOps

AI may recommend products, detect operational anomalies, forecast inventory, classify support requests and assist sellers. It must not autonomously perform sensitive payment/refund/settlement actions unless an explicitly authorized policy permits them.

## Production gate

Commerce is production-ready only after real payment-provider sandbox/production validation, webhook verification, inventory consistency tests, backup/restore tests, security review, load testing, failure recovery and applicable legal/compliance review.
