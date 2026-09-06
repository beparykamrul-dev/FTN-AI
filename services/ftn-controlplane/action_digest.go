package controlplane

import("crypto/sha256";"encoding/hex";"strings")

func ActionDigest(tenantID,serverID,action string,services []string)string{h:=sha256.New();write:=func(v string){h.Write([]byte(strings.TrimSpace(v)));h.Write([]byte{0})};write(tenantID);write(serverID);write(action);for _,s:=range NormalizeServices(services){write(s)};return hex.EncodeToString(h.Sum(nil))}
