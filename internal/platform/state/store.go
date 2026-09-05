package state

import("context";"errors";"strings")
var ErrNotConfigured=errors.New("state store is not configured")
type Decision struct{ServiceID string;PeerID string;Transport string;PolicyVersion string;Status string;Reason string}
type Store interface{SaveDecision(context.Context,Decision)error}
func ValidateDecision(d Decision)error{d.ServiceID=strings.TrimSpace(d.ServiceID);d.PeerID=strings.TrimSpace(d.PeerID);d.Transport=strings.ToLower(strings.TrimSpace(d.Transport));d.PolicyVersion=strings.TrimSpace(d.PolicyVersion);d.Status=strings.ToLower(strings.TrimSpace(d.Status));d.Reason=strings.TrimSpace(d.Reason);if d.ServiceID==""||d.PolicyVersion==""{return errors.New("service_id and policy_version are required")};if d.Status==""{return errors.New("decision status is required")};if d.Transport==""{return errors.New("decision transport is required")};return nil}
