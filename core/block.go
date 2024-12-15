package core

import (
	"Blockchain-Go/crypto"
	"Blockchain-Go/types"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

type Header struct {
	Version       uint32
	PrevBlockHash types.Hash
	DataHash      types.Hash
	Timestamp     uint64
	Height        uint32
}

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// For cache purposes
	hash types.Hash
}

func NewBlock(h *Header, txs []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
	}
}

func (b *Block) Encode(w io.Writer, enc types.Encoder[*Block]) error {
	return enc.Encode(w, b)
}

func (b *Block) Decode(r io.Reader, dec types.Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (b *Block) Hash(hasher types.Hasher[*Block]) types.Hash {

	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}

	return b.hash
}

func (b *Block) Sign(privateKey crypto.PrivateKey) error {
	blockSignature, err := privateKey.Sign(b.HeaderData())

	if err != nil {
		return err
	}

	b.Validator = privateKey.PublicKey()
	b.Signature = blockSignature

	return nil
}

func (b *Block) Verify() (bool, error) {
	if b.Signature == nil {
		return false, fmt.Errorf("block has no signature")
	}

	return b.Signature.Verify(b.Validator, b.HeaderData()), nil
}

func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(b.Header)

	return buf.Bytes()
}
