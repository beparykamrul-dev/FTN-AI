package controlplane

import("crypto/sha256";"encoding/hex";"strings")

type DeploymentEnvelope struct{Plan DeploymentPlan `json:"plan"`;Digest string `json:"digest"`}

func SealPlan(p DeploymentPlan)DeploymentEnvelope{h:=sha256.New();write:=func(v string){h.Write([]byte(v));h.Write([]byte{0})};write(strings.TrimSpace(p.ServerID));write(strings.TrimSpace(p.Reason));for _,service:=range p.Services{write(strings.TrimSpace(service))};if p.Approved{write("approved")}else{write("rejected")};return DeploymentEnvelope{Plan:p,Digest:hex.EncodeToString(h.Sum(nil))}}
func VerifyEnvelope(e DeploymentEnvelope)bool{if !e.Plan.Approved||len(e.Digest)!=64{return false};if _,err:=hex.DecodeString(e.Digest);err!=nil{return false};return e.Digest==SealPlan(e.Plan).Digest}
