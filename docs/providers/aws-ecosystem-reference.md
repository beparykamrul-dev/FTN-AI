# AWS Ecosystem Reference for FTN-AI

FTN-AI treats AWS public/open-source repositories as capability references and optional provider adapters, not as a hard runtime dependency.

## Source families
- AWS
- aws-ia
- awslabs
- aws-samples
- awsdocs
- AWS-Devops-Projects

## Capability mapping
- Infrastructure as Code: CDK / CloudFormation patterns
- Terraform and infrastructure modules
- Kubernetes/EKS deployment patterns
- Cloud networking and multi-region architecture
- Edge-to-cloud operations
- Observability and telemetry
- DevOps automation and CI/CD
- Security automation and IAM-aware workflows
- Private CA / certificate automation patterns
- Object storage and media delivery patterns
- AI/agent operations and runbook automation
- MCP-based infrastructure tools
- Fault-injection and resilience testing
- Cost/usage analysis

## FTN integration boundary
AWS-specific code is not copied into the FTN core blindly. FTN uses provider adapters, normalized capability schemas, and provider-specific credentials/configuration. Public OSS can inform implementation where license-compatible; official APIs/credentials can be attached later without changing the core control plane.

## FTN flow
Provider OSS/API -> Provider Adapter -> Normalizer -> FTN Database -> Metrics/Telemetry -> AI Agent -> Control Plane

## Operational goals
- provider-neutral architecture
- live latency/health/traffic measurements where authorized/observable
- cost-aware routing
- cache-aware delivery
- edge/cloud interoperability
- automatic failover
- auditable provider usage
