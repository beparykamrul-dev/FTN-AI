package controlplane

import (
 "testing"
 "time"
)

func TestAgentSessionValidAtExpiryIsFalse(t *testing.T) {
 issued:=time.Date(2026,1,1,0,0,0,0,time.UTC); s:=AgentSession{ServerID:"s1",Fingerprint:"fp",IssuedAt:issued,ExpiresAt:issued.Add(time.Minute)}
 if s.Valid(issued.Add(time.Minute)) {t.Fatal("expired session accepted")}
}
