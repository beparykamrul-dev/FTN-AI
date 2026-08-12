package authenticator

import (
	"errors"
	"sync"
	"time"
)

var ErrKeyNotFound = errors.New("authentication key not found")
var ErrKeyRevoked = errors.New("authentication key revoked")

// KeyState describes the lifecycle of an FTN authentication key.
type KeyState string

const (
	KeyActive       KeyState = "active"
	KeyVerifyOnly   KeyState = "verify_only"
	KeyRetired      KeyState = "retired"
	KeyRevoked      KeyState = "revoked"
)

// SigningKey contains public verification material and lifecycle metadata.
// Private key material must remain behind the configured secret/key boundary.
type SigningKey struct {
	ID          string
	Version     uint64
	Algorithm   string
	PublicKey   []byte
	State       KeyState
	CreatedAt   time.Time
	ActivatedAt time.Time
	RetiredAt   *time.Time
}

// KeyManager provides concurrency-safe key publication and emergency revocation.
type KeyManager struct {
	mu      sync.RWMutex
	current SigningKey
	keys    map[string]SigningKey
}

func NewKeyManager(initial SigningKey) *KeyManager {
	return &KeyManager{current: initial, keys: map[string]SigningKey{initial.ID: initial}}
}

func (m *KeyManager) Current() (SigningKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current.ID == "" || m.current.State != KeyActive {
		return SigningKey{}, ErrKeyNotFound
	}
	return cloneKey(m.current), nil
}

func (m *KeyManager) Publish(key SigningKey) error {
	if key.ID == "" || key.Version == 0 || len(key.PublicKey) == 0 {
		return errors.New("invalid signing key")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if old, ok := m.keys[key.ID]; ok && old.State == KeyRevoked {
		return ErrKeyRevoked
	}
	key.PublicKey = append([]byte(nil), key.PublicKey...)
	m.keys[key.ID] = key
	if key.State == KeyActive && key.Version >= m.current.Version {
		m.current = key
	}
	return nil
}

func (m *KeyManager) Get(id string) (SigningKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key, ok := m.keys[id]
	if !ok || key.State == KeyRetired || key.State == KeyRevoked {
		return SigningKey{}, ErrKeyNotFound
	}
	return cloneKey(key), nil
}

func (m *KeyManager) Revoke(id string, now time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key, ok := m.keys[id]
	if !ok {
		return ErrKeyNotFound
	}
	key.State = KeyRevoked
	key.RetiredAt = &now
	m.keys[id] = key
	if m.current.ID == id {
		m.current = SigningKey{}
	}
	return nil
}

func cloneKey(k SigningKey) SigningKey {
	k.PublicKey = append([]byte(nil), k.PublicKey...)
	return k
}
