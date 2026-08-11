# FTN Provider Admission Checklist

Use this checklist before promoting an integration from research to production.

- [ ] Upstream identity verified
- [ ] License reviewed
- [ ] Security review completed
- [ ] Capability matrix completed
- [ ] Provider contract tests pass
- [ ] Source-wrapper tests pass, if applicable
- [ ] Import is read-only and non-mutating
- [ ] Snapshot is deterministic
- [ ] Health probe validated
- [ ] Latency probe validated
- [ ] Timeout/cancellation behavior validated
- [ ] Authentication failures classified
- [ ] Retry and idempotency behavior validated
- [ ] Drift detection validated
- [ ] Audit/approval boundary validated
- [ ] Compatibility range recorded
- [ ] Conformance report approved
- [ ] Operational owner recorded

## Final gate

Only after every mandatory item is complete may the provider move to `admitted`. Admission does not bypass the runtime mutation controls.
