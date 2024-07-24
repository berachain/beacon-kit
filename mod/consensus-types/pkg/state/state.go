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

package state

import (
	"reflect"

	types "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types/v2"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// BeaconState is the interface for the beacon state.
type BeaconStateMarshallable[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT any,
] struct {
	// TODO: decouple from deneb.BeaconState
	*types.BeaconState
}

// New creates a new BeaconState.
func (st *BeaconStateMarshallable[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT,
]) New(
	forkVersion uint32,
	genesisValidatorsRoot common.Root,
	slot math.Slot,
	fork ForkT,
	latestBlockHeader BeaconBlockHeaderT,
	blockRoots []common.Root,
	stateRoots []common.Root,
	eth1Data Eth1DataT,
	eth1DepositIndex uint64,
	latestExecutionPayloadHeader ExecutionPayloadHeaderT,
	validators []ValidatorT,
	balances []uint64,
	randaoMixes []common.Bytes32,
	nextWithdrawalIndex uint64,
	nextWithdrawalValidatorIndex math.ValidatorIndex,
	slashings []uint64,
	totalSlashing math.Gwei,
) (*BeaconStateMarshallable[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT,
], error) {
	switch forkVersion {
	case version.Deneb, version.DenebPlus:
		return &BeaconStateMarshallable[
			BeaconBlockHeaderT,
			Eth1DataT,
			ExecutionPayloadHeaderT,
			ForkT,
			ValidatorT,
		]{
			// TODO: Unhack reflection.
			BeaconState: &types.BeaconState{
				Slot:                  slot,
				GenesisValidatorsRoot: genesisValidatorsRoot,
				Fork: reflect.ValueOf(fork).
					Interface().(*types.Fork),
				LatestBlockHeader: reflect.ValueOf(latestBlockHeader).
					Interface().(*types.BeaconBlockHeader),
				BlockRoots: blockRoots,
				StateRoots: stateRoots,
				LatestExecutionPayloadHeader: reflect.
					ValueOf(latestExecutionPayloadHeader).
					Interface().(*types.ExecutionPayloadHeader),
				Eth1Data: reflect.ValueOf(eth1Data).
					Interface().(*types.Eth1Data),
				Eth1DepositIndex: math.U64(eth1DepositIndex),
				Validators: reflect.ValueOf(validators).
					Interface().([]*types.Validator),
				Balances:                     balances,
				RandaoMixes:                  randaoMixes,
				NextWithdrawalIndex:          math.U64(nextWithdrawalIndex),
				NextWithdrawalValidatorIndex: nextWithdrawalValidatorIndex,
				Slashings:                    slashings,
				TotalSlashing:                totalSlashing,
			},
		}, nil
	default:
		return nil, errors.Wrapf(ErrUnsupportedVersion, "%d", forkVersion)
	}
}
