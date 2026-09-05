package gis

import (
	"fmt"
	"net/netip"
	"sort"
	"sync"
)

type IPAM struct {
	mu     sync.RWMutex
	assets map[string]IPAsset
}

func NewIPAM() *IPAM {
	return &IPAM{assets: make(map[string]IPAsset)}
}

func (i *IPAM) Upsert(asset IPAsset) error {
	if i == nil {
		return fmt.Errorf("ipam is required")
	}
	if asset.ID == "" || asset.IP == "" {
		return fmt.Errorf("id and ip are required")
	}
	if _, err := netip.ParseAddr(asset.IP); err != nil {
		return fmt.Errorf("invalid ip: %w", err)
	}
	if asset.CIDR != "" {
		if _, err := netip.ParsePrefix(asset.CIDR); err != nil {
			return fmt.Errorf("invalid cidr: %w", err)
		}
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.assets[asset.ID] = asset
	return nil
}

func (i *IPAM) List() []IPAsset {
	if i == nil {
		return nil
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	out := make([]IPAsset, 0, len(i.assets))
	for _, a := range i.assets {
		out = append(out, a)
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a].IP != out[b].IP {
			return out[a].IP < out[b].IP
		}
		return out[a].ID < out[b].ID
	})
	return out
}
