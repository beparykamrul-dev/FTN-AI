package controlplane

import "strings"

type AdmissionRequest struct{TenantID string;ServerID string;Action string;Approved bool}
type AdmissionDecision struct{Allowed bool;Reason string}

func Admit(r AdmissionRequest,knownServer,healthy bool)AdmissionDecision{if strings.TrimSpace(r.TenantID)==""||strings.TrimSpace(r.ServerID)==""{return AdmissionDecision{Reason:"tenant and server are required"}};if strings.TrimSpace(r.Action)==""{return AdmissionDecision{Reason:"action is required"}};if !knownServer{return AdmissionDecision{Reason:"server is not enrolled"}};if !healthy{return AdmissionDecision{Reason:"server is unhealthy"}};if !r.Approved{return AdmissionDecision{Reason:"explicit approval required"}};return AdmissionDecision{Allowed:true,Reason:"admission accepted"}}
