package storage

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type BlobStore struct {
	Root string
	AEAD cipher.AEAD
}

func NewBlobStore(root string, key []byte) (*BlobStore, error) {
	block, err := aes.NewCipher(key)
	if err != nil { return nil, err }
	aead, err := cipher.NewGCM(block)
	if err != nil { return nil, err }
	if root == "" { return nil, errors.New("storage root is required") }
	return &BlobStore{Root: root, AEAD: aead}, nil
}

func (s *BlobStore) Put(ctx context.Context, key string, plaintext []byte) error {
	select { case <-ctx.Done(): return ctx.Err(); default: }
	if key == "" || s == nil || s.AEAD == nil { return errors.New("invalid blob store") }
	nonce := make([]byte, s.AEAD.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return err }
	ciphertext := s.AEAD.Seal(nonce, nonce, plaintext, nil)
	path := filepath.Join(s.Root, filepath.Clean(key))
	if filepath.Dir(path) != filepath.Clean(s.Root) && !filepath.IsAbs(key) {
		// Nested storage keys are allowed under Root; traversal is rejected below.
	}
	if filepath.IsAbs(key) || filepath.Clean(key) == ".." || len(filepath.Clean(key)) >= 3 && filepath.Clean(key)[:3] == "../" {
		return errors.New("invalid storage key")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil { return err }
	tmp, err := os.CreateTemp(filepath.Dir(path), ".ftn-mail-*")
	if err != nil { return err }
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err = tmp.Chmod(0600); err == nil { _, err = tmp.Write(ciphertext) }
	if closeErr := tmp.Close(); err == nil { err = closeErr }
	if err != nil { return err }
	return os.Rename(tmpName, path)
}

func (s *BlobStore) Get(ctx context.Context, key string) ([]byte, error) {
	select { case <-ctx.Done(): return nil, ctx.Err(); default: }
	if key == "" || filepath.IsAbs(key) { return nil, errors.New("invalid storage key") }
	path := filepath.Join(s.Root, filepath.Clean(key))
	data, err := os.ReadFile(path)
	if err != nil { return nil, err }
	if len(data) < s.AEAD.NonceSize() { return nil, errors.New("corrupt mail blob") }
	return s.AEAD.Open(nil, data[:s.AEAD.NonceSize()], data[s.AEAD.NonceSize():], nil)
}
