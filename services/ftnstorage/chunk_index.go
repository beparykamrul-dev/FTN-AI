package ftnstorage

import "sync"

// ChunkIndex is an in-memory index boundary for content-addressed chunks.
// Persistence belongs to the configured FTN storage metadata backend.
type ChunkIndex struct {
	mu sync.RWMutex
	m  map[[32]byte]ChunkRef
}

func NewChunkIndex() *ChunkIndex { return &ChunkIndex{m: make(map[[32]byte]ChunkRef)} }

func (i *ChunkIndex) Put(ref ChunkRef) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.m[ref.Hash] = ref
}

func (i *ChunkIndex) Get(hash [32]byte) (ChunkRef, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	ref, ok := i.m[hash]
	return ref, ok
}

func (i *ChunkIndex) Has(hash [32]byte) bool {
	_, ok := i.Get(hash)
	return ok
}
