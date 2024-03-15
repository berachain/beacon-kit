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
	sdkbls "github.com/cosmos/cosmos-sdk/crypto/keys/bls12_381"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper implements the StakingKeeper interface
// as a wrapper around Cosmos SDK x/staking keeper.
type Keeper struct {
	*sdkkeeper.Keeper
}

func ProvideStakingKeeper(
	stakingKeeper *sdkkeeper.Keeper,
) *Keeper {
	return NewKeeper(stakingKeeper)
}

// NewKeeper creates a new instance of the staking Keeper.
func NewKeeper(stakingKeeper *sdkkeeper.Keeper) *Keeper {
	if stakingKeeper == nil {
		panic("staking keeper is required")
	}

	return &Keeper{
		Keeper: stakingKeeper,
	}
}

// delegate delegates the deposit to the validator.
func (k *Keeper) ApplyDeposit(
	ctx context.Context,
	deposit *beacontypes.Deposit,
) error {
	pubKey := &sdkbls.PubKey{Key: deposit.Pubkey}
	consAddr := sdk.GetConsAddress(pubKey)

	// Attempt to get the validator
	validator, err := k.GetValidatorByConsAddr(ctx, consAddr)
	switch {
	// if it is not found, then we create a new one.
	case errors.Is(err, sdkstaking.ErrNoValidatorFound):
		_, err = k.createValidator(pubKey, deposit)
		if err != nil {
			return err
		}
		return nil
	// if there is any other error, we return it.
	case err != nil:
		return err
	// Otherwise, we found a validator and we deposit to it.
	default:
		_, err = k.Delegate(
			ctx,
			sdk.AccAddress(consAddr),
			sdkmath.NewIntFromUint64(deposit.Amount),
			sdkstaking.Unbonded,
			validator,
			true,
		)
		return err
	}
}

// undelegate undelegates the validator.
func (k *Keeper) ApplyWithdrawal(
	_ context.Context, _ *enginetypes.Withdrawal,
) error {
	return nil
}

// createValidator creates a new validator with the given public
// key and amount of tokens.
func (k *Keeper) createValidator(
	validatorPubkey *sdkbls.PubKey,
	deposit *beacontypes.Deposit,
) (sdkstaking.Validator, error) {
	// Verify the deposit data against the signature.
	// Deposit message is the deposit without the signature.
	root, err := (&beacontypes.Deposit{
		Index:       deposit.Index,
		Pubkey:      deposit.Pubkey,
		Credentials: deposit.Credentials,
		Amount:      deposit.Amount,
	}).HashTreeRoot()
	if err != nil {
		return sdkstaking.Validator{},
			errors.Wrapf(err, "could not get signing root")
	}

	// TODO: Embed the domain into the signing data.
	//nolint:lll // Will be removed later.
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#domain-types
	if !validatorPubkey.VerifySignature(root[:], deposit.Signature) {
		return sdkstaking.Validator{}, errors.New("could not verify signature")
	}

	// Create a new validator with x/staking.
	stake := sdkmath.NewIntFromUint64(deposit.Amount)
	operator := sdk.ValAddress(deposit.Credentials).String()
	newValidator, err := sdkstaking.NewValidator(
		operator,
		validatorPubkey,
		sdkstaking.Description{Moniker: validatorPubkey.Address().String()},
	)
	newValidator.Tokens = stake
	newValidator.DelegatorShares = sdkmath.LegacyNewDecFromInt(
		newValidator.Tokens,
	)
	return newValidator, err
}
