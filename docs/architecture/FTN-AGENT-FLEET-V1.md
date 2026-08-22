# FTN Agent Fleet V1

FTN uses one logical agent architecture for service agents, user agents, developer assistants, and the FTN personal assistant. The runtime is provider-neutral so intelligence can remain on FTN-controlled infrastructure.

## Scopes

- `service`: scoped to one FTN service and its operational context.
- `user`: scoped to one user/tenant and permitted personal context.
- `developer`: scoped to engineering repositories, code context, tests, and documentation.
- `assistant`: general FTN assistant mode with the same policy and memory boundaries.

## Isolation

Every request carries `service_id`, `user_id`, `tenant_id`, and role where applicable. Memory, tools, permissions, and audit records must be resolved from this scope before execution.

## Control flow

```text
Request
  -> Scope Resolver
  -> Agent Fleet
  -> Local/Approved Runtime
  -> Retrieval + Memory
  -> Policy / Approval Gate
  -> Tool Execution (only when permitted)
  -> Audit Event
  -> Response
```

## Non-negotiable policy

AI may analyze, explain, plan, and suggest without changing external state. Side-effecting operations require an explicit authorization/approval path. The agent must not bypass FTN IAM, service policy, audit logging, or tenant isolation.

## Developer assistant

The developer mode uses the same fleet but receives only repository/workspace context that the developer is authorized to access. It can explain code, trace dependencies, propose patches, run approved checks, and prepare changes; publishing or destructive operations remain behind approval.

## Service agents

Each service receives a logical agent identity, for example `service:ftn-mail`, `service:ftn-sms`, or `service:ftn-dns`. Service agents share common runtime infrastructure but have service-specific tools and policies.

## User agents

Each user receives a logically isolated agent context. User memory is tenant-scoped and must never be mixed with another user's context.
