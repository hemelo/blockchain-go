package core

import (
	"Blockchain-Go/crypto"
	"Blockchain-Go/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Height:        height,
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	tx := Transaction{
		Data: []byte("foo"),
	}

	return NewBlock(header, []Transaction{tx})
}

func TestBlock_Hash(t *testing.T) {
	b := randomBlock(1)
	assert.NotNil(t, b.Hash(types.BlockHasher{}))
}

func TestBlock_Sign(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	b := randomBlock(1)

	assert.Nil(t, b.Sign(privateKey))
	assert.NotNil(t, b.Signature)
	assert.NotNil(t, b.Validator)
}

func TestBlock_Verify(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	b := randomBlock(1)

	assert.Nil(t, b.Sign(privateKey))

	verified, err := b.Verify()

	assert.Nil(t, err)
	assert.True(t, verified)

	otherPrivateKey := crypto.GeneratePrivateKey()

	b.Validator = otherPrivateKey.PublicKey()

	verified2, err2 := b.Verify()

	assert.Nil(t, err2)
	assert.False(t, verified2)

	b.Height = 100
	verified3, err3 := b.Verify()

	assert.Nil(t, err3)
	assert.False(t, verified3)

}
