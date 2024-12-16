package core

import (
	"Blockchain-Go/types"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBlockchain_AddBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	lenBlocks := 1000

	for i := 1; i <= lenBlocks; i++ {
		block := randomBlockWithSignature(t, uint32(i), getPreviousBlockHash(t, bc, uint32(i)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(uint32(i))
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 89, types.Hash{})))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, math.MaxUint32, types.Hash{})))
}

func newBlockChainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlockWithSignature(t, 0, types.Hash{}))

	assert.Nil(t, err)
	assert.NotNil(t, bc)
	assert.NotNil(t, bc.blockValidator)
	assert.NotNil(t, bc.storage)
	assert.Equal(t, bc.Height(), uint32(0))
	assert.True(t, bc.HasBlock(0))
	return bc
}

func getPreviousBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {

	previousHeader, err := bc.GetHeader(height - 1)

	assert.Nil(t, err)

	previousHeaderHash, err := previousHeader.Hash(types.HeaderHasher{})

	assert.Len(t, previousHeaderHash, 32)
	assert.Nil(t, err)

	return previousHeaderHash
}
