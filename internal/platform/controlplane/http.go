package controlplane

import "encoding/json"

type HealthResponse struct{Status string `json:"status"`;API string `json:"api"`}
func HealthJSON()[]byte{b,err:=json.Marshal(HealthResponse{Status:"ok",API:"v1"});if err!=nil{return []byte(`{"status":"ok","api":"v1"}`)};return b}
type EventEnvelope struct{Type string `json:"type"`;Data any `json:"data,omitempty"`}
func EncodeEvent(e EventEnvelope)([]byte,error){return json.Marshal(e)}
