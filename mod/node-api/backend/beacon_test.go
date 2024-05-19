package backend_test

import (
	"context"
	"testing"

	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/node-api/backend/mocks"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/stretchr/testify/require"
)

func TestGetGenesisValidatorsRoot(t *testing.T) {
	sdb := &mocks.StateDB{}
	b := backend.New(func(context.Context, string) backend.StateDB {
		return sdb
	})
	sdb.EXPECT().GetGenesisValidatorsRoot().Return(primitives.Root{0x01}, nil)
	root, err := b.GetGenesis(context.Background())
	require.NoError(t, err)
	require.Equal(t, primitives.Root{0x01}, root)
}