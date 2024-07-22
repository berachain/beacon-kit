// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package backend

import (
	"context"
	"strconv"

	"github.com/berachain/beacon-kit/mod/node-api/backend/storage"
	types "github.com/berachain/beacon-kit/mod/node-api/server/types/beacon"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// GetGenesis returns the genesis state of the beacon chain.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetGenesis(ctx context.Context) (common.Root, error) {
	// needs genesis_time and gensis_fork_version
	st, err := b.StateFromContext(ctx, storage.StateIDGenesis)
	if err != nil {
		return common.Root{}, err
	}
	return st.GetGenesisValidatorsRoot()
}

// GetStateRoot returns the root of the state at the given stateID.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetStateRoot(
	ctx context.Context,
	stateID string,
) (common.Bytes32, error) {
	st, err := b.StateFromContext(ctx, stateID)
	if err != nil {
		return common.Bytes32{}, err
	}
	slot, err := st.GetSlot()
	if err != nil {
		return common.Bytes32{}, err
	}
	// As calculated by the beacon chain. Ideally, this logic
	// should be abstracted by the beacon chain.
	index := slot.Unwrap() % b.cs.SlotsPerHistoricalRoot()
	return st.StateRootAtIndex(index)
}

// GetStateFork returns the fork of the state at the given stateID.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, ForkT, _, _, _, _, _,
]) GetStateFork(
	ctx context.Context,
	stateID string,
) (ForkT, error) {
	var fork ForkT
	st, err := b.StateFromContext(ctx, stateID)
	if err != nil {
		return fork, err
	}
	return st.GetFork()
}

// GetBlockRoot returns the root of the block at the given stateID.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetBlockRoot(
	ctx context.Context,
	stateID string,
) (common.Bytes32, error) {
	st, err := b.StateFromContext(ctx, stateID)
	if err != nil {
		return common.Bytes32{}, err
	}
	slot, err := st.GetSlot()
	if err != nil {
		return common.Bytes32{}, err
	}
	return st.GetBlockRootAtIndex(slot.Unwrap())
}

// TODO: Implement this.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetBlockRewards(
	_ context.Context,
	_ string,
) (*types.BlockRewardsData, error) {
	return &types.BlockRewardsData{
		ProposerIndex:     1,
		Total:             1,
		Attestations:      1,
		SyncAggregate:     1,
		ProposerSlashings: 1,
		AttesterSlashings: 1,
	}, nil
}

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _,
]) GetStateValidator(
	ctx context.Context,
	stateID string,
	validatorID string,
) (*types.ValidatorData[ValidatorT], error) {
	st, err := b.StateFromContext(ctx, stateID)
	if err != nil {
		return nil, err
	}
	index, indexErr := b.getValidatorIndex(st, validatorID)
	if indexErr != nil {
		return nil, indexErr
	}
	validator, validatorErr := st.ValidatorByIndex(index)
	if validatorErr != nil {
		return nil, validatorErr
	}
	balance, balanceErr := st.GetBalance(index)
	if balanceErr != nil {
		return nil, balanceErr
	}
	return &types.ValidatorData[ValidatorT]{
		Index:     index.Unwrap(),
		Balance:   balance.Unwrap(),
		Status:    "active",
		Validator: validator,
	}, nil
}

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _,
]) GetStateValidators(
	ctx context.Context,
	stateID string,
	id []string,
	_ []string,
) ([]*types.ValidatorData[ValidatorT], error) {
	validators := make([]*types.ValidatorData[ValidatorT], 0)
	for _, indexOrKey := range id {
		validatorData, err := b.GetStateValidator(ctx, stateID, indexOrKey)
		if err != nil {
			return nil, err
		}
		validators = append(validators, validatorData)
	}
	return validators, nil
}

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetStateValidatorBalances(
	ctx context.Context,
	stateID string,
	id []string,
) ([]*types.ValidatorBalanceData, error) {
	st, err := b.StateFromContext(ctx, stateID)
	if err != nil {
		return nil, err
	}
	balances := make([]*types.ValidatorBalanceData, 0)
	for _, indexOrKey := range id {
		index, indexErr := b.getValidatorIndex(st, indexOrKey)
		if indexErr != nil {
			return nil, indexErr
		}
		var balance math.U64
		balance, err = st.GetBalance(index)
		if err != nil {
			return nil, err
		}
		balances = append(balances, &types.ValidatorBalanceData{
			Index:   index.Unwrap(),
			Balance: balance.Unwrap(),
		})
	}
	return balances, nil
}

func (b Backend[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _,
]) getValidatorIndex(st BeaconStateT, keyOrIndex string) (math.U64, error) {
	if index, err := strconv.ParseUint(keyOrIndex, 10, 64); err == nil {
		return math.U64(index), nil
	}
	var key crypto.BLSPubkey
	if err := key.UnmarshalText([]byte(keyOrIndex)); err != nil {
		return math.U64(0), err
	}
	return st.ValidatorIndexByPubkey(key)
}
