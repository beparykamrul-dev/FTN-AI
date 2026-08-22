# FTN AI Network Operations

Network AI is intentionally separated from general customer/project AI because network operations are telemetry- and safety-sensitive.

It should consume authorized network signals such as device health, interfaces, routing state, service dependencies, alerts, flow metrics and approved topology data.

Capabilities:
- anomaly detection and correlation
- incident triage
- dependency/failure-path analysis
- capacity and utilization analysis
- service-impact estimation
- maintenance recommendations
- NOC summaries and escalation preparation

Safety boundary:
- read/diagnose by default
- configuration changes require explicit authorization
- destructive or disruptive operations require approval
- every operational action is audited
- network credentials remain server-side
