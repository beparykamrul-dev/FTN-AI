package router

import "testing"
func TestValidateRouteRejectsOversizedPrefix(t *testing.T){r:=Route{Prefix:string(make([]byte,65))};if err:=validateRoute(r);err==nil{t.Fatal("expected oversized prefix rejection")}}
func TestValidateRouteRejectsFamilyMismatch(t *testing.T){r:=Route{Prefix:"10.0.0.0/24",NextHop:"2001:db8::1"};if err:=validateRoute(r);err==nil{t.Fatal("expected address family rejection")}}
