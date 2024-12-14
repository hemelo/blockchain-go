package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr           NetAddr
	consumeChannel chan RPC
	peers          map[NetAddr]*LocalTransport
	lock           sync.RWMutex
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:           addr,
		consumeChannel: make(chan RPC, 1024),
		peers:          make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeChannel
}

func (t *LocalTransport) Connect(transport Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[transport.Addr()] = transport.(*LocalTransport)

	return nil
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]

	if !ok {
		return fmt.Errorf("%s could not send message to %s", t.addr, to)
	}

	peer.consumeChannel <- RPC{
		From:    t.Addr(),
		Payload: payload,
	}

	return nil
}
