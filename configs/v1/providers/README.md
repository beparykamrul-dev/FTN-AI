# FTN Provider CT Catalog

Provider integrations are represented as provider-specific container/agent templates (CTs).

## Lifecycle

`discover -> normalize -> validate -> store -> health-check -> policy-gate -> agent-ready`

The catalog stores metadata and source references only. Credentials are secret-store references and must never be committed.

## Safety

- Provider records do not grant routing authority.
- Global ASN/IP data is observational until validated against authoritative routing data.
- Privileged mutations remain FTN Control Panel approval-gated.
- Experimental or security-sensitive capabilities are `restricted` by default.
- Open-source source repositories are provenance references; FTN does not automatically execute arbitrary repository code.
