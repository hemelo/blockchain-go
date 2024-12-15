package types

import (
	"Blockchain-Go/core"
	"crypto/sha256"
)

type Hasher[T any] interface {
	Hash(T) Hash
}

type BlockHasher struct {
}

func (BlockHasher) Hash(b *core.Block) Hash {
	h := sha256.Sum256(b.HeaderData())
	return Hash(h[:])
}
