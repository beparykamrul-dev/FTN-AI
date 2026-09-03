package main

import (
	"fmt"
	"net/netip"
	"sort"
	"strings"
)

// FTNMainServerService is the authoritative origin of an FTN local service.
// Third-party/local CDN origins are deliberately excluded from this contract.
type FTNMainServerService struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Prefix       string `json:"prefix"`
	Port         uint16 `json:"port"`
	Protocol     string `json:"protocol"`
	MainServerID string `json:"main_server_id"`
	Enabled      bool   `json:"enabled"`
}

type FTNServicePlaneIntent struct {
	SubscriberPrefix string `json:"subscriber_prefix"`
	ServiceID        string `json:"service_id"`
	MainServerID     string `json:"main_server_id"`
	VRF              string `json:"vrf"`
	NextHop          string `json:"next_hop"`
	Approved         bool   `json:"approved"`
}

func NormalizeFTNMainServerService(in FTNMainServerService) (FTNMainServerService, error) {
	in.ID = strings.TrimSpace(in.ID)
	in.Name = strings.TrimSpace(in.Name)
	in.Prefix = strings.TrimSpace(in.Prefix)
	in.Protocol = strings.ToLower(strings.TrimSpace(in.Protocol))
	in.MainServerID = strings.TrimSpace(in.MainServerID)
	if in.ID == "" || in.Name == "" || in.MainServerID == "" || in.Prefix == "" {
		return FTNMainServerService{}, fmt.Errorf("service, main server and prefix are required")
	}
	p, err := netip.ParsePrefix(in.Prefix)
	if err != nil {
		return FTNMainServerService{}, fmt.Errorf("invalid service prefix: %w", err)
	}
	in.Prefix = p.Masked().String()
	if in.Port == 0 || in.Port > 65535 {
		return FTNMainServerService{}, fmt.Errorf("invalid service port")
	}
	if in.Protocol != "tcp" && in.Protocol != "udp" {
		return FTNMainServerService{}, fmt.Errorf("unsupported service protocol")
	}
	return in, nil
}

func NormalizeFTNServicePlaneIntent(in FTNServicePlaneIntent) (FTNServicePlaneIntent, error) {
	in.SubscriberPrefix = strings.TrimSpace(in.SubscriberPrefix)
	in.ServiceID = strings.TrimSpace(in.ServiceID)
	in.MainServerID = strings.TrimSpace(in.MainServerID)
	in.VRF = strings.TrimSpace(in.VRF)
	in.NextHop = strings.TrimSpace(in.NextHop)
	if in.SubscriberPrefix == "" || in.ServiceID == "" || in.MainServerID == "" || in.VRF == "" || in.NextHop == "" {
		return FTNServicePlaneIntent{}, fmt.Errorf("subscriber, service, main server, VRF and next-hop are required")
	}
	sp, err := netip.ParsePrefix(in.SubscriberPrefix)
	if err != nil { return FTNServicePlaneIntent{}, fmt.Errorf("invalid subscriber prefix: %w", err) }
	in.SubscriberPrefix = sp.Masked().String()
	if _, err := netip.ParseAddr(in.NextHop); err != nil { return FTNServicePlaneIntent{}, fmt.Errorf("invalid next-hop: %w", err) }
	if !in.Approved { return FTNServicePlaneIntent{}, fmt.Errorf("service-plane approval required") }
	return in, nil
}

func SortFTNMainServerServices(v []FTNMainServerService) []FTNMainServerService {
	out := append([]FTNMainServerService(nil), v...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
