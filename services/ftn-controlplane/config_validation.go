package controlplane

type ControlPlaneConfig struct{LeaseTTLSeconds int;SessionTTLSeconds int;MaxPageSize int}
func(c ControlPlaneConfig)Valid()bool{return c.LeaseTTLSeconds>=10&&c.LeaseTTLSeconds<=86400&&c.SessionTTLSeconds>=10&&c.SessionTTLSeconds<=86400&&c.MaxPageSize>=1&&c.MaxPageSize<=1000}
