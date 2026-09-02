package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "net"
    "strings"
)

// DNSGuardDataplaneRequest is the minimal resolver input. It intentionally
// carries no query payload beyond the normalized domain and category.
type DNSGuardDataplaneRequest struct {
    TenantID string `json:"tenant_id"`
    SubjectType string `json:"subject_type"`
    SubjectRef string `json:"subject_ref"`
    Domain string `json:"domain"`
    Category string `json:"category,omitempty"`
}

type DNSGuardDataplaneResponse struct {
    Decision string `json:"decision"`
    Reason string `json:"reason"`
    Category string `json:"category"`
    DomainHash string `json:"domain_hash"`
}

// NormalizeDNSName provides deterministic DNS-name canonicalization without
// performing any network access.
func NormalizeDNSName(domain string) string {
    d := strings.TrimSpace(strings.ToLower(domain))
    d = strings.TrimSuffix(d, ".")
    if d == "" || len(d) > 253 { return "" }
    if net.ParseIP(d) != nil { return "" }
    labels := strings.Split(d, ".")
    if len(labels) < 2 { return "" }
    for _, label := range labels {
        if label == "" || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' { return "" }
        for _, r := range label {
            if !(r == '-' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9') { return "" }
        }
    }
    return d
}

func DNSGuardDomainHash(domain string) string {
    d := NormalizeDNSName(domain)
    if d == "" { return "" }
    sum := sha256.Sum256([]byte(d))
    return hex.EncodeToString(sum[:])
}

// CompileDNSGuardDataplane is side-effect-free and suitable for adapters used
// by dnsdist/CoreDNS/Unbound. The actual resolver adapter must call this after
// obtaining the active tenant/profile policy.
func CompileDNSGuardDataplane(p DNSGuardProfile, req DNSGuardDataplaneRequest) DNSGuardDataplaneResponse {
    d := NormalizeDNSName(req.Domain)
    if d == "" { return DNSGuardDataplaneResponse{"allow", "invalid_or_non_domain", "unknown", ""} }
    decision := CompileDNSGuardDecision(p, req.Category, d)
    return DNSGuardDataplaneResponse{decision.Decision, decision.Reason, decision.Category, DNSGuardDomainHash(d)}
}

func EncodeDNSGuardResponse(v DNSGuardDataplaneResponse) []byte { b, _ := json.Marshal(v); return b }
