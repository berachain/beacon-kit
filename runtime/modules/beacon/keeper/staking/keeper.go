// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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
	"errors"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/itsdevbear/bolaris/config"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
)

var _ Staking = &Keeper{}

// Keeper implements the Staking interface
// as a wrapper around Cosmos SDK x/staking keeper.
// with a queue of deposits to be processed.
type Keeper struct {
	stakingKeeper stakingkeeper.Keeper
	deposits      []*consensusv1.Deposit
	beaconCfg     *config.Beacon
}

// NewKeeper creates a new instance of the staking wrapper.
func NewKeeper(stakingKeeper stakingkeeper.Keeper, beaconCfg *config.Beacon) Keeper {
	return Keeper{
		stakingKeeper: stakingKeeper,
		deposits:      make([]*consensusv1.Deposit, 0),
		beaconCfg:     beaconCfg,
	}
}

// AddDeposit queues a deposit to the staking module.
func (k Keeper) AddDeposit(_ context.Context, deposit *consensusv1.Deposit) error {
	k.deposits = append(k.deposits, deposit)
	return nil
}

// ProcessDeposits processes the queued deposits (up to the limit MaxDepositsPerBlock).
func (k Keeper) ProcessDeposits(ctx context.Context) error {
	var processedDeposits uint64
	for processedDeposits < k.beaconCfg.Limits.MaxDepositsPerBlock && len(k.deposits) > 0 {
		deposit := k.deposits[0]
		if err := k.processDeposit(ctx, deposit); err != nil {
			return err
		}
		k.deposits = k.deposits[1:]
		processedDeposits++
	}
	return nil
}

// processDeposit processes a single deposit and delegates the tokens to the validator.
func (k Keeper) processDeposit(ctx context.Context, deposit *consensusv1.Deposit) error {
	validatorPK := &ed25519.PubKey{}
	depositData := deposit.GetData()
	err := validatorPK.Unmarshal(depositData.GetPubkey())
	if err != nil {
		return err
	}
	valConsAddr := sdk.GetConsAddress(validatorPK)
	validator, err := k.stakingKeeper.GetValidator(ctx, sdk.ValAddress(valConsAddr))
	amount := depositData.GetAmount()
	if err != nil {
		if errors.Is(err, stakingtypes.ErrNoValidatorFound) {
			_, err = k.createValidator(ctx, validatorPK, amount)
			return err
		}
		return err
	}
	_, err = k.stakingKeeper.Delegate(ctx, sdk.AccAddress(valConsAddr), math.NewIntFromUint64(amount), stakingtypes.Unbonded, validator, true)
	return err
}

// createValidator creates a new validator with the given public key and amount of tokens.
func (k Keeper) createValidator(ctx context.Context, validatorPK cryptotypes.PubKey, amount uint64) (stakingtypes.Validator, error) {
	stake := math.NewIntFromUint64(amount)
	valConsAddr := sdk.GetConsAddress(validatorPK)
	operator := sdk.ValAddress(valConsAddr).String()
	val, err := stakingtypes.NewValidator(operator, validatorPK, stakingtypes.Description{Moniker: operator})
	val.Tokens = stake
	val.DelegatorShares = math.LegacyNewDecFromInt(val.Tokens)
	return val, err
}
