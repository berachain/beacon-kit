package cosmos

import (
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ChainQuerier struct {
	ContextGetter func(height int64, prove bool) (sdk.Context, error)
	Service       service.BeaconStorageBackend
	ABCI          types.ABCI
}
