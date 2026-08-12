package auth

import "golang.org/x/crypto/argon2"

// HashPassword returns an Argon2id password hash payload.
// Parameters are intentionally conservative and should be benchmarked on the
// production hardware before rollout.
func HashPassword(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, 3, 64*1024, 2, 32)
}
