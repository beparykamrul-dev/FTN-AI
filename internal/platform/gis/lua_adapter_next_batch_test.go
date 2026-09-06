package gis

import "testing"

func TestLuaAdapterRejectsInvalidVersion(t *testing.T) {
	for _, version := range []string{"", "lua 5.4", "lua\n5.4"} {
		if NewLuaAdapter(version).Ready() { t.Fatalf("version %q must be rejected", version) }
	}
	if !NewLuaAdapter("5.4.6").Ready() { t.Fatal("valid Lua version rejected") }
}
