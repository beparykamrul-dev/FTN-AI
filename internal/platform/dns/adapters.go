package dns

import("context";"fmt";"strings")
type AdapterFactory struct{}
func(AdapterFactory)Build(cfg ProviderConfig)(ProviderAdapter,error){cfg.ID=strings.TrimSpace(cfg.ID);if !cfg.Enabled{return nil,fmt.Errorf("provider %q is disabled",cfg.ID)};if cfg.ID==""||len(cfg.ID)>128{return nil,fmt.Errorf("invalid provider ID")};if len(cfg.Endpoint)>4096{return nil,fmt.Errorf("provider endpoint is too large")};switch cfg.Type{case ProviderPowerDNS,ProviderTechnitium,ProviderCoreDNS,ProviderUnbound,ProviderDNSDist,ProviderGoDNS,ProviderAnycast,ProviderDNSPod,ProviderCloudflare,ProviderAkamai:return &genericAdapter{provider:cfg.Type},nil;default:return nil,fmt.Errorf("unsupported DNS provider: %s",cfg.Type)}}
type genericAdapter struct{provider ProviderType}
func(a *genericAdapter)Type()ProviderType{if a==nil{return ""};return a.provider}
func(a *genericAdapter)ApplyZone(ctx context.Context,_ Zone)error{if a==nil{return fmt.Errorf("DNS adapter is required")};if ctx==nil{return fmt.Errorf("context is required")};return fmt.Errorf("DNS provider %q has no mutation implementation",a.provider)}
func(a *genericAdapter)DeleteZone(ctx context.Context,_ string)error{if a==nil{return fmt.Errorf("DNS adapter is required")};if ctx==nil{return fmt.Errorf("context is required")};return fmt.Errorf("DNS provider %q has no mutation implementation",a.provider)}
