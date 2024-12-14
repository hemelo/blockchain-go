package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalTransport(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.peers[trb.Addr()], trb)
	assert.Equal(t, tra.peers[trb.Addr()], trb)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	message := []byte("Hello world")
	assert.Nil(t, tra.SendMessage(trb.Addr(), message))

	rpc := <-trb.Consume()
	assert.Equal(t, rpc.Payload, message)
	assert.Equal(t, rpc.From, tra.Addr())
}
