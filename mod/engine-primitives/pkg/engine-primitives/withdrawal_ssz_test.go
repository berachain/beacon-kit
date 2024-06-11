package engineprimitives_test

import (
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestWithdrawalSSZ(t *testing.T) {
	withdrawal := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(2),
		Address:   [20]byte{},
		Amount:    math.Gwei(100),
	}

	data, err := withdrawal.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	err = withdrawal.UnmarshalSSZ(data)
	require.NoError(t, err)

	size := withdrawal.SizeSSZ()
	require.Equal(t, 44, size)

	tree, errHashTree := withdrawal.HashTreeRoot()
	require.NoError(t, errHashTree)
	require.NotNil(t, tree)
}

func TestWithdrawalGetTree(t *testing.T) {
	withdrawal := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(2),
		Address:   [20]byte{},
		Amount:    math.Gwei(100),
	}

	tree, err := withdrawal.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}
