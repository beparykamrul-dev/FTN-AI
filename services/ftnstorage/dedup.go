package ftnstorage

// DedupResult describes whether bytes were already represented by a known
// content-addressed chunk.
type DedupResult struct {
	Ref     ChunkRef
	Existed bool
}

// PutDeduplicated inserts a chunk reference only when its content hash is new.
func PutDeduplicated(index *ChunkIndex, data []byte, codec string) DedupResult {
	ref := NewChunkRef(data, codec)
	if index.Has(ref.Hash) {
		return DedupResult{Ref: ref, Existed: true}
	}
	index.Put(ref)
	return DedupResult{Ref: ref}
}
