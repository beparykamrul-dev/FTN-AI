package controlplane

import("encoding/base64";"strings")
func EncodeCursor(value string)string{value=strings.TrimSpace(value);if value==""{return ""};return base64.RawURLEncoding.EncodeToString([]byte(value))}
func DecodeCursor(cursor string)(string,bool){cursor=strings.TrimSpace(cursor);if cursor==""{return "",true};b,err:=base64.RawURLEncoding.DecodeString(cursor);if err!=nil||len(b)==0{return "",false};return string(b),true}
