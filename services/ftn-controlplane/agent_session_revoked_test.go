package controlplane

import (
 "testing"
 "time"
)

func TestAgentSessionValidRejectsRevoked(t *testing.T) {
 now:=time.Date(2026,1,1,0,0,0,0,time.UTC); s:=AgentSession{ServerID:"s1",Fingerprint:"fp",IssuedAt:now.Add(-time.Minute),ExpiresAt:now.Add(time.Minute),Revoked:true}
 if s.Valid(now) {t.Fatal("revoked session accepted")}
}
