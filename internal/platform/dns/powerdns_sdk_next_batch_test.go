package dns

import (
	"context"
	"testing"
)

func TestPowerDNSSDKRejectsWrongProviderAndNilContext(t *testing.T) {
	if _, err := NewPowerDNSSDK(ProviderConfig{Type:ProviderKind("other"), Endpoint:"https://dns.example"}); err == nil { t.Fatal("wrong provider type must fail") }
	p, err := NewPowerDNSSDK(ProviderConfig{Type:ProviderPowerDNS, Endpoint:"https://dns.example"})
	if err != nil { t.Fatal(err) }
	if err := p.Health(nil); err == nil { t.Fatal("nil context must fail") }
	if err := p.Health(context.Background()); err != nil { t.Fatal(err) }
}
