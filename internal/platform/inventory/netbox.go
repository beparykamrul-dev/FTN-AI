package inventory

import("fmt";"net/url";"strings")
type NetBoxObject struct{ID string `json:"id"`;Type string `json:"type"`;Name string `json:"name"`;Address string `json:"address,omitempty"`;Site string `json:"site,omitempty"`;Metadata map[string]string `json:"metadata,omitempty"`}
type NetBoxClient struct{BaseURL string;TokenRef string}
func(c NetBoxClient)Validate()error{base:=strings.TrimSpace(c.BaseURL);if base==""{return fmt.Errorf("NetBox base URL is required")};u,err:=url.Parse(base);if err!=nil||u.Scheme!="https"||u.Host==""||u.User!=nil{return fmt.Errorf("NetBox base URL must be a valid HTTPS URL")};if len(base)>2048||strings.TrimSpace(c.TokenRef)==""||len(strings.TrimSpace(c.TokenRef))>256{return fmt.Errorf("invalid NetBox client configuration")};return nil}
func(o NetBoxObject)Valid()bool{return strings.TrimSpace(o.ID)!=""&&strings.TrimSpace(o.Type)!=""&&strings.TrimSpace(o.Name)!=""&&len(strings.TrimSpace(o.ID))<=256&&len(strings.TrimSpace(o.Type))<=128&&len(strings.TrimSpace(o.Name))<=512}
func(o NetBoxObject)Normalize()NetBoxObject{o.ID=strings.TrimSpace(o.ID);o.Type=strings.ToLower(strings.TrimSpace(o.Type));o.Name=strings.TrimSpace(o.Name);o.Address=strings.TrimSpace(o.Address);o.Site=strings.TrimSpace(o.Site);if o.Metadata!=nil{m:=make(map[string]string,len(o.Metadata));for k,v:=range o.Metadata{if k=strings.TrimSpace(k);k!=""&&len(k)<=128&&len(v)<=1024{m[k]=strings.TrimSpace(v)}};o.Metadata=m};return o}
