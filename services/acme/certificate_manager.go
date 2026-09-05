package acme

import("context";"fmt";"strings";"time")
type Certificate struct{ID string;Subject string;Issuer string;ExpiresAt time.Time;Status string}
type Store interface{List(context.Context)([]Certificate,error);Get(context.Context,string)(Certificate,error);Save(context.Context,Certificate)error}
func Validate(c Certificate)error{if strings.TrimSpace(c.ID)==""||strings.TrimSpace(c.Subject)==""||strings.TrimSpace(c.Issuer)==""{return fmt.Errorf("certificate id, subject and issuer are required")};if c.ExpiresAt.IsZero(){return fmt.Errorf("certificate expiry is required")};return nil}
func IsExpiringSoon(c Certificate,now time.Time,within time.Duration)bool{if Validate(c)!=nil||within<0{return false};return !c.ExpiresAt.Before(now)&&c.ExpiresAt.Sub(now)<=within}
