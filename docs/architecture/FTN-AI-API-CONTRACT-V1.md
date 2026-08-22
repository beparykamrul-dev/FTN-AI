# FTN AI API Contract V1

Canonical public base URL: `https://api.familytimenet.com`

## Endpoints

- `GET /v1/health` — public health/version information.
- `POST /v1/agent/chat` — authenticated web/mobile assistant request.
- `GET /v1/agent/usage` — authenticated quota and usage information.
- `GET /v1/agent/summary` — authenticated important summary for the active scope.
- `GET /v1/agent/details` — authenticated drill-down information permitted by scope.
- `POST /v1/service-request` — authenticated customer service-request submission.
- `GET /v1/services` — public service discovery.

## Chat request

```json
{
  "category": "customer",
  "message": "Why is my service unavailable?",
  "client": "android"
}
```

The server resolves the authenticated account, tenant, entitlement, scope, policy, quota and agent category. Clients never submit privileged scope identifiers as authority.

## Chat response

```json
{
  "request_id": "...",
  "answer": "...",
  "summary": [],
  "details_available": true,
  "usage": {
    "requests_remaining": 19,
    "tokens_remaining": 19000
  },
  "approval_required": false
}
```

## Rules

1. HTTPS only.
2. Anonymous access is limited to explicitly public endpoints.
3. Customer and organization data are tenant/user scoped.
4. Usage is checked before model execution.
5. Side-effecting tools require authorization and, where policy requires, explicit approval.
6. Every request receives an auditable request ID.
7. API responses distinguish factual service data from AI recommendations.
8. No secret is embedded in public web or mobile clients.

## Runtime strategy

The API is model-provider neutral. It should use the lightest approved runtime capable of satisfying the request, escalating only when policy/configuration permits.
