package main

import (
	"errors"
	"sort"
	"strings"
)

type FTNPlatformService struct {
	ID string `json:"id"`
	Protocols []string `json:"protocols"`
	Transports []string `json:"transports"`
	Enabled bool `json:"enabled"`
}

type FTNPlatformRegistry struct { Services []FTNPlatformService `json:"services"` }

func NormalizeFTNPlatformRegistry(r FTNPlatformRegistry) (FTNPlatformRegistry, error) {
	seen := map[string]bool{}
	for i := range r.Services {
		s := &r.Services[i]
		s.ID = strings.ToLower(strings.TrimSpace(s.ID))
		if s.ID == "" || seen[s.ID] { return FTNPlatformRegistry{}, errors.New("ftn_service_identity_invalid") }
		seen[s.ID] = true
		s.Protocols = normalizeList(s.Protocols)
		s.Transports = normalizeList(s.Transports)
		if s.Enabled && (len(s.Protocols) == 0 || len(s.Transports) == 0) { return FTNPlatformRegistry{}, errors.New("ftn_service_capabilities_required") }
	}
	sort.Slice(r.Services, func(i,j int) bool { return r.Services[i].ID < r.Services[j].ID })
	return r, nil
}

func (r FTNPlatformRegistry) Service(id string) (FTNPlatformService, bool) {
	id = strings.ToLower(strings.TrimSpace(id))
	for _, s := range r.Services { if s.ID == id { return s, true } }
	return FTNPlatformService{}, false
}
