# FTN Provider Conformance Report

A conformance report records the evidence used to move a provider through the FTN registry lifecycle.

## Required evidence

- upstream identity verified
- license reviewed
- capability matrix completed
- provider contract tests passed
- source-wrapper tests passed where applicable
- security review completed
- compatibility range recorded
- health and latency probes validated
- mutation boundary verified

## Decision states

```text
PASS -> ADMIT
FAIL -> BLOCK
REVIEW -> HOLD
```

A failed or incomplete conformance report must not be promoted to production admission. The report is evidence for the registry; it does not itself grant mutation authority.

## Auditability

Each report should retain a provider version, adapter version, test run identifier, review timestamp, and reviewer/automation identity so future releases can be compared against the admitted baseline.
