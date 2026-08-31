package config

import "sync"

var (
	mu      sync.RWMutex
	current Config
	loaded  bool
)

func Init() (Config, error) {
	cfg := Load()
	if err := Validate(cfg); err != nil {
		return Config{}, err
	}
	mu.Lock()
	current = cfg
	loaded = true
	mu.Unlock()
	return cfg, nil
}

func Get() Config {
	mu.RLock()
	defer mu.RUnlock()
	return current
}

func Loaded() bool {
	mu.RLock()
	defer mu.RUnlock()
	return loaded
}
