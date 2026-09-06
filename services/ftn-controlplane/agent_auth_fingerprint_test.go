package controlplane

import "testing"

func TestFingerprintDeterministic(t *testing.T) {
 a:=Fingerprint([]byte("server-public-material")); b:=Fingerprint([]byte("server-public-material")); if a=="" || a!=b {t.Fatalf("fingerprints differ: %q %q",a,b)}
}
