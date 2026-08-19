# FTN Agent Public API

Canonical public endpoint: `https://api.familytimenet.com`

The agent API is designed for browser and mobile clients. Public availability means the HTTPS API endpoint is reachable from the Internet; it does **not** mean anonymous access to private data or privileged tools.

## Public capabilities

- Agent chat for approved public/user scopes
- Service-specific assistant entry points
- User-scoped assistant sessions
- Developer assistant endpoints for authorized developer clients
- Health/version discovery

## Security boundary

```text
Mobile / Web
    |
    | HTTPS
    v
api.familytimenet.com
    |
    +-- rate limit / abuse protection
    +-- CORS / origin policy
    +-- authentication
    +-- tenant + user scope
    +-- agent policy
    +-- approval gate
    +-- tools
    +-- audit
```

Unauthenticated requests may only access explicitly public capabilities such as health, version and public service discovery. User conversations and service actions require authentication and scoped authorization.

## Client compatibility

The API contract is transport-neutral so the same endpoints can be consumed by Android, iOS, mobile web, desktop web, and future FTN applications.

## Operational requirements

- TLS certificate managed automatically through FTN ACME infrastructure.
- No API key or secret embedded in public frontend/mobile binaries.
- Per-user/tenant rate limits.
- Request and tool audit trail without storing unnecessary sensitive content.
- Explicit approval required for side-effecting actions.
