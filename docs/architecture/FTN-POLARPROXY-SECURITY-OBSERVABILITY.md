# FTN PolarProxy Security Observability

## Decision

PolarProxy is treated as an optional, isolated security-observability adapter, not as the FTN customer traffic proxy.

## Intended scope

- Controlled FTN-owned infrastructure and test/security segments only.
- TLS inspection only where FTN explicitly controls the endpoint/CA and the traffic is authorized for inspection.
- Decrypted session output may feed controlled PCAP/IDS analysis pipelines.
- Keep customer/private traffic outside this inspection path by default.

## Integration boundary

```text
FTN Edge / Security Segment
        |
        v
   PolarProxy Adapter
        |
   +----+----+
   |         |
  PCAP     Metadata
   |         |
   v         v
IDS/Forensics  OTel/Flow Fabric
   |             |
   +------+------+
          v
      ClickHouse
          |
        AIOps
```

## FTN-native controls

- mTLS and FTN PKI remain the trust boundary.
- Secrets and private keys remain server-side.
- No transparent inspection of customer traffic by default.
- No decrypted payload retention unless explicitly enabled for an authorized security workflow.
- Apply retention limits, access controls, audit events and data minimization.
- PolarProxy failure must never become a mandatory dependency for normal FTN traffic.

## Provider-neutral design

PolarProxy is replaceable through the FTN security-observability adapter contract. The core FTN data plane, telemetry schema and AIOps pipeline do not depend on it.

## References

PolarProxy is a TLS/SSL inspection proxy capable of forward, reverse, termination and inline modes and can output decrypted traffic as PCAP for security analysis. The current product is proprietary, so FTN should use it only as an optional adapter rather than treating it as an open-source core dependency.
