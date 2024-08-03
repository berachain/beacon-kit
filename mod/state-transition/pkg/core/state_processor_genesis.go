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

package core

import (
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// InitializePreminedBeaconStateFromEth1 initializes the beacon state.
//
//nolint:gocognit,funlen // todo fix.
func (sp *StateProcessor[
	_, BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT, _, DepositT,
	Eth1DataT, _, ExecutionPayloadHeaderT, ForkT, _, _, ValidatorT, _, _,
]) InitializePreminedBeaconStateFromEth1(
	st BeaconStateT,
	deposits []DepositT,
	executionPayloadHeader ExecutionPayloadHeaderT,
	genesisVersion common.Version,
) (transition.ValidatorUpdates, error) {
	var (
		blkHeader BeaconBlockHeaderT
		blkBody   BeaconBlockBodyT
		fork      ForkT
		eth1Data  Eth1DataT
	)
	fork = fork.New(
		genesisVersion,
		genesisVersion,
		math.U64(constants.GenesisEpoch),
	)

	if err := st.SetSlot(0); err != nil {
		return nil, err
	}

	if err := st.SetFork(fork); err != nil {
		return nil, err
	}

	if err := st.SetEth1DepositIndex(0); err != nil {
		return nil, err
	}

	if err := st.SetEth1Data(eth1Data.New(
		common.Bytes32(gethprimitives.ZeroHash),
		0,
		executionPayloadHeader.GetBlockHash(),
	)); err != nil {
		return nil, err
	}

	// TODO: we need to handle common.Version vs
	// uint32 better.
	bodyRoot, err := blkBody.Empty(
		version.ToUint32(genesisVersion)).HashTreeRoot()
	if err != nil {
		return nil, err
	}

	if err = st.SetLatestBlockHeader(blkHeader.New(
		0, 0, common.Root{}, common.Root{}, bodyRoot,
	)); err != nil {
		return nil, err
	}

	for i := range sp.cs.EpochsPerHistoricalVector() {
		if err = st.UpdateRandaoMixAtIndex(
			i,
			common.Bytes32(executionPayloadHeader.GetBlockHash()),
		); err != nil {
			return nil, err
		}
	}

	for _, deposit := range deposits {
		// TODO: process deposits into eth1 data.
		if err = sp.processDeposit(st, deposit); err != nil {
			return nil, err
		}
	}

	// TODO: process activations.
	var validators []ValidatorT
	validators, err = st.GetValidators()
	if err != nil {
		return nil, err
	}

	var validatorsRoot common.Root

	validatorsRoot, err = merkle.
		NewMerkleizer[[32]byte, ValidatorT]().MerkleizeListComposite(
		validators,
		uint64(len(validators)),
	)
	if err != nil {
		return nil, err
	}

	if err = st.SetGenesisValidatorsRoot(validatorsRoot); err != nil {
		return nil, err
	}

	if err = st.SetLatestExecutionPayloadHeader(
		executionPayloadHeader,
	); err != nil {
		return nil, err
	}

	// Setup a bunch of 0s to prime the DB.
	for i := range sp.cs.HistoricalRootsLimit() {
		//#nosec:G701 // won't overflow in practice.
		if err = st.UpdateBlockRootAtIndex(i, common.Root{}); err != nil {
			return nil, err
		}
		if err = st.UpdateStateRootAtIndex(i, common.Root{}); err != nil {
			return nil, err
		}
	}

	if err = st.SetNextWithdrawalIndex(0); err != nil {
		return nil, err
	}

	if err = st.SetNextWithdrawalValidatorIndex(
		0,
	); err != nil {
		return nil, err
	}

	if err = st.SetTotalSlashing(0); err != nil {
		return nil, err
	}

	var updates transition.ValidatorUpdates
	updates, err = sp.processSyncCommitteeUpdates(st)
	if err != nil {
		return nil, err
	}
	st.Save()
	return updates, nil
}
