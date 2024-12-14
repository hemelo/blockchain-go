package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
}

type Server struct {
	ServerOpts
	rpcChannel  chan RPC
	quitChannel chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts:  opts,
		rpcChannel:  make(chan RPC),
		quitChannel: make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcChannel:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitChannel:
			break free
		case <-ticker.C:
			fmt.Println("tick")
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
