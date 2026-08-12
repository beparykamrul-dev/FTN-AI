package transport

import "golang.org/x/crypto/argon2"

func verifyArgon2id(password string, hash, salt []byte) bool {
	if password == "" || len(hash) == 0 || len(salt) == 0 {
		return false
	}
	candidate := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, uint32(len(hash)))
	if len(candidate) != len(hash) {
		return false
	}
	var diff byte
	for i := range candidate {
		diff |= candidate[i] ^ hash[i]
	}
	return diff == 0
}
