# FTN AI Public Service

## Product

FTN AI is one branded AI assistant available to public users and FTN customers through `https://api.familytimenet.com` and FTN web/mobile clients.

## Access model

- Public users can use the free tier without needing an FTN customer account for explicitly public capabilities.
- Authenticated FTN customers receive a user-scoped agent and customer plan entitlement.
- Paid plans increase the daily request/token allowance after payment is confirmed by the billing system.
- Developer access is separately scoped and must never inherit a customer's private context.

## Entitlement flow

```text
Public / Customer
      -> api.familytimenet.com
      -> Authentication (when required)
      -> User/Tenant scope
      -> Billing entitlement
      -> AI quota gate
      -> Agent Fleet
      -> Model / Retrieval / Tools
```

Payment changes the **entitlement**, not the agent's security policy. A paid user does not automatically receive privileged infrastructure tools. Side-effecting actions still require authorization and approval.

## Plans

Initial service defaults are defined in `internal/platform/agent/quota.go` and are configuration, not a promise of final commercial pricing. Billing must be the source of truth for active entitlement and expiry.

## Public-service requirements

- Mobile and web compatible HTTPS API.
- Rate limiting and abuse protection.
- Per-user/tenant isolation.
- Daily request and token quotas.
- Maximum input/output limits.
- Usage accounting and audit events.
- Graceful quota-exceeded response with upgrade guidance.
- No API secrets embedded in client applications.
- No exposure of internal NOC, customer records, or privileged tools to anonymous users.
