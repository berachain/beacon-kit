package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
	"github.com/stretchr/testify/require"
)

func TestSSZVectorBasic(t *testing.T) {
	t.Run("SizeSSZ for uint8 vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZByte]{1, 2, 3, 4, 5}
		require.Equal(t, 5, vector.SizeSSZ())
	})

	t.Run("SizeSSZ for byte slice vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZUInt8]{1, 2, 3, 4, 5, 6, 7, 8}
		require.Equal(t, 48, vector.SizeSSZ())
	})

	t.Run("SizeSSZ for uint64 vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZUInt64]{1, 2, 3, 4, 5}
		require.Equal(t, 40, vector.SizeSSZ())
	})

	t.Run("SizeSSZ for bool vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZBool]{true, false, true}
		require.Equal(t, 3, vector.SizeSSZ())
	})

	t.Run("SizeSSZ for empty vector", func(t *testing.T) {
		vector := types.SSZVectorBasic[types.SSZUInt64]{}
		require.Equal(t, 0, vector.SizeSSZ())
	})
}
