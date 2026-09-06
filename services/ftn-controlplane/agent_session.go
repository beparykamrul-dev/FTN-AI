package controlplane

import("strings";"time")

type AgentSession struct{ServerID string;Fingerprint string;IssuedAt time.Time;ExpiresAt time.Time;Revoked bool}
func(s AgentSession)Valid(now time.Time)bool{if strings.TrimSpace(s.ServerID)==""||strings.TrimSpace(s.Fingerprint)==""||s.Revoked||s.IssuedAt.IsZero()||s.ExpiresAt.IsZero()||!s.ExpiresAt.After(s.IssuedAt){return false};if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()};return !now.Before(s.IssuedAt)&&now.Before(s.ExpiresAt)}
