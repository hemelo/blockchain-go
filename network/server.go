package network

import (
	"Blockchain-Go/core"
	"Blockchain-Go/crypto"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

type ServerOpts struct {
	Transports []Transport
	PrivateKey *crypto.PrivateKey
	BlockTime  time.Duration
}

type Server struct {
	ServerOpts
	txPool      *TxPool
	isValidator bool
	rpcChannel  chan RPC
	quitChannel chan struct{}
}

func NewServer(opts ServerOpts) *Server {

	if opts.BlockTime == 0 {
		opts.BlockTime = 5 * time.Second
	}

	return &Server{
		ServerOpts:  opts,
		txPool:      NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChannel:  make(chan RPC),
		quitChannel: make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChannel:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitChannel:
			break free
		case <-ticker.C:
			log.Debug().Msg("block time reached")

			if s.isValidator {
				s.CreateNewBlock()
			}
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcChannel <- rpc
			}
		}(tr)
	}
}

func (s *Server) HandleTransaction(tx *core.Transaction) error {
	log.Debug().Msg("handling transaction")

	result, err := tx.Verify()

	if err != nil {
		log.Error().Err(err).Msg("could not verify transaction")
		return err
	}

	if !result {
		log.Error().Msg("transaction is not valid")
		return fmt.Errorf("transaction is not valid")
	}

	return s.txPool.Add(tx)
}

func (s *Server) CreateNewBlock() {
	log.Debug().Msg("creating new block")
}
