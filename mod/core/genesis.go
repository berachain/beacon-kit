// SPDX-License-Identifier: MIT
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

package core

import (
	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/core/types"
	genutiltypes "github.com/berachain/beacon-kit/mod/node-builder/commands/genesis/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/version"
)

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *Genesis {
	return &Genesis{
		Fork: &primitives.Fork{
			PreviousVersion: version.FromUint32[primitives.Version](
				version.Deneb,
			),
			CurrentVersion: version.FromUint32[primitives.Version](
				version.Deneb,
			),
			Epoch: math.U64(constants.GenesisEpoch),
		},
		Eth1BlockHash:          primitives.ExecutionHash{},
		Eth1Timestamp:          0,
		Deposits:               make([]*primitives.Deposit, 0),
		ExecutionPayloadHeader: &engineprimitives.ExecutionHeaderDeneb{},
	}
}

// Genesis is the minimal eth1 genesis state for the beacon chain.
//
//nolint:lll // json tags.
type Genesis struct {
	// Fork is the fork version of the beacon chain.
	Fork *primitives.Fork `json:"fork"`
	// Eth1BlockHash is the hash of the Eth1 block.
	Eth1BlockHash primitives.ExecutionHash `json:"eth1BlockHash"`
	// Eth1Timestamp is the timestamp of the Eth1 block.
	Eth1Timestamp uint64 `json:"eth1Timestamp"`
	// Deposits is the list of genesis deposits.
	Deposits primitives.Deposits `json:"deposits"`
	// ExecutionPayloadHeader is the header of the genesis execution payload.
	ExecutionPayloadHeader engineprimitives.ExecutionPayloadHeader `json:"executionPayloadHeader"`
}

// InitializeBeaconStateFromEth1 initializes the beacon state from the Eth1
// chain.
func (sp *StateProcessor) InitializeBeaconStateFromEth1(
	emptySt state.BeaconState, genesis *Genesis,
) error {
	// Step 1: Setup the initial state
	if err := emptySt.SetGenesisTime(genesis.Eth1Timestamp); err != nil {
		return err
	}

	if err := emptySt.SetFork(genesis.Fork); err != nil {
		return err
	}

	if err := emptySt.SetEth1Data(&primitives.Eth1Data{
		BlockHash:    genesis.Eth1BlockHash,
		DepositCount: uint64(len(genesis.Deposits)),
	}); err != nil {
		return err
	}

	bodyRoot, err := (&types.BeaconBlockBodyDeneb{}).HashTreeRoot()
	if err != nil {
		return err
	}
	if err = emptySt.SetLatestBlockHeader(
		&primitives.BeaconBlockHeader{BodyRoot: bodyRoot},
	); err != nil {
		return err
	}

	for i := range sp.cs.EpochsPerHistoricalVector() {
		if err = emptySt.UpdateRandaoMixAtIndex(
			i, primitives.Bytes32(genesis.Eth1BlockHash),
		); err != nil {
			return err
		}
	}

	// Step 2: Process Deposits
	// leaves = list(map(lambda deposit: deposit.data, deposits))
	for _, deposit := range genesis.Deposits {
		// TODO: Merkle Root stuff
		// deposit_data_list = List[DepositData,
		// 2**DEPOSIT_CONTRACT_TREE_DEPTH](*leaves[:index + 1])
		// state.eth1_data.deposit_root = hash_tree_root(deposit_data_list)
		sp.processDeposit(emptySt, deposit)
	}

	// Step 3: Process Activations
	validators, err := emptySt.GetValidators()
	if err != nil {
		return err
	}

	var balance math.Gwei
	for index, validator := range validators {
		balance, err = emptySt.GetBalance(math.ValidatorIndex(index))
		if err != nil {
			return err
		}
		validator.EffectiveBalance =
			min(
				balance-balance%math.Gwei(sp.cs.EffectiveBalanceIncrement()),
				math.Gwei(sp.cs.MaxEffectiveBalance()),
			)
		if validator.EffectiveBalance == math.Gwei(
			sp.cs.MaxEffectiveBalance(),
		) {
			validator.ActivationEligibilityEpoch = math.U64(
				constants.GenesisEpoch,
			)
			validator.ActivationEpoch = math.U64(constants.GenesisEpoch)
		}
	}

	// Step 4: Set genesis validators root for domain separation and chain
	// versioning
	var genesisValidatorsRoot primitives.Root
	genesisValidatorsRoot, err = (&genutiltypes.ValidatorsMarshaling{
		Validators: validators,
	}).HashTreeRoot()
	if err != nil {
		return err
	}
	if err = emptySt.SetGenesisValidatorsRoot(genesisValidatorsRoot); err != nil {
		return err
	}

	// Step 5: Fill in sync committees
	// TODO: Figure out our own spec.
	// # Note: A duplicate committee is assigned for the current and next
	// committee at genesis
	// state.current_sync_committee = get_next_sync_committee(state)
	// state.next_sync_committee = get_next_sync_committee(state)

	// Step 6: Initialize the execution payload header
	emptySt.SetLatestExecutionPayloadHeader(genesis.ExecutionPayloadHeader)

	return nil
}
