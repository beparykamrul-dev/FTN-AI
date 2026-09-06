package controlplane

import("strings";"github.com/google/uuid")
func CorrelationID(existing string)string{if v:=strings.TrimSpace(existing);v!=""{return v};return uuid.NewString()}
