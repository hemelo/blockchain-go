package core

import (
	"Blockchain-Go/crypto"
	"Blockchain-Go/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlock_Hash(t *testing.T) {
	b := randomBlock(t, 1, types.Hash{})
	hash, err := b.Header.Hash(types.HeaderHasher{})

	assert.Nil(t, err)
	assert.Len(t, hash, 32)
}

func TestBlock_Sign(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 1, types.Hash{})

	assert.Nil(t, b.Sign(privateKey))
	assert.NotNil(t, b.Signature)
	assert.NotNil(t, b.Validator)
}

func TestBlock_Verify(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 1, types.Hash{})

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

func randomBlock(t *testing.T, height uint32, previousBlockHash types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: previousBlockHash,
		Height:        height,
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	block := NewBlock(header, []Transaction{})
	return block
}

func randomBlockWithSignature(t *testing.T, height uint32, previousBlockHash types.Hash) *Block {
	privateKey := crypto.GeneratePrivateKey()
	block := randomBlock(t, height, previousBlockHash)

	tx := randomTxWithSignature(t)
	block.AddTransaction(tx)

	assert.Nil(t, block.Sign(privateKey))

	return block
}
