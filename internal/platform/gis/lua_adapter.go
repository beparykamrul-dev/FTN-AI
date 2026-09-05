package gis

import("fmt";"strings")
type LuaAdapter struct{Name string `json:"name"`;Version string `json:"version"`}
func NewLuaAdapter(version string)*LuaAdapter{return &LuaAdapter{Name:"lua-gis",Version:strings.TrimSpace(version)}}
func(a LuaAdapter)Validate()error{if strings.TrimSpace(a.Name)==""{return fmt.Errorf("lua adapter name is required")};if strings.TrimSpace(a.Version)==""{return fmt.Errorf("lua version is required")};return nil}
