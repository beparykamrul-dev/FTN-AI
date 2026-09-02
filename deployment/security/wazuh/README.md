# FTN Wazuh integration

Wazuh is the FTN security telemetry and detection layer. The integration is intentionally **manual-by-default** for operational changes.

## Modes

- Automatic: observe, correlate, alert, recommend, backup, and bounded health recovery.
- Approval required: configuration, firewall, routing, credential, and service-disable changes.
- Explicit approval: deletion and destructive recovery.

Credentials are supplied through the deployment secret store; they are never committed to Git.

## Deployment boundary

The repository provides the FTN integration contract and deployment boundary. The Wazuh manager/indexer/dashboard deployment is only enabled on FTN-owned infrastructure and must be validated in the target environment before being marked live.

## Runtime checks

After Wazuh is installed on the target node, validate:

1. manager/API health;
2. enrolled FTN agents;
3. TLS certificate validation;
4. alert ingestion into FTN;
5. control-plane correlation;
6. audit persistence;
7. approval gate enforcement;
8. restart/recovery behavior.
