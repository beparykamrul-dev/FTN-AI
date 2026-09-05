package proxy

import("strings";"time")
type SecurityPolicy struct{MaxConnectionsPerClient int;MaxRequestsPerWindow int;RateWindow time.Duration;RequireTLS bool;RejectPrivateTargets bool;MaxHeaderBytes int}
func DefaultSecurityPolicy()SecurityPolicy{return SecurityPolicy{MaxConnectionsPerClient:64,MaxRequestsPerWindow:600,RateWindow:time.Minute,RequireTLS:true,RejectPrivateTargets:true,MaxHeaderBytes:32<<10}}
func(p SecurityPolicy)Valid()bool{return p.MaxConnectionsPerClient>0&&p.MaxRequestsPerWindow>0&&p.RateWindow>0&&p.MaxHeaderBytes>0}
func(p SecurityPolicy)ValidateRequest(scheme string,connections,requests,headerBytes int,privateTarget bool)bool{scheme=strings.ToLower(strings.TrimSpace(scheme));if !p.Valid()||connections<0||requests<0||headerBytes<0{return false};if p.RequireTLS&&scheme!="https"{return false};if connections>p.MaxConnectionsPerClient||requests>p.MaxRequestsPerWindow||headerBytes>p.MaxHeaderBytes{return false};if p.RejectPrivateTargets&&privateTarget{return false};return true}
