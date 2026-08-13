# FTN-AI V1 Configuration Bundle

Central, safe configuration templates for Account/IAM, AI governance, Mail, DNS, PostgreSQL, and runtime health.

## Rules

- No production passwords, API keys, private keys, or tokens are committed here.
- Copy `.env.example` to a deployment secret store/environment and fill real values there.
- Existing DNS/PostgreSQL schemas are preserved; migrations must be explicitly reviewed.
- PostgreSQL is the authoritative persistence layer; event streams are delivery mechanisms.
- AI privileged actions require authorization and human approval.
- Service adapters must use registered tools; AI must not receive arbitrary infrastructure credentials.

## Files

- `env.example` — runtime environment contract.
- `postgres.yaml` — PostgreSQL pool/readiness policy.
- `ai-governance.yaml` — identity, tool, risk, approval and audit policy.
- `mail.yaml` — Mail storage/streaming policy.
- `dns.yaml` — DNS health/metrics policy.

These are configuration contracts/templates, not secret material and not a replacement for existing service-specific configuration.
