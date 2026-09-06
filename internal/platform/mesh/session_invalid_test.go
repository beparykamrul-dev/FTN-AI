package mesh

import "testing"

func TestSessionRegistryRejectsInvalidIDs(t *testing.T) {
 r:=NewSessionRegistry(); if r.Heartbeat(""){t.Fatal("empty heartbeat accepted")}; if r.Subscribe("",EventType("x")){t.Fatal("empty subscription accepted")}; if s:=r.Open(" "); s.ID!="" {t.Fatal("empty node opened")}
}
