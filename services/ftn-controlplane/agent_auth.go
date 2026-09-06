package controlplane
import("crypto/sha256";"crypto/subtle";"encoding/hex";"strings")
type AgentIdentity struct{ServerID string;Fingerprint string;Enrolled bool;Revoked bool}
func AuthorizeAgent(expected,presented AgentIdentity)bool{es:=strings.TrimSpace(expected.ServerID);ps:=strings.TrimSpace(presented.ServerID);ef:=strings.ToLower(strings.TrimSpace(expected.Fingerprint));pf:=strings.ToLower(strings.TrimSpace(presented.Fingerprint));if es==""||ps==""||es!=ps{return false};if !expected.Enrolled||!presented.Enrolled||expected.Revoked||presented.Revoked||ef==""||pf==""||len(ef)!=64||len(pf)!=64{return false};eb,err:=hex.DecodeString(ef);if err!=nil{return false};pb,err:=hex.DecodeString(pf);if err!=nil{return false};return subtle.ConstantTimeCompare(eb,pb)==1}
func Fingerprint(identityMaterial []byte)string{sum:=sha256.Sum256(identityMaterial);return hex.EncodeToString(sum[:])}
