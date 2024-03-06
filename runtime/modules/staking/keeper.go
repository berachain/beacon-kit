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
	"bytes"
	"context"
	"encoding/binary"
	"errors"

	sdkmath "cosmossdk.io/math"
	sdkkeeper "cosmossdk.io/x/staking/keeper"
	sdkstaking "cosmossdk.io/x/staking/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	beacontypesv1 "github.com/itsdevbear/bolaris/beacon/core/types/v1"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
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
	ctx context.Context, deposit *beacontypesv1.Deposit,
) (uint64, error) {
	// Get validator's sdk public key.
	validatorPK := &ed25519.PubKey{}
	err := validatorPK.Unmarshal(deposit.GetValidatorPubkey())
	if err != nil {
		return 0, err
	}
	amount := deposit.GetAmount()
	valConsAddr := sdk.GetConsAddress(validatorPK)
	validator, err := k.stakingKeeper.GetValidator(
		ctx, sdk.ValAddress(valConsAddr),
	)
	if err != nil {
		if errors.Is(err, sdkstaking.ErrNoValidatorFound) {
			// Create a new validator on the first deposit.
			validator, err = k.createValidator(
				validatorPK, deposit,
			)
			return validator.DelegatorShares.BigInt().Uint64(), err
		}
		return 0, err
	}
	delegatorAddr := sdk.AccAddress(deposit.GetStakingCredentials())
	newShares, err := k.stakingKeeper.Delegate(
		ctx,
		delegatorAddr,
		sdkmath.NewIntFromUint64(amount),
		sdkstaking.Unbonded,
		validator,
		true,
	)
	return newShares.BigInt().Uint64(), err
}

// undelegate undelegates the validator.
func (k *Keeper) undelegate(
	_ context.Context, _ *enginev1.Withdrawal,
) (uint64, error) {
	// TODO: implement undelegate
	return 0, nil
}

// createValidator creates a new validator with the given public
// key and amount of tokens.
func (k *Keeper) createValidator(
	pubkey sdkcrypto.PubKey,
	deposit *beacontypesv1.Deposit,
) (sdkstaking.Validator, error) {
	validatorPK := deposit.GetValidatorPubkey()
	stakingCredentials := deposit.GetStakingCredentials()
	amount := deposit.GetAmount()

	// Verify the deposit data against the signature
	msg := make([]byte, 0)
	msg = append(msg, validatorPK...)
	msg = append(msg, stakingCredentials...)
	// Execution layer uses big endian encoding
	msg = binary.BigEndian.AppendUint64(msg, amount)
	sigPK, err := ethcrypto.Ecrecover(msg, deposit.GetSignature())
	if err != nil {
		return sdkstaking.Validator{}, err
	}
	if !bytes.Equal(sigPK, deposit.GetValidatorPubkey()) {
		return sdkstaking.Validator{}, errors.New("invalid signature")
	}

	// Create a new validator.
	stake := sdkmath.NewIntFromUint64(amount)
	operator := sdk.ValAddress(stakingCredentials).String()
	newValidator, err := sdkstaking.NewValidator(
		operator,
		pubkey,
		sdkstaking.Description{Moniker: pubkey.String()},
	)
	newValidator.Tokens = stake
	newValidator.DelegatorShares = sdkmath.LegacyNewDecFromInt(newValidator.Tokens)
	return newValidator, err
}

// ApplyChanges applies the deposits and withdrawals to the underlying
// staking module.
func (k *Keeper) ApplyChanges(
	ctx context.Context,
	deposits []*beacontypesv1.Deposit,
	withdrawals []*enginev1.Withdrawal,
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
