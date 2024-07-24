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

	types "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types/v2"
	serverType "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func (h Backend) GetGenesis(ctx context.Context) (common.Root, error) {
	// needs genesis_time and gensis_fork_version
	return h.getNewStateDB(ctx, "stateID").GetGenesisValidatorsRoot()
}

func (h Backend) GetStateRoot(
	ctx context.Context,
	stateID string,
) (common.Bytes32, error) {
	stateDB := h.getNewStateDB(ctx, stateID)
	slot, err := stateDB.GetSlot()
	if err != nil {
		return common.Bytes32{}, err
	}
	block, err := stateDB.StateRootAtIndex(slot.Unwrap())
	if err != nil {
		return common.Bytes32{}, err
	}
	root, err := block.HashTreeRoot()
	if err != nil {
		return common.Bytes32{}, err
	}
	return root, nil
}

func (h Backend) GetStateFork(
	ctx context.Context,
	stateID string,
) (*types.Fork, error) {
	return h.getNewStateDB(ctx, stateID).GetFork()
}

func (h Backend) GetStateValidators(
	ctx context.Context,
	stateID string,
	id []string,
	_ []string,
) ([]*serverType.ValidatorData, error) {
	stateDB := h.getNewStateDB(ctx, stateID)
	validators := make([]*serverType.ValidatorData, 0)
	for _, indexOrKey := range id {
		index, indexErr := getValidatorIndex(stateDB, indexOrKey)
		if indexErr != nil {
			return nil, indexErr
		}
		validator, validatorErr := stateDB.ValidatorByIndex(index)
		if validatorErr != nil {
			return nil, validatorErr
		}
		balance, balanceErr := stateDB.GetBalance(index)
		if balanceErr != nil {
			return nil, balanceErr
		}
		validators = append(validators, &serverType.ValidatorData{
			Index:     index.Unwrap(),
			Balance:   balance.Unwrap(),
			Status:    "active",
			Validator: validator,
		})
	}
	return validators, nil
}

func getValidatorIndex(stateDB StateDB, keyOrIndex string) (math.U64, error) {
	if index, err := strconv.ParseUint(keyOrIndex, 10, 64); err == nil {
		return math.U64(index), nil
	}
	key := crypto.BLSPubkey{}
	err := key.UnmarshalText([]byte(keyOrIndex))
	if err != nil {
		return math.U64(0), err
	}
	index, err := stateDB.ValidatorIndexByPubkey(key)
	if err == nil {
		return index, nil
	}
	return math.U64(0), err
}

func (h Backend) GetStateValidator(
	ctx context.Context,
	stateID string,
	validatorID string,
) (*serverType.ValidatorData, error) {
	stateDB := h.getNewStateDB(ctx, stateID)
	index, indexErr := getValidatorIndex(stateDB, validatorID)
	if indexErr != nil {
		return nil, indexErr
	}
	validator, validatorErr := stateDB.ValidatorByIndex(index)
	if validatorErr != nil {
		return nil, validatorErr
	}
	balance, balanceErr := stateDB.GetBalance(index)
	if balanceErr != nil {
		return nil, balanceErr
	}
	return &serverType.ValidatorData{
		Index:     index.Unwrap(),
		Balance:   balance.Unwrap(),
		Status:    "active",
		Validator: validator,
	}, nil
}

func (h Backend) GetStateValidatorBalances(
	ctx context.Context,
	stateID string,
	id []string,
) ([]*serverType.ValidatorBalanceData, error) {
	stateDB := h.getNewStateDB(ctx, stateID)
	balances := make([]*serverType.ValidatorBalanceData, 0)
	for _, indexOrKey := range id {
		index, indexErr := getValidatorIndex(stateDB, indexOrKey)
		if indexErr != nil {
			return nil, indexErr
		}
		balance, err := stateDB.GetBalance(index)
		if err != nil {
			return nil, err
		}
		balances = append(balances, &serverType.ValidatorBalanceData{
			Index:   index.Unwrap(),
			Balance: balance.Unwrap(),
		})
	}
	return balances, nil
}

func (h Backend) GetBlockRoot(
	ctx context.Context,
	_ string,
) (common.Bytes32, error) {
	stateDB := h.getNewStateDB(ctx, "stateID")
	slot, err := stateDB.GetSlot()
	if err != nil {
		return common.Bytes32{}, err
	}
	block, err := stateDB.GetBlockRootAtIndex(slot.Unwrap())
	if err != nil {
		return common.Bytes32{}, err
	}
	root, err := block.HashTreeRoot()
	if err != nil {
		return common.Bytes32{}, err
	}
	return root, nil
}
