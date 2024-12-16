package core

import (
	"Blockchain-Go/types"
	"fmt"
	"github.com/rs/zerolog/log"
)

type Blockchain struct {
	storage        types.Storage[*Block]
	headers        []*Header
	blockValidator types.Validator[*Block]
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {

	log.Debug().Msg("creating new blockchain")

	bc := &Blockchain{
		headers: []*Header{},
	}

	bc.storage = types.NewBlockMemStorage()
	bc.blockValidator = types.NewBlockValidator(bc)

	err := bc.AddBlockWithoutValidation(genesis)

	if err != nil {
		log.Error().Err(err).Msg("failed to add genesis block")
	} else {
		log.Debug().Msg("genesis block added")
	}

	log.Debug().Msg("blockchain created")

	return bc, err
}

func (bc *Blockchain) SetBlockValidator(blockValidator types.Validator[*Block]) {
	bc.blockValidator = blockValidator
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height < bc.Height()
}

func (bc *Blockchain) Height() uint32 {

	if len(bc.headers) > 0 {
		return uint32(len(bc.headers) - 1)
	}

	return 0
}

func (bc *Blockchain) AddBlockWithoutValidation(block *Block) error {
	bc.headers = append(bc.headers, block.Header)
	res := bc.storage.Put(block)

	if res != nil {
		log.Error().Err(res).Uint32("height", block.Height).Msg("failed to store block")
	} else {
		log.Debug().Uint32("height", block.Height).Msg("block stored")
	}

	return res
}

func (bc *Blockchain) AddBlock(block *Block) error {

	log.Debug().Uint32("height", block.Height).Msg("adding block %d")

	if err := bc.blockValidator.Validate(block); err != nil {
		log.Error().Err(err).Uint32("height", block.Height).Msgf("failed to validate block")
		return err
	}

	return bc.AddBlockWithoutValidation(block)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height < bc.Height() {
		return bc.headers[height], nil
	}

	return nil, fmt.Errorf("chain does not contain block at height %d", height)
}

func (bc *Blockchain) GetBlock(height uint32) (*Block, error) {
	if height < bc.Height() {
		return bc.storage.Get(height)
	}

	return nil, fmt.Errorf("chain does not contain block at height %d", height)
}

func (bc *Blockchain) GetBlocks(from uint32, to uint32) ([]*Block, error) {
	if from > to {
		return nil, fmt.Errorf("invalid range")
	}

	if to > bc.Height() {
		to = bc.Height()
	}

	return bc.storage.GetChunk(from, to)
}

func (bc *Blockchain) GetAllBlocks() ([]*Block, error) {
	return bc.storage.GetAll()
}

func (bc *Blockchain) getBlockFromArray(blocks []*Block, height uint32) (*Block, error) {

	if blocks == nil {
		log.Error().Msg("array of blocks on parameter is nil")
		return nil, fmt.Errorf("blocks is nil")
	}

	for _, block := range blocks {
		if block.Height == height {
			return block, nil
		}
	}

	return nil, fmt.Errorf("block not found")
}

func (bc *Blockchain) ValidateChain() error {
	const chunkSize = 200
	totalBlocks := bc.Height()

	log.Debug().Uint32("blocks", totalBlocks).Uint32("per_chunk", chunkSize).Uint32("total_chunks", (totalBlocks/chunkSize)+1).Msg("validating blockchain")

	for start := uint32(0); start < totalBlocks; start += chunkSize {
		end := start + chunkSize
		if end > totalBlocks {
			end = totalBlocks
		}

		blocks, err := bc.GetBlocks(start, end)

		if err != nil {
			log.Error().Err(err).Uint32("start", start).Uint32("end", end).Msg("failed to get blocks")
			return err
		}

		log.Debug().Uint32("start", start).Uint32("end", end).Msg("validating chunk")

		for i := start; i < end; i++ {
			log.Debug().Uint32("height", i).Msg("validating block")

			block, err := bc.getBlockFromArray(blocks, i)
			if err != nil {
				log.Error().Err(err).Uint32("height", i).Msg("failed to get block")
				return err
			}

			if err := bc.blockValidator.Validate(block); err != nil {
				log.Error().Err(err).Uint32("height", i).Msg("block is invalid")
				return fmt.Errorf("block at height %d is invalid: %v", i, err)
			}
		}
	}

	log.Debug().Msg("blockchain validation completed")
	return nil
}
