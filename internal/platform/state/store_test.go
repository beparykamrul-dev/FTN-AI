package state

import "testing"

func TestValidateDecision(t *testing.T) {
	if err := ValidateDecision(Decision{ServiceID: "dns", PolicyVersion: "v1", Status: "SELECTED"}); err != nil {
		t.Fatalf("valid decision rejected: %v", err)
	}

	cases := []Decision{
		{PolicyVersion: "v1", Status: "SELECTED"},
		{ServiceID: "dns", Status: "SELECTED"},
		{ServiceID: "dns", PolicyVersion: "v1"},
	}
	for _, tc := range cases {
		if err := ValidateDecision(tc); err == nil {
			t.Fatalf("invalid decision accepted: %+v", tc)
		}
	}
}
