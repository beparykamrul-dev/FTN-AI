package security

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

// CanonicalFingerprint is platform-independent: every numeric field is encoded
// explicitly instead of relying on Go string/rune conversion semantics.
func CanonicalFingerprint(runID, commit, policyID string, version uint64) string {
	buf := make([]byte, 0, len(runID)+len(commit)+len(policyID)+8+3)
	buf = append(buf, runID...)
	buf = append(buf, 0)
	buf = append(buf, commit...)
	buf = append(buf, 0)
	buf = append(buf, policyID...)
	buf = append(buf, 0)
	var v [8]byte
	binary.BigEndian.PutUint64(v[:], version)
	buf = append(buf, v[:]...)
	sum := sha256.Sum256(buf)
	return hex.EncodeToString(sum[:])
}
