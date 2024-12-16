package core

import (
	"Blockchain-Go/crypto"
	"fmt"
)

type Transaction struct {
	Data []byte

	From      crypto.PublicKey
	Signature *crypto.Signature
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
