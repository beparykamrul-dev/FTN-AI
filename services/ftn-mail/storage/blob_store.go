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
	"strings"
)

type BlobStore struct { Root string; AEAD cipher.AEAD }

func NewBlobStore(root string, key []byte) (*BlobStore, error) {
	root = strings.TrimSpace(root)
	if root == "" { return nil, errors.New("storage root is required") }
	block, err := aes.NewCipher(key); if err != nil { return nil, err }
	aead, err := cipher.NewGCM(block); if err != nil { return nil, err }
	return &BlobStore{Root:root, AEAD:aead}, nil
}

func safeBlobPath(root, key string) (string, error) {
	root, key = strings.TrimSpace(root), strings.TrimSpace(key)
	if root == "" || key == "" || filepath.IsAbs(key) { return "", errors.New("invalid storage key") }
	clean := filepath.Clean(key)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) { return "", errors.New("invalid storage key") }
	rootClean, err := filepath.Abs(root); if err != nil { return "", err }
	path, err := filepath.Abs(filepath.Join(rootClean, clean)); if err != nil { return "", err }
	if path != rootClean && !strings.HasPrefix(path, rootClean+string(filepath.Separator)) { return "", errors.New("invalid storage key") }
	return path,nil
}

func (s *BlobStore) Put(ctx context.Context, key string, plaintext []byte) error {
	if ctx == nil { return errors.New("context is required") }
	select { case <-ctx.Done(): return ctx.Err(); default: }
	if s == nil || s.AEAD == nil { return errors.New("invalid blob store") }
	path, err := safeBlobPath(s.Root,key); if err != nil { return err }
	nonce := make([]byte,s.AEAD.NonceSize()); if _,err=io.ReadFull(rand.Reader,nonce); err!=nil{return err}
	ciphertext := s.AEAD.Seal(nonce,nonce,plaintext,nil)
	if err=os.MkdirAll(filepath.Dir(path),0700); err!=nil{return err}
	tmp,err:=os.CreateTemp(filepath.Dir(path),".ftn-mail-*"); if err!=nil{return err}; tmpName:=tmp.Name(); defer os.Remove(tmpName)
	if err=tmp.Chmod(0600); err==nil { _,err=tmp.Write(ciphertext) }; if closeErr:=tmp.Close(); err==nil {err=closeErr}; if err!=nil{return err}
	return os.Rename(tmpName,path)
}

func (s *BlobStore) Get(ctx context.Context, key string) ([]byte,error) {
	if ctx == nil { return nil, errors.New("context is required") }
	select { case <-ctx.Done(): return nil,ctx.Err(); default: }
	if s==nil || s.AEAD==nil{return nil,errors.New("invalid blob store")}
	path,err:=safeBlobPath(s.Root,key);if err!=nil{return nil,err};data,err:=os.ReadFile(path);if err!=nil{return nil,err}
	if len(data)<s.AEAD.NonceSize(){return nil,errors.New("corrupt mail blob")}
	return s.AEAD.Open(nil,data[:s.AEAD.NonceSize()],data[s.AEAD.NonceSize():],nil)
}
