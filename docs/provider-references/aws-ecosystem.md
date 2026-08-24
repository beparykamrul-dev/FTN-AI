# AWS Ecosystem Reference

FTN-AI provider-neutral reference profile for AWS public/open-source ecosystems.

## Sources
- AWS
- AWS Labs
- AWS Samples
- AWS Docs
- AWS IA
- AWS DevOps Projects

## Capability mapping
- Infrastructure as Code: Terraform, AWS CDK, CloudFormation
- Cloud operations: deployment, lifecycle, automation, rollback
- Kubernetes: EKS-oriented patterns and controllers
- Serverless: Lambda/event-driven execution patterns
- Containers: ECS/EKS workload orchestration
- Networking: VPC, routing, load balancing, service connectivity
- Storage: object/block/file storage abstraction
- Databases: relational, NoSQL, cache abstraction
- Observability: metrics, logs, traces, alarms, dashboards
- AI operations: DevOps-agent style investigation and root-cause workflows
- Security: IAM-oriented identity, least privilege, secrets and policy boundaries
- Multi-account/multi-region: federation and cross-environment management patterns
- CI/CD: pipeline, deployment, artifact and release automation
- Resilience: health checks, fault isolation, recovery and controlled rollback

## FTN integration rule
AWS source material is treated as an architectural/open-source reference. FTN keeps its own provider adapter, normalized data model, control plane and service implementations. Provider-specific credentials and official APIs remain isolated from public-source metadata.

## Traffic and service telemetry
Where FTN has authorized access to provider telemetry/API, normalize latency, throughput, request volume, cache performance, health, route state and service usage into the unified FTN traffic model. Public/open-source examples are not treated as permission to access private provider data.

## AI operations
AWS DevOps Agent patterns are useful references for FTN AI Agent investigation: correlate logs, metrics, traces, configuration and service health; identify probable root causes; recommend or execute approved remediation; record audit history.

## Deployment
FTN may implement reusable deployment modules and one-click workflows inspired by public AWS IaC and sample patterns, while retaining FTN branding, service boundaries and provider-neutral interfaces.

## Notes
- Review each repository's license before incorporating source code.
- Prefer independent FTN implementations when only an architectural pattern is needed.
- Archived, experimental or demonstration repositories must not automatically become production dependencies.
