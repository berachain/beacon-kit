// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// InitializePreminedBeaconStateFromEth1 initializes the beacon state.
//
//nolint:gocognit,funlen // todo fix.
func (sp *StateProcessor) InitializePreminedBeaconStateFromEth1(
	st *statedb.StateDB,
	deposits ctypes.Deposits,
	execPayloadHeader *ctypes.ExecutionPayloadHeader,
	genesisVersion common.Version,
) (transition.ValidatorUpdates, error) {
	if !version.Equals(genesisVersion, execPayloadHeader.GetForkVersion()) {
		return nil, fmt.Errorf("fork mismatch between genesisVersion %s and execPayloadHeader %s",
			genesisVersion, execPayloadHeader.GetForkVersion())
	}
	if !version.Equals(sp.cs.GenesisForkVersion(), genesisVersion) {
		return nil, fmt.Errorf("fork mismatch between chain spec genesis version %s and genesisVersion %s",
			sp.cs.GenesisForkVersion(), genesisVersion)
	}

	if err := st.SetSlot(constants.GenesisSlot); err != nil {
		return nil, err
	}

	fork := ctypes.NewFork(genesisVersion, genesisVersion, constants.GenesisEpoch)
	if err := st.SetFork(fork); err != nil {
		return nil, err
	}
	if err := sp.PrepareStateForFork(st, genesisVersion, 0); err != nil {
		return nil, err
	}

	eth1Data := &ctypes.Eth1Data{
		DepositRoot:  deposits.HashTreeRoot(),
		DepositCount: 0,
		BlockHash:    execPayloadHeader.GetBlockHash(),
	}
	if err := st.SetEth1Data(eth1Data); err != nil {
		return nil, err
	}

	versionable := ctypes.NewVersionable(genesisVersion)
	blkBody := &ctypes.BeaconBlockBody{
		Versionable: versionable,
		Eth1Data:    &ctypes.Eth1Data{},
		ExecutionPayload: &ctypes.ExecutionPayload{
			Versionable: versionable,
			ExtraData:   make([]byte, ctypes.ExtraDataSize),
		},
	}

	blkHeader := &ctypes.BeaconBlockHeader{
		Slot:            constants.GenesisSlot,
		ProposerIndex:   0,
		ParentBlockRoot: common.Root{},
		StateRoot:       common.Root{},
		BodyRoot:        blkBody.HashTreeRoot(),
	}
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

	// Before processing deposits, set the eth1 deposit index to 0.
	if err := st.SetEth1DepositIndex(constants.FirstDepositIndex); err != nil {
		return nil, err
	}
	if err := validateGenesisDeposits(st, deposits, sp.cs.ValidatorSetCap()); err != nil {
		return nil, err
	}
	for _, deposit := range deposits {
		if err := sp.processDeposit(st, deposit); err != nil {
			return nil, err
		}
	}

	// process activations
	if err := sp.processGenesisActivation(st); err != nil {
		return nil, err
	}

	validators, err := st.GetValidators()
	if err != nil {
		return nil, err
	}
	if err = st.SetGenesisValidatorsRoot(validators.HashTreeRoot()); err != nil {
		return nil, err
	}

	if err = st.SetLatestExecutionPayloadHeader(execPayloadHeader); err != nil {
		return nil, err
	}

	// Setup a bunch of 0s to prime the DB.
	for i := range sp.cs.HistoricalRootsLimit() {
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

	if err = st.SetNextWithdrawalValidatorIndex(0); err != nil {
		return nil, err
	}

	if err = st.SetTotalSlashing(0); err != nil {
		return nil, err
	}

	activeVals, err := getActiveVals(st, constants.GenesisEpoch)
	if err != nil {
		return nil, err
	}
	return validatorSetsDiffs(nil, activeVals), nil
}

func (sp *StateProcessor) processGenesisActivation(st *statedb.StateDB) error {
	vals, err := st.GetValidators()
	if err != nil {
		return fmt.Errorf("genesis activation, failed listing validators: %w", err)
	}
	minEffectiveBalance := math.Gwei(
		sp.cs.EjectionBalance() +
			sp.cs.EffectiveBalanceIncrement(),
	)

	var idx math.ValidatorIndex
	for _, val := range vals {
		if val.GetEffectiveBalance() < minEffectiveBalance {
			continue
		}
		val.SetActivationEligibilityEpoch(constants.GenesisEpoch)
		val.SetActivationEpoch(constants.GenesisEpoch)
		idx, err = st.ValidatorIndexByPubkey(val.GetPubkey())
		if err != nil {
			return err
		}
		if err = st.UpdateValidatorAtIndex(idx, val); err != nil {
			return err
		}
	}

	return nil
}
