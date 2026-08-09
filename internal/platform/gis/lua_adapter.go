package gis

import "fmt"

// LuaAdapter defines a safe integration boundary for Lua-based GIS/network
// extensions. Lua source is treated as an external plugin and must be executed
// in an isolated worker with an explicit capability policy.
type LuaAdapter struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewLuaAdapter(version string) *LuaAdapter {
	return &LuaAdapter{Name: "lua-gis", Version: version}
}

func (a LuaAdapter) Validate() error {
	if a.Version == "" { return fmt.Errorf("lua version is required") }
	return nil
}
