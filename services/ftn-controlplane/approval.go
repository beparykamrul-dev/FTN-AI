package controlplane
import("strings";"time")
type Approval struct{ID string;TenantID string;ServerID string;Action string;ApprovedBy string;ApprovedAt time.Time;Revoked bool}
func(a Approval)Valid(now time.Time)bool{if strings.TrimSpace(a.ID)==""||strings.TrimSpace(a.TenantID)==""||strings.TrimSpace(a.ServerID)==""||strings.TrimSpace(a.Action)==""||strings.TrimSpace(a.ApprovedBy)==""||a.Revoked||a.ApprovedAt.IsZero(){return false};if now.IsZero(){return true};return !a.ApprovedAt.After(now)}
func(a Approval)Normalized()Approval{a.ID=strings.TrimSpace(a.ID);a.TenantID=strings.TrimSpace(a.TenantID);a.ServerID=strings.TrimSpace(a.ServerID);a.Action=strings.ToLower(strings.TrimSpace(a.Action));a.ApprovedBy=strings.TrimSpace(a.ApprovedBy);if !a.ApprovedAt.IsZero(){a.ApprovedAt=a.ApprovedAt.UTC()};return a}
