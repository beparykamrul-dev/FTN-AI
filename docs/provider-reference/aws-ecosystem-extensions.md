# FTN-AI AWS Ecosystem Extensions

Provider-neutral reference mapping for AWS-related open-source ecosystems supplied for FTN-AI architecture research.

## Sources

- terraform-aws-modules
- aws-cloudformation
- aws-solutions
- aws-amplify
- aws-robotics
- aws-actions
- aws-controllers-k8s
- aws-observability
- aws-neuron
- awskrug
- AWSCookbook
- awsfundamentals-hq

## FTN capability mapping

- Infrastructure as Code: reusable Terraform module patterns, version pinning, composable modules.
- Cloud resource lifecycle: CloudFormation resource/module/hook patterns and declarative provisioning.
- Solutions architecture: production deployment patterns, security reference architectures, multi-account and automation patterns.
- Application platform: frontend/mobile/backend integration patterns where relevant.
- Robotics/edge: device and edge workload patterns where applicable.
- CI/CD: GitHub Actions and deployment automation patterns.
- Kubernetes: AWS Controllers for Kubernetes patterns for declarative cloud resources from Kubernetes.
- Observability: OpenTelemetry, Prometheus, Grafana, logs, traces, metrics and alerting patterns.
- AI/accelerated compute: AWS Neuron architecture references for accelerator-aware workloads; use only where FTN has a compatible workload.
- Education/reference: cookbook and fundamentals repositories as implementation guidance, not runtime dependencies.

## FTN integration boundary

AWS projects remain external references/provider adapters. FTN core remains provider-neutral. Any actual provider integration must use the applicable public API, SDK, CLI, or authorized credential. Open-source code is used only when its license permits the intended use.

## Data model

Provider-specific capabilities are normalized into FTN's provider registry, service catalog, traffic/latency telemetry, health model, deployment metadata, and control-plane interfaces. Provider-specific fields remain namespaced so one provider's assumptions do not leak into another provider's core model.

## Operational goals

- low-latency path selection based on measured conditions
- provider-aware service discovery
- cache/edge/compute capability classification
- infrastructure lifecycle automation
- observability normalization
- cost/usage metadata where an authorized provider interface exposes it
- failover and health-aware routing
- future official API/credential/agent integration without changing the FTN core
