package adapters

import dnsglobal "github.com/beparykamrul-dev/FTN-AI/services/ftn-dns/global"

// DNSSECObservation is the canonical provider-neutral observation shared with
// the global DNS control-plane layer.
type DNSSECObservation = dnsglobal.DNSSECObservation

// DNSSECPolicy defines the minimum trust conditions required before a DNS
// response can be accepted by the FTN DNS control plane.
type DNSSECPolicy struct {
	RequireValidation    bool
	RequireAuthentication bool
	AllowInsecure         bool
}

// DNSSECDecision is the normalized result consumed by routing/selection.
type DNSSECDecision struct {
	Accepted   bool
	Validated  bool
	Authenticated bool
	Reason     string
}

// EvaluateDNSSEC applies policy to adapter-produced DNSSEC observations.
// Cryptographic verification itself must be performed by the DNSSEC adapter;
// this function only evaluates its verified result.
func EvaluateDNSSEC(policy DNSSECPolicy, observation DNSSECObservation) DNSSECDecision {
	if observation.Error != "" {
		return DNSSECDecision{Reason: observation.Error}
	}
	if policy.RequireValidation && !observation.Validated {
		return DNSSECDecision{Validated: observation.Validated, Authenticated: observation.Authenticated, Reason: "DNSSEC validation required"}
	}
	if policy.RequireAuthentication && !observation.Authenticated {
		return DNSSECDecision{Validated: observation.Validated, Authenticated: observation.Authenticated, Reason: "DNSSEC authentication required"}
	}
	if !observation.Validated && !policy.AllowInsecure {
		return DNSSECDecision{Validated: observation.Validated, Authenticated: observation.Authenticated, Reason: "insecure DNS response rejected"}
	}
	return DNSSECDecision{Accepted: true, Validated: observation.Validated, Authenticated: observation.Authenticated, Reason: "accepted"}
}
