# FTN QKD Security Layer

FTN treats Quantum Key Distribution (QKD) as an optional quantum-safe key-establishment layer for future optical/core deployments. It is not a software substitute for physical QKD hardware.

## Architecture

- QKD-Tx / QKD-Rx: external certified QKD modules
- Quantum channel: dedicated optical path
- Classical channel: authenticated control/key-distillation traffic
- Key Manager (KM): receives and manages QKD key material
- FTN control plane: inventory, health, policy and audit only
- Crypto consumers: IPsec/TLS/storage or other approved symmetric-key consumers
- PQC remains enabled as the software-compatible quantum-safe baseline; QKD is complementary.

## Safety boundaries

1. No raw quantum keys are stored in Git, logs, metrics, OpenSearch, or PostgreSQL.
2. QKD hardware credentials are references to an external secret manager, never plaintext.
3. QKD enable/disable and key-consumer changes are approval-gated.
4. Classical-channel authentication is mandatory.
5. A QKD outage must fail closed for consumers explicitly requiring QKD, or fall back only when an explicit policy permits a PQC/hybrid mode.
6. Production interoperability requires vendor/device evidence; this repository provides contracts and policy, not a simulated QKD implementation.

## Standards alignment

The design follows the QKD network concepts in ITU-T Y.3802 and the QKD protocol framework in ITU-T X.1711, including quantum and classical channels, key distillation, key management, authentication and operational monitoring. ETSI QKD work also defines interoperable key-management interfaces.
