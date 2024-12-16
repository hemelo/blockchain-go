package core

import (
	"Blockchain-Go/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction_Verify(t *testing.T) {

	tx := randomTxWithSignature(t)

	verified, err := tx.Verify()

	assert.Nil(t, err)
	assert.True(t, verified)

	otherPrivateKey := crypto.GeneratePrivateKey()

	tx.From = otherPrivateKey.PublicKey()

	verified2, err2 := tx.Verify()

	assert.Nil(t, err2)
	assert.False(t, verified2)

}

func randomTxWithSignature(t *testing.T) *Transaction {
	privateKey := crypto.GeneratePrivateKey()
	data := []byte("Foo")

	tx := &Transaction{
		Data: data,
	}

	assert.Nil(t, tx.Sign(privateKey))
	assert.NotNil(t, tx.Signature)
	assert.NotNil(t, tx.From)

	return tx
}
