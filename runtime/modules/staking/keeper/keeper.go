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
	sdkbls "github.com/cosmos/cosmos-sdk/crypto/keys/bls12_381"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

var _ ValsetChangeProvider = &Keeper{}

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
func (k *Keeper) delegate(
	ctx context.Context,
	deposit *beacontypes.Deposit,
) (uint64, error) {
	validatorPubkey := &sdkbls.PubKey{Key: deposit.Pubkey}
	validator, err := k.GetValidatorByPubKey(ctx, validatorPubkey)
	if err != nil {
		if errors.Is(err, sdkstaking.ErrNoValidatorFound) {
			validator, err = k.createValidator(validatorPubkey, deposit)
			return validator.DelegatorShares.BigInt().Uint64(), err
		}
		return 0, err
	}
	valConsAddr := sdk.GetConsAddress(validatorPubkey)
	newShares, err := k.Delegate(
		ctx, sdk.AccAddress(valConsAddr),
		sdkmath.NewIntFromUint64(deposit.Amount),
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
	validatorPubkey *sdkbls.PubKey,
	deposit *beacontypes.Deposit,
) (sdkstaking.Validator, error) {
	// Verify the deposit data against the signature.
	// Deposit message is the deposit without the signature.
	depositMsg := &beacontypes.Deposit{
		Pubkey:      deposit.Pubkey,
		Credentials: deposit.Credentials,
		Amount:      deposit.Amount,
	}
	root, err := depositMsg.HashTreeRoot()
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
	newValidator.DelegatorShares = sdkmath.LegacyNewDecFromInt(newValidator.Tokens)
	return newValidator, err
}

// GetValidatorByPubKey returns the validator with the given public key.
func (k *Keeper) GetValidatorByPubKey(
	ctx context.Context,
	pubkey types.PubKey,
) (sdkstaking.Validator, error) {
	return k.GetValidatorByConsAddr(
		ctx,
		sdk.GetConsAddress(pubkey),
	)
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
