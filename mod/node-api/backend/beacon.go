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
	"github.com/berachain/beacon-kit/mod/node-api/backend/utils"
	types "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// GetGenesis returns the genesis state of the beacon chain.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GenesisValidatorsRoot(
	slot uint64,
) (common.Root, error) {
	// needs genesis_time and gensis_fork_version
	st, err := b.StateFromSlot(slot)
	if err != nil {
		return common.Root{}, err
	}
	return st.GetGenesisValidatorsRoot()
}

// GetStateRoot returns the root of the state at the given stateID.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) StateRootAtSlot(
	slot uint64,
) (common.Root, error) {
	st, err := b.StateFromSlot(slot)
	if err != nil {
		return common.Root{}, err
	}
	// This is required to handle the semantical expectation that
	// 0 -> latest despite 0 != latest.
	latestSlot, err := st.GetSlot()
	if err != nil {
		return common.Root{}, err
	}
	// As calculated by the beacon chain. Ideally, this logic
	// should be abstracted by the beacon chain.
	index := latestSlot.Unwrap() % b.cs.SlotsPerHistoricalRoot()
	return st.StateRootAtIndex(index)
}

// GetStateFork returns the fork of the state at the given stateID.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, ForkT, _, _, _, _, _, _,
]) StateForkAtSlot(
	slot uint64,
) (ForkT, error) {
	var fork ForkT
	st, err := b.StateFromSlot(slot)
	if err != nil {
		return fork, err
	}
	return st.GetFork()
}

// GetBlockRoot returns the root of the block at the given stateID.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) BlockRootAtSlot(
	slot uint64,
) (common.Root, error) {
	st, err := b.StateFromSlot(slot)
	if err != nil {
		return common.Root{}, err
	}
	latestSlot, err := st.GetSlot()
	if err != nil {
		return common.Root{}, err
	}
	return st.GetBlockRootAtIndex(latestSlot.Unwrap())
}

// TODO: Implement this.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) BlockRewardsAtSlot(
	_ uint64,
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
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _,
]) ValidatorByID(
	slot uint64,
	id string,
) (*types.ValidatorData[ValidatorT], error) {
	// TODO: to adhere to the spec, this shouldn't error if the error
	// is not found, but i can't think of a way to do that without coupling
	// db impl to the api impl.
	st, err := b.StateFromSlot(slot)
	if err != nil {
		return nil, err
	}
	index, err := utils.ValidatorIndexByID(st, id)
	if err != nil {
		return nil, err
	}
	validator, err := st.ValidatorByIndex(index)
	if err != nil {
		return nil, err
	}
	balance, err := st.GetBalance(index)
	if err != nil {
		return nil, err
	}
	return &types.ValidatorData[ValidatorT]{
		ValidatorBalanceData: types.ValidatorBalanceData{
			Index:   index.Unwrap(),
			Balance: balance.Unwrap(),
		},
		Status:    "active", // TODO: fix
		Validator: validator,
	}, nil
}

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _,
]) ValidatorsByIDs(
	slot uint64,
	ids []string,
	_ []string, // TODO: filter by status
) ([]*types.ValidatorData[ValidatorT], error) {
	validatorsData := make([]*types.ValidatorData[ValidatorT], 0)
	for _, id := range ids {
		// TODO: we can probably optimize this via a getAllValidators
		// query and then filtering but blocked by the fact that IDs
		// can be indices and the hard type only holds its own pubkey.
		validatorData, err := b.ValidatorByID(slot, id)
		if err != nil {
			return nil, err
		}
		validatorsData = append(validatorsData, validatorData)
	}
	return validatorsData, nil
}

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) ValidatorBalancesByIDs(
	slot uint64,
	ids []string,
) ([]*types.ValidatorBalanceData, error) {
	var index math.U64
	st, err := b.StateFromSlot(slot)
	if err != nil {
		return nil, err
	}
	balances := make([]*types.ValidatorBalanceData, 0)
	for _, id := range ids {
		index, err = utils.ValidatorIndexByID(st, id)
		if err != nil {
			return nil, err
		}
		var balance math.U64
		// TODO: same issue as above, shouldn't error on not found.
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
