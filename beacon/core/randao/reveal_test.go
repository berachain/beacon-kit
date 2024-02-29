package randao_test

import (
	"github.com/itsdevbear/bolaris/beacon/core/randao"
	blst2 "github.com/itsdevbear/bolaris/beacon/core/randao/blst"
	"testing"

	"github.com/berachain/comet-bls12-381/bls/blst"
	"github.com/stretchr/testify/require"
)

func TestReveal_Verify(t *testing.T) {
	key, err := blst.RandKey()
	require.NoError(t, err)

	someSigningData := randao.SigningData{
		Epoch:   1234,
		ChainID: "berachain-1",
	}

	reveal, err := blst2.NewRandaoReveal(someSigningData, key)
	require.NoError(t, err)

	require.True(t, reveal.Verify(key.PublicKey().Marshal(), someSigningData))

	// Test with wrong signing data
	anotherSigningData := randao.SigningData{
		Epoch:   1234,
		ChainID: "berachain-2",
	}

	require.False(t, reveal.Verify(key.PublicKey().Marshal(), anotherSigningData))
}
