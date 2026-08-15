package ftnstorage

import "crypto/sha256"

// ChunkRef identifies content by its cryptographic digest instead of by a
// large metadata object. The actual bytes are owned by a storage backend.
type ChunkRef struct {
	Hash   [32]byte
	Size   uint64
	Codec  string
}

func NewChunkRef(data []byte, codec string) ChunkRef {
	return ChunkRef{Hash: sha256.Sum256(data), Size: uint64(len(data)), Codec: codec}
}

// SameContent reports whether two references identify the same content.
func SameContent(a, b ChunkRef) bool {
	return a.Hash == b.Hash && a.Size == b.Size
}
