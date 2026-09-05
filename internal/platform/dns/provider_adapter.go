package dns

import("context";"fmt";"strings";"sync")
type ProviderType = ProviderKind
type ProviderConfig struct{ID string `json:"id"`;Type ProviderType `json:"type"`;Endpoint string `json:"endpoint,omitempty"`;Enabled bool `json:"enabled"`}
type ProviderAdapter interface{Type() ProviderType;ApplyZone(context.Context,Zone)error;DeleteZone(context.Context,string)error}
type AdapterRegistry struct{mu sync.RWMutex;adapters map[ProviderType]ProviderAdapter}
func NewAdapterRegistry()*AdapterRegistry{return &AdapterRegistry{adapters:make(map[ProviderType]ProviderAdapter)}}
func(r *AdapterRegistry)Register(adapter ProviderAdapter)error{if r==nil{return fmt.Errorf("adapter registry is required")};if adapter==nil{return fmt.Errorf("adapter is required")};typ:=ProviderType(strings.ToLower(strings.TrimSpace(string(adapter.Type()))));if typ==""{return fmt.Errorf("adapter type is required")};r.mu.Lock();defer r.mu.Unlock();if r.adapters==nil{r.adapters=make(map[ProviderType]ProviderAdapter)};if _,exists:=r.adapters[typ];exists{return fmt.Errorf("adapter already registered: %s",typ)};r.adapters[typ]=adapter;return nil}
func(r *AdapterRegistry)Get(t ProviderType)(ProviderAdapter,bool){if r==nil{return nil,false};t=ProviderType(strings.ToLower(strings.TrimSpace(string(t))));if t==""{return nil,false};r.mu.RLock();a,ok:=r.adapters[t];r.mu.RUnlock();return a,ok}
func NormalizeProviderEndpoint(endpoint string)string{return strings.TrimRight(strings.TrimSpace(endpoint),"/")}
