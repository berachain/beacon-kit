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

package types

import (
	"unsafe"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlockForValidator represents a block in the beacon chain during the
// Deneb fork, with only the minimally required values to prove a validator
// exists in this block.
type BeaconBlockForValidator struct {
	// Slot represents the position of the block in the chain.
	Slot math.Slot
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.ValidatorIndex
	// ParentBlockRoot is the hash of the parent block.
	ParentBlockRoot common.Root
	// StateRoot is the summary of the BeaconState type with only the required
	// raw values to prove a validator's pubkey.
	StateRoot *BeaconStateForValidator
	// BodyRoot is the root of the block body.
	BodyRoot common.Root
}

// NewBeaconBlockForValidator creates a new BeaconBlock SSZ summary with only
// the required raw values to prove a validator exists in this block.
func NewBeaconBlockForValidator[BeaconBlockHeaderT constraints.SSZRootable](
	bbh BeaconBlockHeader[BeaconBlockHeaderT],
	bsv *BeaconStateForValidator,
) (*BeaconBlockForValidator, error) {
	return &BeaconBlockForValidator{
		Slot:            bbh.GetSlot(),
		ProposerIndex:   bbh.GetProposerIndex(),
		ParentBlockRoot: bbh.GetParentBlockRoot(),
		StateRoot:       bsv,
		BodyRoot:        bbh.GetBodyRoot(),
	}, nil
}

// BeaconStateForValidator is the SSZ summary of the BeaconState type with only
// the required raw values to prove a validator exists in this state.
type BeaconStateForValidator struct {
	GenesisValidatorsRoot common.Root
	Slot                  math.Slot
	// Fork is the hash tree root of the Fork.
	Fork common.Root
	// LatestBlockHeader is the hash tree root of the latest block header.
	LatestBlockHeader common.Root
	// BlockRoots is the hash tree root of the block headers buffer.
	BlockRoots common.Root
	// StateRoots is the hash tree root of the beacon states buffer.
	StateRoots common.Root
	// Eth1Data is the hash tree root of the eth1 data.
	Eth1Data         common.Root
	Eth1DepositIndex uint64
	// LatestExecutionPayloadHeader is the hash tree root of the latest
	// execution payload header.
	LatestExecutionPayloadHeader common.Root
	// Validators is the list of Validators with the field Pubkey to prove.
	Validators []*types.Validator `ssz-max:"1099511627776"`
	// Balances is the hash tree root of the validator balances.
	Balances common.Root
	// RandaoMixes is the hash tree root of the randao mixes.
	RandaoMixes                  common.Root
	NextWithdrawalIndex          uint64
	NextWithdrawalValidatorIndex math.ValidatorIndex
	// Slashings is the hash tree root of the slashings.
	Slashings     common.Root
	TotalSlashing math.Gwei
}

// NewBeaconStateForValidator creates a new BeaconState SSZ summary with only
// the required raw values to prove a validator exists in this state.
//
//nolint:funlen,gocognit // all lines are required to pack the entire beacon
func NewBeaconStateForValidator[
	BeaconBlockHeaderT constraints.SSZRootable,
	BeaconStateT BeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		ValidatorT,
	],
	Eth1DataT constraints.SSZRootable,
	ExecutionPayloadHeaderT constraints.SSZRootable,
	ForkT constraints.SSZRootable,
	ValidatorT any,
](
	bs BeaconStateT,
	cs common.ChainSpec,
) (*BeaconStateForValidator, error) {
	var (
		bsv                       = &BeaconStateForValidator{}
		err                       error
		slotsPerHistoricalRoot    = cs.SlotsPerHistoricalRoot()
		epochsPerHistoricalVector = cs.EpochsPerHistoricalVector()
	)

	bsv.GenesisValidatorsRoot, err = bs.GetGenesisValidatorsRoot()
	if err != nil {
		return nil, err
	}

	if bsv.Slot, err = bs.GetSlot(); err != nil {
		return nil, err
	}

	var fork ForkT
	if fork, err = bs.GetFork(); err != nil {
		return nil, err
	}
	if bsv.Fork, err = fork.HashTreeRoot(); err != nil {
		return nil, err
	}

	var latestBlockHeader BeaconBlockHeaderT
	if latestBlockHeader, err = bs.GetLatestBlockHeader(); err != nil {
		return nil, err
	}
	bsv.LatestBlockHeader, err = latestBlockHeader.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	blockRoots := make([]common.Root, slotsPerHistoricalRoot)
	for i := range slotsPerHistoricalRoot {
		if blockRoots[i], err = bs.GetBlockRootAtIndex(i); err != nil {
			return nil, err
		}
	}
	if bsv.BlockRoots, err = ssz.ListFromElements(
		MaxBlockRoots, blockRoots...,
	).HashTreeRoot(); err != nil {
		return nil, err
	}

	stateRoots := make([]common.Root, slotsPerHistoricalRoot)
	for i := range slotsPerHistoricalRoot {
		if stateRoots[i], err = bs.StateRootAtIndex(i); err != nil {
			return nil, err
		}
	}
	if bsv.StateRoots, err = ssz.ListFromElements(
		MaxStateRoots, stateRoots...,
	).HashTreeRoot(); err != nil {
		return nil, err
	}

	var eth1Data Eth1DataT
	if eth1Data, err = bs.GetEth1Data(); err != nil {
		return nil, err
	}
	bsv.Eth1Data, err = eth1Data.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	if bsv.Eth1DepositIndex, err = bs.GetEth1DepositIndex(); err != nil {
		return nil, err
	}

	var leph ExecutionPayloadHeaderT
	leph, err = bs.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}
	if bsv.LatestExecutionPayloadHeader, err = leph.HashTreeRoot(); err != nil {
		return nil, err
	}

	var validators []ValidatorT
	if validators, err = bs.GetValidators(); err != nil {
		return nil, err
	}
	//#nosec:G103 // on purpose.
	bsv.Validators = *(*[]*types.Validator)(unsafe.Pointer(&validators))

	var balances []uint64
	if balances, err = bs.GetBalances(); err != nil {
		return nil, err
	}
	if bsv.Balances, err = ssz.ListFromElements(
		//#nosec:G103 // on purpose.
		MaxBalances, *(*[]math.U64)(unsafe.Pointer(&balances))...,
	).HashTreeRoot(); err != nil {
		return nil, err
	}

	randaoMixes := make([]common.Bytes32, epochsPerHistoricalVector)
	for i := range epochsPerHistoricalVector {
		if randaoMixes[i], err = bs.GetRandaoMixAtIndex(i); err != nil {
			return nil, err
		}
	}
	if bsv.RandaoMixes, err = ssz.ListFromElements(
		MaxRandaoMixes, randaoMixes...,
	).HashTreeRoot(); err != nil {
		return nil, err
	}

	if bsv.NextWithdrawalIndex, err = bs.GetNextWithdrawalIndex(); err != nil {
		return nil, err
	}

	bsv.NextWithdrawalValidatorIndex, err = bs.GetNextWithdrawalValidatorIndex()
	if err != nil {
		return nil, err
	}

	var slashings []uint64
	if slashings, err = bs.GetSlashings(); err != nil {
		return nil, err
	}
	if bsv.Slashings, err = ssz.ListFromElements(
		MaxSlashings, *(*[]math.U64)(unsafe.Pointer(&slashings))...,
	).HashTreeRoot(); err != nil {
		return nil, err
	}

	if bsv.TotalSlashing, err = bs.GetTotalSlashing(); err != nil {
		return nil, err
	}

	return bsv, nil
}
