package inventory

import("fmt";"net/url";"strings")
type NetBoxObject struct{ID string `json:"id"`;Type string `json:"type"`;Name string `json:"name"`;Address string `json:"address,omitempty"`;Site string `json:"site,omitempty"`;Metadata map[string]string `json:"metadata,omitempty"`}
type NetBoxClient struct{BaseURL string;TokenRef string}
func(c NetBoxClient)Validate()error{base:=strings.TrimSpace(c.BaseURL);if base==""{return fmt.Errorf("NetBox base URL is required")};u,err:=url.Parse(base);if err!=nil||u.Scheme!="https"||u.Host==""{return fmt.Errorf("NetBox base URL must be a valid HTTPS URL")};if strings.TrimSpace(c.TokenRef)==""{return fmt.Errorf("NetBox token reference is required")};return nil}
func(o NetBoxObject)Valid()bool{return strings.TrimSpace(o.ID)!=""&&strings.TrimSpace(o.Type)!=""&&strings.TrimSpace(o.Name)!=""}
func(o NetBoxObject)Normalize()NetBoxObject{o.ID=strings.TrimSpace(o.ID);o.Type=strings.ToLower(strings.TrimSpace(o.Type));o.Name=strings.TrimSpace(o.Name);o.Address=strings.TrimSpace(o.Address);o.Site=strings.TrimSpace(o.Site);return o}
