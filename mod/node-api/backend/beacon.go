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
	sszTypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	response "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func (h Backend) GetGenesis(
	ctx context.Context,
) (*response.GenesisData, error) {
	return h.getNewStateDB(ctx, "genesis").GetGenesisDetails()
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
) ([]*response.ValidatorData, error) {
	stateDB := h.getNewStateDB(ctx, stateID)
	validators := make([]*response.ValidatorData, 0)
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
		validators = append(validators, &response.ValidatorData{
			Index:     index.Unwrap(),
			Balance:   balance.Unwrap(),
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
) (*response.ValidatorData, error) {
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
	return &response.ValidatorData{
		Index:     index.Unwrap(),
		Balance:   balance.Unwrap(),
		Validator: validator,
	}, nil
}

func (h Backend) GetStateValidatorBalances(
	ctx context.Context,
	stateID string,
	id []string,
) ([]*response.ValidatorBalanceData, error) {
	stateDB := h.getNewStateDB(ctx, stateID)
	balances := make([]*response.ValidatorBalanceData, 0)
	for _, indexOrKey := range id {
		index, indexErr := getValidatorIndex(stateDB, indexOrKey)
		if indexErr != nil {
			return nil, indexErr
		}
		balance, err := stateDB.GetBalance(index)
		if err != nil {
			return nil, err
		}
		balances = append(balances, &response.ValidatorBalanceData{
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

func (h Backend) GetStateCommittees(
	ctx context.Context,
	stateID string,
	_ string,
	epoch string,
	_ string,
) ([]*response.CommitteeData, error) {
	stateDB := h.getNewStateDB(ctx, stateID)
	epochU64, epochErr := resolveEpoch(stateDB, epoch)
	if epochErr != nil {
		return nil, epochErr
	}
	committees, committeeErr := stateDB.GetStateCommittees(epochU64)
	if committeeErr != nil {
		return nil, committeeErr
	}
	return committees, nil
}

func (h Backend) GetStateSyncCommittees(
	ctx context.Context,
	stateID string,
	epoch string,
) (*response.SyncCommitteeData, error) {
	stateDB := h.getNewStateDB(ctx, stateID)
	epochU64, epochErr := resolveEpoch(stateDB, epoch)
	if epochErr != nil {
		return nil, epochErr
	}
	committees, committeeErr := stateDB.GetStateSyncCommittees(epochU64)
	if committeeErr != nil {
		return nil, committeeErr
	}
	// filter committees by index and slot
	return committees, nil
}

func (h Backend) GetBlockHeaders(
	ctx context.Context,
	slot string,
	parentRoot primitives.Root,
) ([]*response.BlockHeaderData, error) {
	blockDB := h.getNewBlockDB(ctx, "blockID")
	slotU64, slotErr := resolveUint64[math.Slot](slot)
	if slotErr != nil {
		return nil, slotErr
	}
	headers, err := blockDB.GetBlockHeaders(slotU64, parentRoot)
	if err != nil {
		return nil, err
	}
	return headers, nil
}

func (h Backend) GetBlockHeader(
	ctx context.Context,
	blockID string,
) (*response.BlockHeaderData, error) {
	blockDB := h.getNewBlockDB(ctx, blockID)
	header, err := blockDB.GetBlockHeader()
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (h Backend) GetBlock(
	ctx context.Context,
	blockID string,
) (*types.BeaconBlock, error) {
	blockDB := h.getNewBlockDB(ctx, blockID)
	block, err := blockDB.GetBlock()
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (h Backend) GetBlockBlobSidecars(
	ctx context.Context,
	blockID string,
	indicies []string,
) ([]*sszTypes.BlobSidecar, error) {
	blockDB := h.getNewBlockDB(ctx, blockID)
	sidecars, err := blockDB.GetBlockBlobSidecars(indicies)
	if err != nil {
		return nil, err
	}
	return sidecars, nil
}

func (h Backend) GetPoolVoluntaryExits(
	ctx context.Context,
) ([]*response.MessageSignature, error) {
	nodeState := h.getNodeState(ctx)
	exits, err := nodeState.GetVoluntaryExits()
	if err != nil {
		return nil, err
	}
	return exits, nil
}

func (h Backend) GetPoolBtsToExecutionChanges(
	ctx context.Context,
) ([]*response.MessageSignature, error) {
	nodeState := h.getNodeState(ctx)
	blsToExecutionChanges, err := nodeState.GetBlsToExecutionChanges()
	if err != nil {
		return nil, err
	}
	return blsToExecutionChanges, nil
}
