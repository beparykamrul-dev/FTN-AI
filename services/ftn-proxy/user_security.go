package proxy

// UserSecurityPolicy is a baseline for protecting end-user sessions and
// sensitive financial/e-commerce traffic. It deliberately stores no payment
// credentials or raw authentication secrets.
type UserSecurityPolicy struct {
	RequireTLS            bool
	RequireSecureCookies  bool
	RequireRequestID      bool
	MaxSessionAgeSeconds  int64
	MaxClockSkewSeconds   int64
	RejectReplay          bool
	RequireIdempotencyKey bool
}

func DefaultUserSecurityPolicy() UserSecurityPolicy {
	return UserSecurityPolicy{
		RequireTLS: true,
		RequireSecureCookies: true,
		RequireRequestID: true,
		MaxSessionAgeSeconds: 3600,
		MaxClockSkewSeconds: 120,
		RejectReplay: true,
		RequireIdempotencyKey: true,
	}
}

// ValidateUserRequest checks transport/session invariants. Identity and
// payment authorization remain the responsibility of the identity/payment
// service, not the reverse proxy.
func (p UserSecurityPolicy) ValidateUserRequest(scheme string, secureCookie, requestID, replayDetected, idempotencyKeyPresent bool) bool {
	if p.RequireTLS && scheme != "https" { return false }
	if p.RequireSecureCookies && !secureCookie { return false }
	if p.RequireRequestID && !requestID { return false }
	if p.RejectReplay && replayDetected { return false }
	if p.RequireIdempotencyKey && !idempotencyKeyPresent { return false }
	return true
}
