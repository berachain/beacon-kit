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
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

//nolint:lll // temporary.
const (
	bArtioValRoot = "0x9147586693b6e8faa837715c0f3071c2000045b54233901c2e7871b15872bc43"
	bArtioChainID = 80084
)

// InitializePreminedBeaconStateFromEth1 initializes the beacon state.
//
//nolint:gocognit,funlen // todo fix.
func (sp *StateProcessor[
	_, BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT, _, DepositT,
	Eth1DataT, _, ExecutionPayloadHeaderT, ForkT, _, _, ValidatorT, _, _, _, _,
]) InitializePreminedBeaconStateFromEth1(
	st BeaconStateT,
	deposits []DepositT,
	execPayloadHeader ExecutionPayloadHeaderT,
	genesisVersion common.Version,
) (transition.ValidatorUpdates, error) {
	sp.processingGenesis = true
	defer func() {
		sp.processingGenesis = false
	}()

	if err := st.SetSlot(0); err != nil {
		return nil, err
	}

	var fork ForkT
	fork = fork.New(
		genesisVersion,
		genesisVersion,
		math.U64(constants.GenesisEpoch),
	)
	if err := st.SetFork(fork); err != nil {
		return nil, err
	}

	// Eth1DepositIndex will be set in processDeposit

	var eth1Data Eth1DataT
	eth1Data = eth1Data.New(
		common.Root{},
		0,
		execPayloadHeader.GetBlockHash(),
	)
	if err := st.SetEth1Data(eth1Data); err != nil {
		return nil, err
	}

	// TODO: we need to handle common.Version vs uint32 better.
	var blkBody BeaconBlockBodyT
	blkBody = blkBody.Empty(version.ToUint32(genesisVersion))

	var blkHeader BeaconBlockHeaderT
	blkHeader = blkHeader.New(
		0,                      // slot
		0,                      // proposer index
		common.Root{},          // parent block root
		common.Root{},          // state root
		blkBody.HashTreeRoot(), // body root

	)
	if err := st.SetLatestBlockHeader(blkHeader); err != nil {
		return nil, err
	}

	for i := range sp.cs.EpochsPerHistoricalVector() {
		if err := st.UpdateRandaoMixAtIndex(
			i,
			common.Bytes32(execPayloadHeader.GetBlockHash()),
		); err != nil {
			return nil, err
		}
	}

	// BeaconKit enforces a cap on the validator set size.
	// If genesis deposits breaches the cap we return an error.
	//#nosec:G701 // can't overflow.
	if uint32(len(deposits)) > sp.cs.GetValidatorSetCapSize() {
		return nil, fmt.Errorf("validator set cap %d, deposits count %d: %w",
			sp.cs.GetValidatorSetCapSize(),
			len(deposits),
			ErrHitValidatorsSetCap,
		)
	}

	// Process deposits
	for _, deposit := range deposits {
		if err := sp.processDeposit(st, deposit); err != nil {
			return nil, err
		}
	}

	// Currently we don't really process activations for validator.
	// We do not update ActivationEligibilityEpoch nor ActivationEpoch
	// for validators.
	// A validator is created with its EffectiveBalance duly set
	// (as in Eth 2.0 specs). The EffectiveBalance is updated at the
	// turn of the epoch, when the consensus is made aware of the
	// validator existence as well.
	// TODO: this is likely to change once we introduce a cap on
	// the validators set, in which case some validators may be evicted
	// from the validator set because the cap is reached.

	// Handle special case bartio genesis.
	if sp.cs.DepositEth1ChainID() == bArtioChainID {
		validatorsRoot := common.Root(hex.MustToBytes(bArtioValRoot))
		if err := st.SetGenesisValidatorsRoot(validatorsRoot); err != nil {
			return nil, err
		}
	} else {
		validators, err := st.GetValidators()
		if err != nil {
			return nil, err
		}
		if err = st.
			SetGenesisValidatorsRoot(validators.HashTreeRoot()); err != nil {
			return nil, err
		}
	}

	if err := st.SetLatestExecutionPayloadHeader(execPayloadHeader); err != nil {
		return nil, err
	}

	// Setup a bunch of 0s to prime the DB.
	for i := range sp.cs.HistoricalRootsLimit() {
		//#nosec:G701 // won't overflow in practice.
		if err := st.UpdateBlockRootAtIndex(i, common.Root{}); err != nil {
			return nil, err
		}
		if err := st.UpdateStateRootAtIndex(i, common.Root{}); err != nil {
			return nil, err
		}
	}

	if err := st.SetNextWithdrawalIndex(0); err != nil {
		return nil, err
	}

	if err := st.SetNextWithdrawalValidatorIndex(0); err != nil {
		return nil, err
	}

	if err := st.SetTotalSlashing(0); err != nil {
		return nil, err
	}

	return sp.processSyncCommitteeUpdates(st)
}
