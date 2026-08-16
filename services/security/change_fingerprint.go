package security

import "crypto/sha256"

// ChangeFingerprint produces a stable identifier from a security scan identity.
func ChangeFingerprint(runID, commit, policyID string, version uint64) string {
	input := runID + "|" + commit + "|" + policyID + "|" + string(rune(version))
	sum := sha256.Sum256([]byte(input))
	return stringHex(sum[:])
}

func stringHex(b []byte) string {
	const hex = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hex[v>>4]
		out[i*2+1] = hex[v&15]
	}
	return string(out)
}
