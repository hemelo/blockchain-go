package crypto

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeyPairGenerateAndSignVerify(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()

	require.NotPanics(t, func() {
		publicKey.Address()
	})

	msg := []byte("Hello")

	sig, err := privateKey.Sign(msg)
	require.NoError(t, err)
	require.NotNil(t, sig)

	assert.True(t, sig.Verify(publicKey, msg))
}

func TestKeyPairGenerateAndSignVerifyWithDifferentKey(t *testing.T) {
	privateKey := GeneratePrivateKey()

	msg := []byte("Hello")

	sig, err := privateKey.Sign(msg)
	require.NoError(t, err)
	require.NotNil(t, sig)

	diffPrivateKey := GeneratePrivateKey()
	diffPublicKey := diffPrivateKey.PublicKey()

	assert.False(t, sig.Verify(diffPublicKey, msg))
}
