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
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// InitializeBeaconStateFromEth1 initializes the beacon state. Modified from the ETH 2.0 spec:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#genesis
func (sp *StateProcessor) InitializeBeaconStateFromEth1(
	st *statedb.StateDB,
	deposits ctypes.Deposits,
	execPayloadHeader *ctypes.ExecutionPayloadHeader,
	genesisVersion common.Version,
) (transition.ValidatorUpdates, error) {
	if err := st.SetSlot(constants.GenesisSlot); err != nil {
		return nil, err
	}

	fork := ctypes.NewFork(genesisVersion, genesisVersion, constants.GenesisEpoch)
	if err := st.SetFork(fork); err != nil {
		return nil, err
	}
	if err := sp.ProcessFork(st, execPayloadHeader.GetTimestamp(), true); err != nil {
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

	if err := sp.seedRandaoMix(
		st,
		execPayloadHeader.GetBlockHash(),
	); err != nil {
		return nil, err
	}

	// ingest deposits & do genesis‐activation
	if err := sp.processGenesisDepositsAndActivations(st, deposits); err != nil {
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

	// seed historical block‑ and state‑roots
	if err = sp.seedHistoricalRoots(st); err != nil {
		return nil, err
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

// seedRandaoMix writes the initial RANDAO mixes.
func (sp *StateProcessor) seedRandaoMix(
	st *statedb.StateDB,
	hash common.ExecutionHash,
) error {
	for i := range sp.cs.EpochsPerHistoricalVector() {
		if err := st.UpdateRandaoMixAtIndex(
			i, common.Bytes32(hash),
		); err != nil {
			return err
		}
	}
	return nil
}

// seedHistoricalRoots zero‑primes the block and state‐roots.
func (sp *StateProcessor) seedHistoricalRoots(st *statedb.StateDB) error {
	for i := range sp.cs.HistoricalRootsLimit() {
		if err := st.UpdateBlockRootAtIndex(i, common.Root{}); err != nil {
			return err
		}
		if err := st.UpdateStateRootAtIndex(i, common.Root{}); err != nil {
			return err
		}
	}
	return nil
}

// processGenesisDepositsAndActivations handles the eth1 deposit index,
// validates and ingests each deposit, then does the genesis activation pass.
func (sp *StateProcessor) processGenesisDepositsAndActivations(
	st *statedb.StateDB,
	deposits ctypes.Deposits,
) error {
	// Before processing deposits, set the eth1 deposit index to 0.
	if err := st.SetEth1DepositIndex(constants.FirstDepositIndex); err != nil {
		return err
	}
	if err := validateGenesisDeposits(
		st, deposits, sp.cs.ValidatorSetCap(),
	); err != nil {
		return err
	}
	for _, dep := range deposits {
		if err := sp.processDeposit(st, dep); err != nil {
			return err
		}
	}
	return sp.processGenesisActivation(st)
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
