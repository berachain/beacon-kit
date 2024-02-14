package keeper

import (
	"context"
	"fmt"

	"github.com/itsdevbear/bolaris/runtime/modules/beacon/types"
)

// Querier is a struct that holds the keeper.
type Querier struct {
	*Keeper
}

// NewQuerier creates a new querier.
func NewQuerier(k *Keeper) Querier {
	return Querier{
		Keeper: k,
	}
}

// FinalizedEth1Block returns the finalized eth1 block.
func (q *Querier) FinalizedEth1Block(ctx context.Context, req *types.FinalizedEth1BlockRequest) (*types.FinalizedEth1BlockResponse, error) {
	fmt.Println("HENLO")
	return &types.FinalizedEth1BlockResponse{
		Eth1BlockHash: q.BeaconState(ctx).GetFinalizedEth1BlockHash().Hex(),
	}, nil
}
