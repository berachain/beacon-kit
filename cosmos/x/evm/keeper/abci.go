package keeper

import (
	"context"
)

// Precommit runs on the Cosmos-SDK lifecycle Precommit().
func (k *Keeper) EndBlock(_ context.Context) error {
	// sCtx := sdk.UnwrapSDKContext(ctx)
	return nil
}

// PrepareCheckState runs on the Cosmos-SDK lifecycle PrepareCheckState().
func (k *Keeper) PrepareCheckState(_ context.Context) error {
	return nil
}
