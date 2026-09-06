package controlplane

import "strings"
func ValidTransport(v string)bool{switch strings.ToLower(strings.TrimSpace(v)){case "agent","ssh","api":return true;default:return false}}
