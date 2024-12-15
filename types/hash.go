package types

import (
	"encoding/hex"
	"fmt"
	"math/rand"
)

type Hash [32]uint8

func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("given bytes with length %d should be 32 bytes ", len(b))
		panic(msg)
	}

	var value [32]uint8

	for i := 0; i < 32; i++ {
		value[i] = uint8(b[i])
	}

	return Hash(value)
}

func RandomBytes(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}

func RandomHash() Hash {
	return HashFromBytes(RandomBytes(32))
}

func (h Hash) IsZero() bool {
	for i := 0; i < len(h); i++ {
		if h[i] != 0 {
			return false
		}
	}

	return true
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}
