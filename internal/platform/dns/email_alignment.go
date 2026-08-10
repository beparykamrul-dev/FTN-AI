package dns

import "strings"

type AlignmentMode string

const (
	AlignmentRelaxed AlignmentMode = "relaxed"
	AlignmentStrict  AlignmentMode = "strict"
)

type EmailSecurityStatus struct {
	SPFPresent bool `json:"spf_present"`
	DKIMPresent bool `json:"dkim_present"`
	DMARCPresent bool `json:"dmarc_present"`
	SPFAligned bool `json:"spf_aligned"`
	DKIMAligned bool `json:"dkim_aligned"`
	DMARCEnforcing bool `json:"dmarc_enforcing"`
}

func DomainMatches(authDomain, headerDomain string, mode AlignmentMode) bool {
	a := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(authDomain), "."))
	h := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(headerDomain), "."))
	if a == "" || h == "" { return false }
	if mode == AlignmentStrict { return a == h }
	return a == h || strings.HasSuffix(h, "."+a)
}

func EvaluateEmailSecurity(spf, dkim, dmarc bool, spfAligned, dkimAligned bool, dmarcPolicy DMARCPolicy) EmailSecurityStatus {
	return EmailSecurityStatus{
		SPFPresent: spf,
		DKIMPresent: dkim,
		DMARCPresent: dmarc,
		SPFAligned: spfAligned,
		DKIMAligned: dkimAligned,
		DMARCEnforcing: dmarc && dmarcPolicy != DMARCNone,
	}
}
