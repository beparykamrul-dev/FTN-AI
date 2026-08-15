package ftnstorage

import "crypto/sha256"

// ScrubResult reports integrity verification for one stored chunk.
type ScrubResult struct {
	Hash    [32]byte
	Size    uint64
	Healthy bool
}

// Scrub verifies bytes against their content address. It does not repair or
// mutate storage; a recovery worker can act on an unhealthy result.
func Scrub(ref ChunkRef, data []byte) ScrubResult {
	h := sha256.Sum256(data)
	return ScrubResult{Hash: h, Size: uint64(len(data)), Healthy: h == ref.Hash && uint64(len(data)) == ref.Size}
}
