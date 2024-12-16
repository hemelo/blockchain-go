package crypto

import (
	"Blockchain-Go/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/rs/zerolog/log"
	"math/big"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

type PublicKey struct {
	key *ecdsa.PublicKey
}

type Signature struct {
	s, r *big.Int
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate private key")
		panic(err)
	}

	return PrivateKey{key}
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.key, k.key.X, k.key.Y)
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())

	return types.AddressFromBytes(h[:20])
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		key: &k.key.PublicKey,
	}
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)

	if err != nil {
		return nil, err
	}

	return &Signature{s: s, r: r}, nil
}

func (sig Signature) Verify(publicKey PublicKey, data []byte) bool {
	return ecdsa.Verify(publicKey.key, data, sig.r, sig.s)
}
