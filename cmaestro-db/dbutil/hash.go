package dbutil

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
)

// HasherReader calculates the hash of a byte stream
// As an underlying io.Reader is read from, the hash is updated
type HasherReader struct {
	hash   hash.Hash
	reader io.Reader
}

// NewHasherReader creates a new HasherReader from a provided io.Raeder
func NewHasherReader(r io.Reader) HasherReader {
	hash := sha1.New()
	reader := io.TeeReader(r, hash)
	return HasherReader{hash, reader}
}

// Hash returns the hash value
// Ensure all contents of the underlying io.Reader have been read
func (h HasherReader) Hash() string {
	return hex.EncodeToString(h.hash.Sum(nil))
}

// Read allows HasherReader to conform to io.Reader interface
func (h HasherReader) Read(p []byte) (n int, err error) {
	return h.reader.Read(p)
}
