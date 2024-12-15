package core

import (
	"Blockchain-Go/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction_Sign(t *testing.T) {

	privateKey := crypto.GeneratePrivateKey()
	data := []byte("Foo")

	tx := Transaction{
		Data: data,
	}

	assert.Nil(t, tx.Sign(privateKey))
	assert.NotNil(t, tx.Signature)
	assert.NotNil(t, tx.PublicKey)

}

func TestTransaction_Verify(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	data := []byte("Foo")

	tx := Transaction{
		Data: data,
	}

	assert.Nil(t, tx.Sign(privateKey))

	verified, err := tx.Verify()

	assert.Nil(t, err)
	assert.True(t, verified)

	otherPrivateKey := crypto.GeneratePrivateKey()

	tx.PublicKey = otherPrivateKey.PublicKey()

	verified2, err2 := tx.Verify()

	assert.Nil(t, err2)
	assert.False(t, verified2)

}
