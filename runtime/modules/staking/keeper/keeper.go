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
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/cockroachdb/errors"
	sdkbls "github.com/cosmos/cosmos-sdk/crypto/keys/bls12_381"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper implements the StakingKeeper interface
// as a wrapper around Cosmos SDK x/staking keeper.
type Keeper struct {
	*sdkkeeper.Keeper
	bk BankKeeper
}

// NewKeeper creates a new instance of the staking Keeper.
func NewKeeper(
	stakingKeeper *sdkkeeper.Keeper,
	bankKeeper *BankKeeper,
) *Keeper {
	if stakingKeeper == nil {
		panic("staking keeper is required")
	}

	if bankKeeper == nil {
		panic("bank keeper is required")
	}

	return &Keeper{
		Keeper: stakingKeeper,
		bk:     *bankKeeper,
	}
}

// delegate delegates the deposit to the validator.
func (k *Keeper) IncreaseConsensusPower(
	ctx context.Context,
	delegator beacontypes.DepositCredentials,
	pubkey [bls12381.PubKeyLength]byte,
	amount uint64,
	signature []byte,
	index uint64,
) error {
	validator, err := k.getValidatorFromPubkey(
		ctx,
		&sdkbls.PubKey{Key: pubkey[:]},
	)
	switch {
	// if it is not found, then we create a new one.
	case errors.Is(err, sdkstaking.ErrNoValidatorFound):
		_, err = k.createValidator(
			ctx,
			delegator,
			pubkey,
			amount,
			signature,
			index,
		)
		return err
	case err != nil:
		return err
	// Otherwise, we found a validator and we deposit to it.
	default:
		return k.mintAndDelegate(ctx, delegator[:], validator, amount)
	}
}

// undelegate undelegates the validator.
func (k *Keeper) DecreaseConsensusPower(
	ctx context.Context,
	delegator primitives.ExecutionAddress,
	pubkey [bls12381.PubKeyLength]byte,
	amount uint64,
) error {
	validator, err := k.getValidatorFromPubkey(
		ctx,
		&sdkbls.PubKey{Key: pubkey[:]},
	)
	if err != nil {
		return err
	}

	return k.withdrawAndBurn(
		ctx,
		delegator[:],
		validator,
		amount,
	)
}

// createValidator creates a new validator with the given public
// key and amount of tokens.
func (k *Keeper) createValidator(
	ctx context.Context,
	delegator beacontypes.DepositCredentials,
	validatorPubkey [bls12381.PubKeyLength]byte,
	amount uint64,
	signature []byte,
	index uint64,
) (sdkstaking.Validator, error) {
	// Verify the deposit data against the signature.
	// Deposit message is the deposit without the signature.
	root, err := (&beacontypes.Deposit{
		Index:       index,
		Pubkey:      validatorPubkey[:],
		Credentials: delegator,
		Amount:      amount,
	}).HashTreeRoot()
	if err != nil {
		return sdkstaking.Validator{},
			errors.Wrapf(err, "could not get signing root")
	}

	// TODO: Embed the domain into the signing data.
	//nolint:lll // Will be removed later.
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#domain-types
	sdkblsPubKey := &sdkbls.PubKey{Key: validatorPubkey[:]}
	if !sdkblsPubKey.VerifySignature(root[:], signature) {
		return sdkstaking.Validator{}, errors.New("could not verify signature")
	}

	delegatorAddress, err := delegator.ToExecutionAddress()
	if err != nil {
		return sdkstaking.Validator{}, err
	}

	// Create a new validator with x/staking.
	newValidator, err := sdkstaking.NewValidator(
		// TODO: make the byte prefixed credentials into a hard type.
		sdk.AccAddress(delegatorAddress[:]).String(),
		sdkblsPubKey,
		sdkstaking.Description{Moniker: sdkblsPubKey.Address().String()},
	)
	if err != nil {
		return sdkstaking.Validator{}, err
	}

	if err = k.mintAndDelegate(
		ctx,
		delegator[:],
		newValidator,
		amount,
	); err != nil {
		return sdkstaking.Validator{}, err
	}

	newValidator.Tokens = sdkmath.NewIntFromUint64(amount)
	newValidator.DelegatorShares = sdkmath.LegacyNewDecFromInt(
		newValidator.Tokens,
	)
	return newValidator, err
}

// GetValidatorFromPubkey returns the validator from the given public key.
func (k *Keeper) getValidatorFromPubkey(
	ctx context.Context,
	pubkey *sdkbls.PubKey,
) (sdkstaking.Validator, error) {
	consAddr := sdk.GetConsAddress(pubkey)
	return k.GetValidatorByConsAddr(ctx, consAddr)
}

// mintAndDelegate mints the staking coins and delegates them to the
// specified validator.
func (k *Keeper) mintAndDelegate(
	ctx context.Context,
	delegator []byte,
	validator sdkstaking.Validator,
	amount uint64,
) error {
	var err error
	coins := sdk.Coins{
		sdk.NewCoin(StakingUnit, sdkmath.NewIntFromUint64(amount)),
	}

	// Mint the coins to the bonded pool.
	if err = k.bk.MintCoins(
		ctx,
		StakingModuleName,
		coins,
	); err != nil {
		return err
	}

	// Transfer the coins from the module account to the delegator.
	if err = k.bk.SendCoinsFromModuleToAccount(
		ctx,
		StakingModuleName,
		sdk.AccAddress(delegator[12:]),
		coins,
	); err != nil {
		return err
	}

	_, err = k.Delegate(
		ctx,
		sdk.AccAddress(delegator[12:]),
		sdkmath.NewIntFromUint64(amount),
		sdkstaking.Unbonded, // TODO: Check if this is the correct value.
		validator,
		true,
	)
	return err
}

// withdrawAndBurn undelegates the staking coins from the validator
// and burns them.
func (k *Keeper) withdrawAndBurn(
	ctx context.Context,
	delegator []byte,
	validator sdkstaking.Validator,
	amount uint64,
) error {
	var err error
	valBz, err := k.ValidatorAddressCodec().
		StringToBytes(validator.GetOperator())
	if err != nil {
		return err
	}

	shares, err := validator.SharesFromTokens(sdkmath.NewIntFromUint64(amount))
	if err != nil {
		return err
	}

	_, err = k.ValidateUnbondAmount(
		ctx,
		sdk.AccAddress(delegator),
		valBz,
		shares.TruncateInt(),
	)
	if err != nil {
		return err
	}

	_, _, err = k.Undelegate(
		ctx,
		sdk.AccAddress(delegator),
		valBz,
		shares,
	)
	if err != nil {
		return err
	}

	coinsToBurn, err := k.CompleteUnbonding(
		ctx,
		sdk.AccAddress(delegator),
		valBz,
	)
	if err != nil {
		return err
	}

	return k.bk.BurnCoins(ctx, sdk.AccAddress(delegator), coinsToBurn)
}
