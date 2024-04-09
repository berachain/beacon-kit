package cosmos

import (
	rpc "github.com/berachain/beacon-kit/mod/api"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ rpc.ChainQuerier = ChainQuerier{}
)

type ChainQuerier struct {
	ContextGetter func(height int64, prove bool) (sdk.Context, error)
	Service       service.BeaconStorageBackend
	ABCI          types.ABCI
}
