// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package keeper

import (
	"context"

	stakingtypes "cosmossdk.io/x/staking/types"
	cosmoslib "github.com/berachain/beacon-kit/lib/cosmos"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakingHooks struct.
type StakingHooks struct {
	*cosmoslib.UnimplementedStakingHooks
	k *Keeper
}

// Verify that the Hooks struct implements the stakingtypes.StakingHooks
// interface.
var _ stakingtypes.StakingHooks = StakingHooks{}

// Create new stakinghooks hooks.
func (k *Keeper) Hooks() StakingHooks {
	return StakingHooks{k: k}
}

// initialize validator distribution record.
func (h StakingHooks) AfterValidatorCreated(
	ctx context.Context,
	valAddr sdk.ValAddress,
) error {
	pubkey, err := h.k.vcp.GetValidatorPubkeyFromValAddress(ctx, valAddr)
	if err != nil {
		return err
	}

	return h.k.beaconStore.AddValidator(ctx, pubkey[:])
}

// AfterConsensusPubKeyUpdate does nothing and returns nil.
func (h StakingHooks) AfterConsensusPubKeyUpdate(
	ctx context.Context,
	fromPubkey cryptotypes.PubKey,
	toPubkey cryptotypes.PubKey,
	_ sdk.Coin,
) error {
	return h.k.beaconStore.UpdateValidator(
		ctx,
		fromPubkey.Bytes(),
		toPubkey.Bytes(),
	)
}
