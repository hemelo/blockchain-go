package types

import (
	"Blockchain-Go/core"
	"fmt"
)

type Validator[T any] interface {
	Validate(T) error
}

type BlockValidator struct {
	blockchain *core.Blockchain
}

func NewBlockValidator(blockchain *core.Blockchain) *BlockValidator {
	return &BlockValidator{blockchain}
}

func (bv *BlockValidator) Validate(block *core.Block) error {

	hash, err := block.Header.Hash(HeaderHasher{})

	if err != nil {
		return fmt.Errorf("failed to hash block (%d) header: %w", block.Height, err)
	}

	if bv.blockchain.HasBlock(block.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", block.Height, hash)
	}

	if block.Height != bv.blockchain.Height()+1 {
		return fmt.Errorf("block (%d) is too high", block.Height)
	}

	innerBlockValidation, innerBlockValidationErr := block.Verify()

	if innerBlockValidationErr != nil {
		return innerBlockValidationErr
	}

	if !innerBlockValidation {
		return fmt.Errorf("block validation failed")
	}

	previousHeader, previousHeaderErr := bv.blockchain.GetHeader(block.Height - 1)

	if previousHeaderErr != nil {
		return previousHeaderErr
	}

	previousHeaderHash, previousHeaderHashErr := previousHeader.Hash(HeaderHasher{})

	if previousHeaderHashErr != nil {
		return fmt.Errorf("failed to hash previous block (%d) header: %w", block.Height-1, previousHeaderHashErr)
	}

	if previousHeaderHash != block.PrevBlockHash {
		return fmt.Errorf("block (%d) has invalid previous block hash", block.Height)
	}

	return nil
}
