package auth

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/argon2"
)

const (
	passwordTime    uint32 = 3
	passwordMemory  uint32 = 64 * 1024
	passwordThreads uint8  = 2
	passwordKeyLen  uint32 = 32
	passwordSaltLen        = 16
)

func HashCredential(password string) (hash, salt []byte, err error) {
	if password == "" {
		return nil, nil, errors.New("password is required")
	}
	salt = make([]byte, passwordSaltLen)
	if _, err = rand.Read(salt); err != nil {
		return nil, nil, err
	}
	hash = argon2.IDKey([]byte(password), salt, passwordTime, passwordMemory, passwordThreads, passwordKeyLen)
	return hash, salt, nil
}

func VerifyCredential(password string, hash, salt []byte) bool {
	if password == "" || len(hash) == 0 || len(salt) == 0 {
		return false
	}
	candidate := argon2.IDKey([]byte(password), salt, passwordTime, passwordMemory, passwordThreads, passwordKeyLen)
	if len(candidate) != len(hash) {
		return false
	}
	var diff byte
	for i := range candidate {
		diff |= candidate[i] ^ hash[i]
	}
	return diff == 0
}
