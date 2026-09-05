package gis

import("fmt";"strings")
type LuaAdapter struct{Name string `json:"name"`;Version string `json:"version"`}
func NewLuaAdapter(version string)*LuaAdapter{return &LuaAdapter{Name:"lua-gis",Version:strings.TrimSpace(version)}}
func(a LuaAdapter)Validate()error{if strings.TrimSpace(a.Name)!="lua-gis"{return fmt.Errorf("invalid lua adapter name")};version:=strings.TrimSpace(a.Version);if version==""{return fmt.Errorf("lua version is required")};if len(version)>128||strings.ContainsAny(version,"\r\n\x00")||strings.ContainsAny(version," "){return fmt.Errorf("invalid lua version")};return nil}
func(a LuaAdapter)Ready()bool{return a.Validate()==nil}
