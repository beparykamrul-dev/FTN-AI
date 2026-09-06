package controlplane
import("encoding/base64";"strings")
func EncodeCursor(value string)string{value=strings.TrimSpace(value);if value==""||len(value)>2048{return ""};return base64.RawURLEncoding.EncodeToString([]byte(value))}
func DecodeCursor(cursor string)(string,bool){cursor=strings.TrimSpace(cursor);if cursor==""{return "",true};if len(cursor)>4096{return "",false};b,err:=base64.RawURLEncoding.DecodeString(cursor);if err!=nil||len(b)==0||len(b)>2048{return "",false};return string(b),true}
