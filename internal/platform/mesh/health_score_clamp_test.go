package mesh

import (
 "testing"
 "time"
)

func TestHealthObserveClampsScore(t *testing.T) {
 r:=NewHealthRegistry(DefaultHealthPolicy()); p:=r.Observe("peer",time.Now().UTC(),255)
 if p.Score!=100 {t.Fatalf("score=%d",p.Score)}
}
