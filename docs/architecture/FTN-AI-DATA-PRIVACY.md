# FTN AI Data Privacy Boundary

FTN AI uses scoped context. Customer, developer, organization and network contexts are isolated by identity, tenant and policy.

- minimum necessary context is passed to an AI layer
- sensitive records remain behind FTN authorization
- external/specialist AI layers receive only policy-approved context
- prompts and responses are not automatically retained as audit data
- data retention is configurable by FTN policy
- customer memory is isolated from developer and network memory

```text
FTN Data
  ↓
IAM + Tenant Policy
  ↓
Minimum Context
  ↓
Selected AI Layer
  ↓
Validated Response
```
