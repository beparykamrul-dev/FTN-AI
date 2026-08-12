package authenticator

import (
	"errors"
	"sync"
	"time"
)

var ErrKeyVersionUnavailable = errors.New("verification key version unavailable")
var ErrKeyVersionRevoked = errors.New("verification key version revoked")

// VerificationKeyPublication is the public portion of an authentication signing key.
// Private signing material never belongs in this structure.
type VerificationKeyPublication struct {
	KeyID       string
	Version     uint64
	Algorithm   string
	PublicKey   []byte
	PublishedAt time.Time
	RevokedAt   *time.Time
}

// KeyPublisher publishes verification keys to authenticator nodes.
// Implementations may persist this state in FTN's distributed registry.
type KeyPublisher interface {
	Publish(VerificationKeyPublication) error
	Revoke(keyID string, version uint64, at time.Time) error
	Get(keyID string, version uint64) (VerificationKeyPublication, error)
}

// KeyDistribution coordinates publication and local verification of public keys.
// It is intentionally independent of any specific storage or transport.
type KeyDistribution struct {
	mu        sync.RWMutex
	publisher KeyPublisher
	keys      map[string]VerificationKeyPublication
}

func NewKeyDistribution(publisher KeyPublisher) *KeyDistribution {
	return &KeyDistribution{publisher: publisher, keys: make(map[string]VerificationKeyPublication)}
}

func (d *KeyDistribution) Publish(key VerificationKeyPublication) error {
	if key.KeyID == "" || key.Version == 0 || len(key.PublicKey) == 0 {
		return errors.New("invalid verification key publication")
	}
	if key.PublishedAt.IsZero() {
		key.PublishedAt = time.Now().UTC()
	}
	if d.publisher != nil {
		if err := d.publisher.Publish(key); err != nil {
			return err
		}
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.keys[key.KeyID+":"+versionString(key.Version)] = clonePublication(key)
	return nil
}

func (d *KeyDistribution) Revoke(keyID string, version uint64, at time.Time) error {
	if keyID == "" || version == 0 {
		return errors.New("invalid key revocation")
	}
	if at.IsZero() {
		at = time.Now().UTC()
	}
	if d.publisher != nil {
		if err := d.publisher.Revoke(keyID, version, at); err != nil {
			return err
		}
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	key, ok := d.keys[keyID+":"+versionString(version)]
	if !ok {
		return ErrKeyVersionUnavailable
	}
	key.RevokedAt = &at
	d.keys[keyID+":"+versionString(version)] = clonePublication(key)
	return nil
}

func (d *KeyDistribution) Resolve(keyID string, version uint64) (VerificationKeyPublication, error) {
	key := keyID + ":" + versionString(version)
	d.mu.RLock()
	value, ok := d.keys[key]
	d.mu.RUnlock()
	if !ok && d.publisher != nil {
		published, err := d.publisher.Get(keyID, version)
		if err != nil {
			return VerificationKeyPublication{}, err
		}
		value = published
		d.mu.Lock()
		d.keys[key] = clonePublication(published)
		d.mu.Unlock()
	}
	if !ok && value.KeyID == "" {
		return VerificationKeyPublication{}, ErrKeyVersionUnavailable
	}
	if value.RevokedAt != nil {
		return VerificationKeyPublication{}, ErrKeyVersionRevoked
	}
	return clonePublication(value), nil
}

func clonePublication(in VerificationKeyPublication) VerificationKeyPublication {
	out := in
	out.PublicKey = append([]byte(nil), in.PublicKey...)
	if in.RevokedAt != nil {
		t := *in.RevokedAt
		out.RevokedAt = &t
	}
	return out
}

func versionString(v uint64) string {
	// Avoid fmt allocation in the hot authentication path by using a fixed buffer.
	const digits = "0123456789"
	var buf [20]byte
	i := len(buf)
	for v >= 10 {
		i--
		buf[i] = digits[v%10]
		v /= 10
	}
	i--
	buf[i] = digits[v]
	return string(buf[i:])
}
