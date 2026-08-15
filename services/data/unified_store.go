package data

import "context"

// Backend is the minimal contract used by FTN's unified data layer.
type Backend interface {
	Name() string
	Ping(context.Context) error
}

// Policy describes how a service selects a database backend.
type Policy struct {
	Primary   string
	Fallbacks []string
	LocalOnly bool
	Global    bool
}

// UnifiedStore keeps service code independent from a specific database engine.
type UnifiedStore struct {
	backends map[string]Backend
}

func NewUnifiedStore(backends ...Backend) *UnifiedStore {
	m := make(map[string]Backend, len(backends))
	for _, b := range backends {
		if b != nil && b.Name() != "" {
			m[b.Name()] = b
		}
	}
	return &UnifiedStore{backends: m}
}

func (s *UnifiedStore) Backend(name string) (Backend, bool) {
	b, ok := s.backends[name]
	return b, ok
}
