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

	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// InitializePreminedBeaconStateFromEth1 initializes the beacon state.
//
//nolint:gocognit,funlen // todo fix.
func (sp *StateProcessor[
	_, _,
]) InitializePreminedBeaconStateFromEth1(
	st *statedb.StateDB,
	deposits ctypes.Deposits,
	execPayloadHeader *ctypes.ExecutionPayloadHeader,
	genesisVersion common.Version,
) (transition.ValidatorUpdates, error) {
	if err := st.SetSlot(0); err != nil {
		return nil, err
	}

	fork := ctypes.NewFork(
		genesisVersion,
		genesisVersion,
		math.U64(constants.GenesisEpoch),
	)
	if err := st.SetFork(fork); err != nil {
		return nil, err
	}

	var eth1Data *ctypes.Eth1Data
	eth1Data = eth1Data.New(
		deposits.HashTreeRoot(),
		0,
		execPayloadHeader.GetBlockHash(),
	)
	if err := st.SetEth1Data(eth1Data); err != nil {
		return nil, err
	}

	// TODO: we need to handle common.Version vs uint32 better.
	var blkBody *ctypes.BeaconBlockBody
	blkBody = blkBody.Empty(version.ToUint32(genesisVersion))

	var blkHeader *ctypes.BeaconBlockHeader
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

	// Before processing deposits, set the eth1 deposit index to 0.
	if err := st.SetEth1DepositIndex(0); err != nil {
		return nil, err
	}
	if err := sp.validateGenesisDeposits(st, deposits); err != nil {
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

	// Handle special case bartio genesis.
	validatorsRoot := common.Root(hex.MustToBytes(spec.BartioValRoot))
	if sp.cs.DepositEth1ChainID() != spec.BartioChainID {
		validators, err := st.GetValidators()
		if err != nil {
			return nil, err
		}
		validatorsRoot = validators.HashTreeRoot()
	}
	if err := st.SetGenesisValidatorsRoot(validatorsRoot); err != nil {
		return nil, err
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

	activeVals, err := getActiveVals(sp.cs, st, 0)
	if err != nil {
		return nil, err
	}
	return validatorSetsDiffs(nil, activeVals), nil
}

func (sp *StateProcessor[
	_, _,
]) processGenesisActivation(
	st *statedb.StateDB,
) error {
	switch {
	case sp.cs.DepositEth1ChainID() == spec.BartioChainID:
		// nothing to do
		return nil
	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID:
		// nothing to do
		return nil
	default:
		vals, err := st.GetValidators()
		if err != nil {
			return fmt.Errorf(
				"genesis activation, failed listing validators: %w",
				err,
			)
		}
		minEffectiveBalance := math.Gwei(
			sp.cs.EjectionBalance() + sp.cs.EffectiveBalanceIncrement(),
		)

		var idx math.ValidatorIndex
		for _, val := range vals {
			if val.GetEffectiveBalance() < minEffectiveBalance {
				continue
			}
			val.SetActivationEligibilityEpoch(0)
			val.SetActivationEpoch(0)
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
}
