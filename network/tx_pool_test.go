package network

import (
	"Blockchain-Go/core"
	"Blockchain-Go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTxPool(t *testing.T) {
	txPool := NewTxPool()

	tx := core.NewTransaction([]byte("Foo"))

	hash, err := tx.Hash(types.TransactionHasher{})

	assert.Nil(t, err)
	assert.Nil(t, txPool.Add(tx))
	assert.True(t, txPool.Has(hash))
	assert.Equal(t, 1, txPool.Size())

	txPool.Flush()
	assert.Equal(t, 0, txPool.Size())
}
