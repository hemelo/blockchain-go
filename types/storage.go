package types

import (
	"Blockchain-Go/core"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
)

type Storage[T any] interface {
	Put(T) error
	Get(height uint32) (T, error)
	GetAll() ([]T, error)
	GetChunk(from, to uint32) ([]T, error)
	Count() (uint32, error)
}

type BlockMemStorage struct {
	lock   sync.RWMutex
	blocks map[uint32]*core.Block
}

func NewBlockMemStorage() *BlockMemStorage {

	log.Debug().Msg("creating new block memory storage")

	return &BlockMemStorage{
		blocks: make(map[uint32]*core.Block),
		lock:   sync.RWMutex{},
	}
}

func (st *BlockMemStorage) Put(b *core.Block) error {

	log.Debug().Msg("putting block into memory storage")

	st.lock.Lock()
	defer st.lock.Unlock()
	st.blocks[b.Height] = b

	return nil
}

func (st *BlockMemStorage) Get(height uint32) (*core.Block, error) {

	log.Debug().Msg("getting block from memory storage")

	st.lock.RLock()
	block, ok := st.blocks[height]
	st.lock.RUnlock()

	if !ok {
		log.Debug().Uint32("height", height).Msg("block not found")
		return nil, fmt.Errorf("block (%d) not found", height)
	}

	return block, nil
}

func (st *BlockMemStorage) GetAll() ([]*core.Block, error) {

	log.Debug().Msg("getting all blocks from memory storage")

	blocks := make([]*core.Block, 0, len(st.blocks))

	st.lock.RLock()
	defer st.lock.RUnlock()

	for _, block := range st.blocks {
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (st *BlockMemStorage) GetChunk(from, to uint32) ([]*core.Block, error) {

	log.Debug().Msg("getting chunk of blocks from memory storage")

	blocks := make([]*core.Block, 0, to-from+1)

	if from > to {
		log.Debug().Uint32("from", from).Uint32("to", to).Msg("invalid range")
		return nil, fmt.Errorf("invalid range")
	}

	count, err := st.Count()

	if err != nil {
		log.Error().Err(err).Msg("failed to count blocks")
		return nil, err
	}

	if to > count {
		log.Debug().Uint32("to", to).Uint32("length", uint32(len(st.blocks))).Msg("invalid range, to is greater than length, adjusting to length")
		to = uint32(len(st.blocks))
	}

	st.lock.RLock()

	for i := from; i <= to; i++ {
		block, err := st.Get(i)

		if err != nil {
			st.lock.RUnlock()
			log.Error().Err(err).Uint32("height", i).Msg("failed to get block")
			return nil, err
		}

		blocks = append(blocks, block)
	}

	st.lock.RUnlock()

	return blocks, nil
}

func (st *BlockMemStorage) Count() (uint32, error) {

	log.Debug().Msg("counting blocks in memory storage")

	st.lock.RLock()
	defer st.lock.RUnlock()
	return uint32(len(st.blocks)), nil
}
