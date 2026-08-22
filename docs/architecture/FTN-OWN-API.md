# FTN Own API

`api.familytimenet.com` is the FTN-owned API namespace and public integration boundary.

The API contract, authentication, tenant isolation, quotas, service routing, AI orchestration, billing entitlements, audit and operational policies belong to FTN. External model/provider APIs are optional implementation dependencies, not the public API identity.

## Ownership boundary

```text
FTN Web / Android / Apps / Developer Clients
                    |
                    v
          api.familytimenet.com
                    |
          FTN API Gateway / Auth
                    |
       +------------+-------------+
       |            |             |
    Services       AI          Billing
       |            |             |
       +------------+-------------+
                    |
             FTN Database/Storage
```

## Design rules

- FTN owns the public API contract and versioning.
- Clients never need to know internal service addresses.
- API credentials and service secrets stay server-side.
- Customer data is tenant-scoped.
- Usage is metered by FTN.
- Payment changes FTN entitlements/limits.
- Privileged service operations require IAM/policy/approval.
- Audit records are generated for security-sensitive actions.
- AI provider changes do not require changing client applications.
- Web and mobile clients use the same FTN API surface.

## Canonical endpoint

`https://api.familytimenet.com`

Internal services can evolve behind this boundary without breaking FTN applications.
