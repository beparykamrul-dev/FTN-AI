package global

import "testing"
func TestSameAuthoritySetNormalizesNames(t *testing.T){if !SameAuthoritySet([]string{"NS1.Example.com.","ns2.example.com"},[]string{"ns2.example.com.","ns1.example.com."}){t.Fatal("expected normalized authority sets to match")}}
func TestSameAuthoritySetRejectsDifferentSets(t *testing.T){if SameAuthoritySet([]string{"ns1.example.com"},[]string{"ns2.example.com"}){t.Fatal("expected authority sets to differ")}}
