package dns

import (
	"sort"
	"strings"
	"sync"
)

type ProviderKind string
const ( ProviderPowerDNS ProviderKind="powerdns"; ProviderTechnitium ProviderKind="technitium"; ProviderCoreDNS ProviderKind="coredns"; ProviderUnbound ProviderKind="unbound"; ProviderDNSDist ProviderKind="dnsdist"; ProviderGoDNS ProviderKind="godns"; ProviderAnycast ProviderKind="anycast"; ProviderDNSPod ProviderKind="dnspod"; ProviderCloudflare ProviderKind="cloudflare"; ProviderAkamai ProviderKind="akamai"; ProviderFTN ProviderKind="ftn"; ProviderGoBGP ProviderKind="gobgp" )
type Provider struct { ID string `json:"id"`; Kind ProviderKind `json:"kind"`; Name string `json:"name"`; Endpoint string `json:"endpoint,omitempty"`; Region string `json:"region,omitempty"`; Scope string `json:"scope,omitempty"`; Enabled bool `json:"enabled"`; Primary bool `json:"primary"`; Capabilities []string `json:"capabilities,omitempty"`; DNSZone string `json:"dnsZone,omitempty"`; ConfigRef string `json:"configRef,omitempty"` }
type ProviderRegistry struct { mu sync.RWMutex; providers map[string]Provider }
func NewProviderRegistry()*ProviderRegistry{return &ProviderRegistry{providers:make(map[string]Provider)}}
func(r *ProviderRegistry)Upsert(p Provider)bool{if r==nil{return false};p.ID=strings.TrimSpace(p.ID);p.Name=strings.TrimSpace(p.Name);p.Endpoint=strings.TrimRight(strings.TrimSpace(p.Endpoint),"/");p.Region=strings.TrimSpace(p.Region);p.Scope=strings.ToLower(strings.TrimSpace(p.Scope));p.DNSZone=strings.TrimSpace(p.DNSZone);if p.ID==""||p.Name==""||p.Kind==""||p.DNSZone==""||(p.Scope!="local"&&p.Scope!="global"){return false};caps:=make([]string,0,len(p.Capabilities));seen:=map[string]struct{}{};for _,c:=range p.Capabilities{c=strings.TrimSpace(c);if c!=""{if _,ok:=seen[c];!ok{seen[c]=struct{}{};caps=append(caps,c)}}};sort.Strings(caps);p.Capabilities=caps;r.mu.Lock();if r.providers==nil{r.providers=make(map[string]Provider)};r.providers[p.ID]=p;r.mu.Unlock();return true}
func(r *ProviderRegistry)Get(id string)(Provider,bool){if r==nil{return Provider{},false};r.mu.RLock();defer r.mu.RUnlock();p,ok:=r.providers[strings.TrimSpace(id)];if !ok{return Provider{},false};p.Capabilities=append([]string(nil),p.Capabilities...);return p,true}
func(r *ProviderRegistry)List()[]Provider{if r==nil{return []Provider{}};r.mu.RLock();defer r.mu.RUnlock();out:=make([]Provider,0,len(r.providers));for _,p:=range r.providers{p.Capabilities=append([]string(nil),p.Capabilities...);out=append(out,p)};sort.Slice(out,func(i,j int)bool{return out[i].ID<out[j].ID});return out}
func(r *ProviderRegistry)Select(scope,zone string)[]Provider{if r==nil{return []Provider{}};scope=strings.ToLower(strings.TrimSpace(scope));zone=strings.TrimSpace(zone);r.mu.RLock();defer r.mu.RUnlock();out:=make([]Provider,0);for _,p:=range r.providers{if p.Enabled&&p.Scope==scope&&p.DNSZone==zone{p.Capabilities=append([]string(nil),p.Capabilities...);out=append(out,p)}};sort.Slice(out,func(i,j int)bool{if out[i].Primary!=out[j].Primary{return out[i].Primary};return out[i].ID<out[j].ID});return out}
