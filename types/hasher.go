package types

import (
	"Blockchain-Go/core"
	"crypto/sha256"
)

type Hasher[T any] interface {
	Hash(T) (Hash, error)
}

type HeaderHasher struct {
}

func (HeaderHasher) Hash(header *core.Header) (Hash, error) {

	headerBytes, err := header.Bytes()

	if err != nil {
		return Hash{}, err
	}

	hash := sha256.Sum256(headerBytes)

	return Hash(hash[:]), nil
}
