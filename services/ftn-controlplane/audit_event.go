package controlplane

import("strings";"time")

type AuditEvent struct{ID string;TenantID string;ActorID string;Action string;TargetID string;Outcome string;At time.Time}
func(a AuditEvent)Valid()bool{return strings.TrimSpace(a.ID)!=""&&strings.TrimSpace(a.TenantID)!=""&&strings.TrimSpace(a.ActorID)!=""&&strings.TrimSpace(a.Action)!=""&&strings.TrimSpace(a.TargetID)!=""&&strings.TrimSpace(a.Outcome)!=""&&!a.At.IsZero()}
