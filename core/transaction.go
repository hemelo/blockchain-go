package core

import (
	"Blockchain-Go/crypto"
	"Blockchain-Go/types"
	"fmt"
)

type Transaction struct {
	Data []byte

	From      crypto.PublicKey
	Signature *crypto.Signature

	// For cache purposes
	hash types.Hash
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Sign(privateKey crypto.PrivateKey) error {

	txSignature, err := privateKey.Sign(tx.Data)

	if err != nil {
		return err
	}

	tx.From = privateKey.PublicKey()
	tx.Signature = txSignature

	return nil
}

func (tx *Transaction) Verify() (bool, error) {

	if tx.Signature == nil {
		return false, fmt.Errorf("tx has no signature")
	}

	return tx.Signature.Verify(tx.From, tx.Data), nil
}

func (tx *Transaction) Hash(hasher types.Hasher[*Transaction]) (types.Hash, error) {

	if tx.hash.IsZero() {
		hash, err := hasher.Hash(tx)

		if err != nil {
			return hash, err
		}

		tx.hash = hash
	}

	return tx.hash, nil
}
