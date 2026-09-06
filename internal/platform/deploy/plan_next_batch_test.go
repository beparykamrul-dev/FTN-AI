package deploy

import "testing"

func TestPlanValidationRejectsMissingIdentity(t *testing.T) {
	if (Plan{}).Valid() { t.Fatal("empty deploy plan must be invalid") }
}
