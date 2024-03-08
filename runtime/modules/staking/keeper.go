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

package staking

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdkkeeper "cosmossdk.io/x/staking/keeper"
	sdkstaking "cosmossdk.io/x/staking/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/cockroachdb/errors"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ ValsetChangeProvider = &Keeper{}

// Keeper implements the StakingKeeper interface
// as a wrapper around Cosmos SDK x/staking keeper.
type Keeper struct {
	stakingKeeper *sdkkeeper.Keeper
}

// NewKeeper creates a new instance of the staking Keeper.
func NewKeeper(stakingKeeper *sdkkeeper.Keeper) *Keeper {
	return &Keeper{
		stakingKeeper: stakingKeeper,
	}
}

// delegate delegates the deposit to the validator.
func (k *Keeper) delegate(
	ctx context.Context, deposit *beacontypes.Deposit,
) (uint64, error) {
	validatorPK := &ed25519.PubKey{}
	err := validatorPK.Unmarshal(deposit.Pubkey)
	if err != nil {
		return 0, err
	}
	amount := deposit.Amount
	valConsAddr := sdk.GetConsAddress(validatorPK)
	validator, err := k.stakingKeeper.GetValidator(
		ctx, sdk.ValAddress(valConsAddr),
	)
	if err != nil {
		if errors.Is(err, sdkstaking.ErrNoValidatorFound) {
			validator, err = k.createValidator(validatorPK, amount)
			return validator.DelegatorShares.BigInt().Uint64(), err
		}
		return 0, err
	}
	newShares, err := k.stakingKeeper.Delegate(
		ctx, sdk.AccAddress(valConsAddr),
		sdkmath.NewIntFromUint64(amount),
		sdkstaking.Unbonded, validator, true)
	return newShares.BigInt().Uint64(), err
}

// undelegate undelegates the validator.
func (k *Keeper) undelegate(
	_ context.Context, _ *enginetypes.Withdrawal,
) (uint64, error) {
	// TODO: implement undelegate
	return 0, nil
}

// createValidator creates a new validator with the given public
// key and amount of tokens.
func (k *Keeper) createValidator(
	validatorPK sdkcrypto.PubKey,
	amount uint64) (sdkstaking.Validator, error) {
	stake := sdkmath.NewIntFromUint64(amount)
	valConsAddr := sdk.GetConsAddress(validatorPK)
	operator := sdk.ValAddress(valConsAddr).String()
	val, err := sdkstaking.NewValidator(
		operator, validatorPK,
		sdkstaking.Description{Moniker: validatorPK.String()})
	val.Tokens = stake
	val.DelegatorShares = sdkmath.LegacyNewDecFromInt(val.Tokens)
	return val, err
}

// ApplyChanges applies the deposits and withdrawals to the underlying
// staking module.
func (k *Keeper) ApplyChanges(
	ctx context.Context,
	deposits []*beacontypes.Deposit,
	withdrawals []*enginetypes.Withdrawal,
) error {
	for _, deposit := range deposits {
		_, err := k.delegate(ctx, deposit)
		if err != nil {
			return err
		}
	}
	for _, withdrawal := range withdrawals {
		_, err := k.undelegate(ctx, withdrawal)
		if err != nil {
			return err
		}
	}
	return nil
}
