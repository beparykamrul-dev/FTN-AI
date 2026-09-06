package gis

import "testing"

func TestLuaAdapterRejectsInvalidVersion(t *testing.T) {
	if NewLuaAdapter("bad version").Validate() == nil { t.Fatal("Lua version containing spaces must be rejected") }
}

func TestLuaAdapterReady(t *testing.T) {
	if !NewLuaAdapter("5.4.7").Ready() { t.Fatal("valid Lua adapter must report ready") }
}
