package dns

import (
	"fmt"
	"net"
	"strings"
)

type SPFRecord struct {
	Domain string `json:"domain"`
	Mechanisms []string `json:"mechanisms,omitempty"`
	IP4 []string `json:"ip4,omitempty"`
	IP6 []string `json:"ip6,omitempty"`
	Include []string `json:"include,omitempty"`
	All string `json:"all"`
}

func (r SPFRecord) Validate() error {
	if strings.TrimSpace(r.Domain) == "" { return fmt.Errorf("domain is required") }
	if r.All != "" && r.All != "pass" && r.All != "fail" && r.All != "softfail" && r.All != "neutral" { return fmt.Errorf("invalid all qualifier") }
	for _, cidr := range append(append([]string{}, r.IP4...), r.IP6...) {
		if _, _, err := net.ParseCIDR(cidr); err != nil { return fmt.Errorf("invalid network: %s", cidr) }
	}
	return nil
}

func (r SPFRecord) TXTValue() (string, error) {
	if err := r.Validate(); err != nil { return "", err }
	parts := []string{"v=spf1"}
	parts = append(parts, r.Mechanisms...)
	for _, ip := range r.IP4 { parts = append(parts, "ip4:"+ip) }
	for _, ip := range r.IP6 { parts = append(parts, "ip6:"+ip) }
	for _, inc := range r.Include { if strings.TrimSpace(inc) != "" { parts = append(parts, "include:"+strings.TrimSpace(inc)) } }
	if r.All != "" { parts = append(parts, "~all") }
	return strings.Join(parts, " "), nil
}
