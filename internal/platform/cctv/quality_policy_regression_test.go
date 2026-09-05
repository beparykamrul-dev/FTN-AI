package cctv

import "testing"

func TestNormalizeProfilesCanonicalizesAndDeduplicates(t *testing.T) {
	out := NormalizeProfiles([]string{" 1080P ","1080p","","4K"})
	if len(out) != 2 || out[0] != "1080p" || out[1] != "4k" { t.Fatalf("unexpected profiles: %#v", out) }
}

func TestDefaultQualityPolicyIsValid(t *testing.T) {
	if !DefaultQualityPolicy().Valid() { t.Fatal("default CCTV quality policy must be valid") }
}
