package gis

import "testing"

func TestLuaAdapterValidation(t *testing.T) {
 a := NewLuaAdapter("5.4.7")
 if err := a.Validate(); err != nil || !a.Ready() { t.Fatalf("valid Lua adapter rejected: err=%v ready=%v", err, a.Ready()) }
 if err := NewLuaAdapter("").Validate(); err == nil { t.Fatal("empty Lua version accepted") }
 if err := NewLuaAdapter("5.4\n7").Validate(); err == nil { t.Fatal("control character in Lua version accepted") }
}
