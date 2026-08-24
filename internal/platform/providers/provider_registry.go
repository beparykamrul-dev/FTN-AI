package providers

import "sort"

type Capability string
const (
 CapabilityDNS Capability = "dns"
 CapabilityCDN Capability = "cdn"
 CapabilityEdge Capability = "edge"
 CapabilityCache Capability = "cache"
 CapabilityHosting Capability = "hosting"
 CapabilityCloud Capability = "cloud"
 CapabilityTraffic Capability = "traffic"
 CapabilityObservability Capability = "observability"
)

type Provider struct {
 ID string `json:"id"`
 Name string `json:"name"`
 Capabilities []Capability `json:"capabilities"`
 Endpoint string `json:"endpoint,omitempty"`
 Enabled bool `json:"enabled"`
 Metadata map[string]string `json:"metadata,omitempty"`
}

type Registry struct { providers map[string]Provider }
func NewRegistry() *Registry { return &Registry{providers: map[string]Provider{}} }
func (r *Registry) Upsert(p Provider) { if r.providers==nil { r.providers=map[string]Provider{} }; r.providers[p.ID]=p }
func (r *Registry) Get(id string) (Provider,bool) { p,ok:=r.providers[id]; return p,ok }
func (r *Registry) List() []Provider { out:=make([]Provider,0,len(r.providers)); for _,p:=range r.providers { out=append(out,p) }; sort.Slice(out,func(i,j int)bool{return out[i].ID<out[j].ID}); return out }
func (r *Registry) Supports(id string, c Capability) bool { p,ok:=r.providers[id]; if !ok || !p.Enabled{return false}; for _,x:=range p.Capabilities {if x==c{return true}}; return false }
