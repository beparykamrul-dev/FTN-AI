package data

import "context"

// Route resolves a backend according to policy and availability.
func (s *UnifiedStore) Route(ctx context.Context, p Policy) (Backend, bool) {
	candidates := append([]string{p.Primary}, p.Fallbacks...)
	for _, name := range candidates {
		b, ok := s.backends[name]
		if !ok {
			continue
		}
		if b.Ping(ctx) == nil {
			return b, true
		}
	}
	return nil, false
}
