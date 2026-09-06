package controlplane

import "testing"

func TestFingerprintIsDeterministicAndDistinct(t *testing.T) {
 a:=Fingerprint([]byte("a")); if a!=Fingerprint([]byte("a")) {t.Fatal("fingerprint not deterministic")}; if a==Fingerprint([]byte("b")) {t.Fatal("fingerprint collision for distinct inputs")}
}
