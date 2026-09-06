package controlplane
import "strings"
type AdmissionRequest struct{TenantID string;ServerID string;Action string;Approved bool}
type AdmissionDecision struct{Allowed bool;Reason string}
func Admit(r AdmissionRequest,knownServer,healthy bool)AdmissionDecision{tenant:=strings.TrimSpace(r.TenantID);server:=strings.TrimSpace(r.ServerID);action:=strings.ToLower(strings.TrimSpace(r.Action));if tenant==""||server==""{return AdmissionDecision{Reason:"tenant and server are required"}};if len(tenant)>256||len(server)>256{return AdmissionDecision{Reason:"tenant or server is too large"}};if action==""||len(action)>128{return AdmissionDecision{Reason:"action is required"}};if !knownServer{return AdmissionDecision{Reason:"server is not enrolled"}};if !healthy{return AdmissionDecision{Reason:"server is unhealthy"}};if !r.Approved{return AdmissionDecision{Reason:"explicit approval required"}};return AdmissionDecision{Allowed:true,Reason:"admission accepted"}}
