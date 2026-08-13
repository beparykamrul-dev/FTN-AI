package runtime

import (
	"context"
	"fmt"
	"sync"
)

type Service interface {
	Name() string
	Start(context.Context) error
	Stop(context.Context) error
	Health(context.Context) error
}

type ServiceManager struct {
	mu       sync.RWMutex
	services map[string]Service
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{services: make(map[string]Service)}
}

func (m *ServiceManager) Register(service Service) error {
	if service == nil || service.Name() == "" {
		return fmt.Errorf("invalid service")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.services[service.Name()]; exists {
		return fmt.Errorf("service already registered: %s", service.Name())
	}
	m.services[service.Name()] = service
	return nil
}

func (m *ServiceManager) StartAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for name, service := range m.services {
		if err := service.Start(ctx); err != nil {
			return fmt.Errorf("start %s: %w", name, err)
		}
	}
	return nil
}

func (m *ServiceManager) StopAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for name, service := range m.services {
		if err := service.Stop(ctx); err != nil {
			return fmt.Errorf("stop %s: %w", name, err)
		}
	}
	return nil
}

func (m *ServiceManager) Health(ctx context.Context) map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]error, len(m.services))
	for name, service := range m.services {
		result[name] = service.Health(ctx)
	}
	return result
}
