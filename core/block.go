package core

import (
	"Blockchain-Go/crypto"
	"Blockchain-Go/types"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
)

type Header struct {
	Version       uint32
	PrevBlockHash types.Hash
	DataHash      types.Hash
	Timestamp     uint64
	Height        uint32

	// For cache purposes
	hash types.Hash
}

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
}

func NewBlock(h *Header, txs []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
	}
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

func (b *Block) Encode(w io.Writer, enc types.Encoder[*Block]) error {
	return enc.Encode(w, b)
}

func (b *Block) Decode(r io.Reader, dec types.Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (h *Header) Hash(hasher types.Hasher[*Header]) (types.Hash, error) {

	if h.hash.IsZero() {
		hash, err := hasher.Hash(h)

		if err != nil {
			return hash, err
		}

		h.hash = hash
	}

	return h.hash, nil
}

func (b *Block) Sign(privateKey crypto.PrivateKey) error {

	log.Debug().Msgf("signing block %d", b.Height)

	headerBytes, err := b.Header.Bytes()

	if err != nil {
		return err
	}

	blockSignature, err := privateKey.Sign(headerBytes)

	if err != nil {
		return err
	}

	b.Validator = privateKey.PublicKey()
	b.Signature = blockSignature

	return nil
}

func (b *Block) Verify() (bool, error) {

	log.Debug().Uint32("height", b.Height).Msg("Verifying block")

	if b.Signature == nil {
		return false, fmt.Errorf("block has no signature")
	}

	headerBytes, err := b.Header.Bytes()

	if err != nil {

		return false, err
	}

	signValidation := b.Signature.Verify(b.Validator, headerBytes)

	if !signValidation {
		return false, nil
	}

	for _, tx := range b.Transactions {

		log.Debug().Uint32("height", b.Height).Str("from", tx.From.Address().String()).Msg("verifying transaction")

		txValidation, txValidationErr := tx.Verify()

		if txValidationErr != nil {
			log.Debug().Err(txValidationErr).Uint32("height", b.Height).Str("from", tx.From.Address().String()).Msg("block has invalid transaction")
			return false, txValidationErr
		}

		if !txValidation {
			log.Debug().Uint32("height", b.Height).Str("from", tx.From.Address().String()).Msg("block has invalid transaction, sign validation failed")
			return false, nil
		}
	}

	if signValidation {
		log.Debug().Uint32("height", b.Height).Msg("block is valid")
	} else {
		log.Debug().Uint32("height", b.Height).Msg("block is invalid")
	}

	return signValidation, nil
}

func (h *Header) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	if enc.Encode(h) != nil {
		return nil, fmt.Errorf("could not encode header")
	}

	return buf.Bytes(), nil
}
