package network

import (
	"Blockchain-Go/core"
	"Blockchain-Go/types"
	"fmt"
	"github.com/rs/zerolog/log"
)

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (tp *TxPool) Add(tx *core.Transaction) error {

	hash, err := tx.Hash(types.TransactionHasher{})

	if err != nil {
		log.Error().Err(err).Msg("could not hash transaction")
		return err
	}

	if tp.Has(hash) {
		return fmt.Errorf("transaction already in pool")
	}

	tp.transactions[hash] = tx

	return nil
}

func (tp *TxPool) Get(hash types.Hash) (*core.Transaction, bool) {
	tx, ok := tp.transactions[hash]
	return tx, ok
}

func (tp *TxPool) GetAll() []*core.Transaction {
	txs := make([]*core.Transaction, 0, len(tp.transactions))
	for _, tx := range tp.transactions {
		txs = append(txs, tx)
	}
	return txs
}

func (tp *TxPool) Remove(hash types.Hash) {
	delete(tp.transactions, hash)
}

func (tp *TxPool) Size() int {
	return len(tp.transactions)
}

func (tp *TxPool) Flush() {
	tp.transactions = make(map[types.Hash]*core.Transaction)
}

func (tp *TxPool) Has(hash types.Hash) bool {
	_, ok := tp.transactions[hash]
	return ok
}

func (tp *TxPool) HasTx(tx *core.Transaction) (bool, error) {

	hash, err := tx.Hash(types.TransactionHasher{})

	if err != nil {
		log.Error().Err(err).Msg("could not hash transaction")
		return false, err
	}

	return tp.Has(hash), nil
}
