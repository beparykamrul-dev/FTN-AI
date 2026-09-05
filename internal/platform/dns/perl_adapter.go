package dns

import("fmt";"strings")
type PerlAdapter struct{Name string;Version string}
func NewPerlAdapter(version string)*PerlAdapter{return &PerlAdapter{Name:"perl-dns",Version:strings.TrimSpace(version)}}
func(p *PerlAdapter)Validate()error{if p==nil{return fmt.Errorf("perl adapter is required")};if strings.TrimSpace(p.Name)==""{return fmt.Errorf("perl adapter name is required")};if strings.TrimSpace(p.Version)==""{return fmt.Errorf("perl version is required")};return nil}
