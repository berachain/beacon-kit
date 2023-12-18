package keeper

import (
	"context"

	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/beacon/execution"
)

var LatestForkChoiceKey = []byte("latestForkChoice")

type (
	Keeper struct {
		// consensusAPI is the consensus API
		executionClient execution.EngineCaller
		storeKey        storetypes.StoreKey
		forkchoiceState *enginev1.ForkchoiceState
	}
)

// NewKeeper creates new instances of the polaris Keeper.
func NewKeeper(
	executionClient execution.EngineCaller,
	storeKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		executionClient: executionClient,
		storeKey:        storeKey,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx context.Context) log.Logger {
	return sdk.UnwrapSDKContext(ctx).Logger()
}

func (k *Keeper) UpdateHoodForkChoice(forkchoiceState *enginev1.ForkchoiceState) {
	k.forkchoiceState = forkchoiceState
}
