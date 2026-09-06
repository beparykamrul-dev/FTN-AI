package automation

import "testing"

func TestAnsiblePlanValidationAndNormalization(t *testing.T) {
 p := AnsiblePlan{Inventory: " inventory ", Playbook: " playbook.yml ", Limit: []string{"b", " a ", "b", ""}}
 if err := p.Validate(); err != nil { t.Fatalf("valid Ansible plan rejected: %v", err) }
 n := p.Normalize(); if n.Inventory != "inventory" || n.Playbook != "playbook.yml" || len(n.Limit) != 2 || n.Limit[0] != "a" || n.Limit[1] != "b" { t.Fatalf("normalized plan=%+v", n) }
 if err := (AnsiblePlan{Inventory: "inv\n", Playbook: "pb"}).Validate(); err == nil { t.Fatal("control character accepted") }
}
