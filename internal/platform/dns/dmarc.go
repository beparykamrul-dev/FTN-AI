package dns

import (
	"fmt"
	"strings"
)

type DMARCPolicy string

const (
	DMARCNone       DMARCPolicy = "none"
	DMARCQuarantine DMARCPolicy = "quarantine"
	DMARCReject     DMARCPolicy = "reject"
)

type DMARCRecord struct {
	Domain             string       `json:"domain"`
	Policy             DMARCPolicy  `json:"policy"`
	SubdomainPolicy    DMARCPolicy  `json:"subdomain_policy,omitempty"`
	Percentage         uint8        `json:"percentage"`
	AggregateReportURI []string     `json:"aggregate_report_uri,omitempty"`
	ForensicReportURI  []string     `json:"forensic_report_uri,omitempty"`
}

func (r DMARCRecord) Validate() error {
	if strings.TrimSpace(r.Domain) == "" {
		return fmt.Errorf("domain is required")
	}
	switch r.Policy {
	case DMARCNone, DMARCQuarantine, DMARCReject:
	default:
		return fmt.Errorf("invalid DMARC policy")
	}
	if r.SubdomainPolicy != "" {
		switch r.SubdomainPolicy {
		case DMARCNone, DMARCQuarantine, DMARCReject:
		default:
			return fmt.Errorf("invalid subdomain policy")
		}
	}
	if r.Percentage > 100 {
		return fmt.Errorf("percentage must be between 0 and 100")
	}
	return nil
}

func (r DMARCRecord) TXTValue() (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}
	parts := []string{"v=DMARC1", "p=" + string(r.Policy)}
	if r.SubdomainPolicy != "" {
		parts = append(parts, "sp="+string(r.SubdomainPolicy))
	}
	parts = append(parts, fmt.Sprintf("pct=%d", r.Percentage))
	if len(r.AggregateReportURI) > 0 {
		parts = append(parts, "rua="+strings.Join(r.AggregateReportURI, ","))
	}
	if len(r.ForensicReportURI) > 0 {
		parts = append(parts, "ruf="+strings.Join(r.ForensicReportURI, ","))
	}
	return strings.Join(parts, "; "), nil
}
