# FTN Diagnostic Engine — Source Preview

## Purpose

Preview of the production module boundary for Incident → Evidence → Dependency → Root Cause → AI Advisory.

## Flow

```text
Telemetry
  -> Incident Detector
  -> Evidence Correlator
  -> Dependency Graph
  -> Candidate Causes
  -> Evidence Validation
  -> Confidence
  -> Incident Record
  -> AI Advisory
  -> Approval / Policy
  -> Remediation
  -> Verification
```

## Package layout

```text
backend/internal/diagnostics/
├── model/
│   ├── incident.go
│   ├── evidence.go
│   ├── dependency.go
│   └── diagnosis.go
├── correlate/
│   └── correlator.go
├── graph/
│   └── dependency_graph.go
├── rootcause/
│   └── analyzer.go
├── ai/
│   └── advisor.go
├── policy/
│   └── policy.go
└── service.go
```

## Source preview

```go
package diagnostics

type Incident struct {
    ID          string
    Severity    string
    NodeID      string
    MAC         string
    IP          string
    Interface   string
    Service     string
    PathID      string
    EvidenceIDs []string
}

type Diagnosis struct {
    IncidentID string
    Cause      string
    Confidence float64
    Evidence   []string
    Impact     []string
}

func (s *Service) Diagnose(i Incident) Diagnosis {
    evidence := s.correlator.Collect(i)
    graph := s.graph.Resolve(i)
    candidates := s.rootCause.Find(evidence, graph)
    return s.rootCause.Validate(candidates, evidence)
}
```

## Safety boundary

```text
AI = advisory
Privileged action = approval required
External telemetry export = disabled
Secrets in telemetry = prohibited
Destructive migration = disabled
```

The preview is intentionally a source-level architecture preview; implementation must follow the repository's existing Go module conventions and tests before production activation.
