// SPDX-License-IDentifier: MIT
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

package backend

import (
	"context"
	"strconv"

	types "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	serverType "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func (h Backend) GetGenesis(ctx context.Context) (primitives.Root, error) {
	// needs genesis_time and gensis_fork_version
	return h.getNewStateDB(ctx, "stateID").GetGenesisValidatorsRoot()
}

func (h Backend) GetStateRoot(
	ctx context.Context,
	stateID string,
) (primitives.Bytes32, error) {
	stateDB := h.getNewStateDB(ctx, stateID)
	slot, err := stateDB.GetSlot()
	if err != nil {
		return primitives.Bytes32{}, err
	}
	block, err := stateDB.StateRootAtIndex(slot.Unwrap())
	if err != nil {
		return primitives.Bytes32{}, err
	}
	root, err := block.HashTreeRoot()
	if err != nil {
		return primitives.Bytes32{}, err
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
) (primitives.Bytes32, error) {
	stateDB := h.getNewStateDB(ctx, "stateID")
	slot, err := stateDB.GetSlot()
	if err != nil {
		return primitives.Bytes32{}, err
	}
	block, err := stateDB.GetBlockRootAtIndex(slot.Unwrap())
	if err != nil {
		return primitives.Bytes32{}, err
	}
	root, err := block.HashTreeRoot()
	if err != nil {
		return primitives.Bytes32{}, err
	}
	return root, nil
}
