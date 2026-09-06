package controlplane

import("strings";"time")

type TargetServer struct{ID string `json:"id"`;Name string `json:"name"`;Address string `json:"address"`;Environment string `json:"environment"`;Transport string `json:"transport"`;AgentID string `json:"agent_id,omitempty"`;Enabled bool `json:"enabled"`;LastSeen time.Time `json:"last_seen,omitempty"`}
type DeploymentTarget struct{ServerID string `json:"server_id"`;Path string `json:"path"`;Version string `json:"version"`}
func ValidateTarget(servers []TargetServer,id string)bool{id=strings.TrimSpace(id);if id==""{return false};for _,s:=range servers{if strings.TrimSpace(s.ID)==id&&s.Enabled{return true}};return false}
