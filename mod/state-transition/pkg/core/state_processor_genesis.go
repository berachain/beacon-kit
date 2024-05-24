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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// InitializePreminedBeaconStateFromEth1 initializes the beacon state from the
// eth1 block hash.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) InitializePreminedBeaconStateFromEth1(
	st BeaconStateT,
	eth1BlockHash common.ExecutionHash,
	eth1Timestamp math.U64,
	deposits []*types.Deposit,
	executionPayloadHeader engineprimitives.ExecutionPayloadHeader,
	genesisVersion primitives.Version,
) ([]*transition.ValidatorUpdate, error) {
	fork := &types.Fork{
		PreviousVersion: genesisVersion,
		CurrentVersion:  genesisVersion,
		Epoch:           0,
	}

	if err := st.SetSlot(0); err != nil {
		return nil, err
	}

	if err := st.SetFork(fork); err != nil {
		return nil, err
	}

	if err := st.SetEth1DepositIndex(0); err != nil {
		return nil, err
	}

	if err := st.SetEth1Data(&types.Eth1Data{
		DepositRoot:  bytes.B32(common.ZeroHash),
		DepositCount: 0,
		BlockHash:    eth1BlockHash,
	}); err != nil {
		return nil, err
	}

	bodyRoot, err := (&types.BeaconBlockBodyDeneb{
		BeaconBlockBodyBase: types.BeaconBlockBodyBase{},
		ExecutionPayload: &types.ExecutableDataDeneb{
			LogsBloom: make([]byte, 256),
			ExtraData: make([]byte, 32),
		},
	}).HashTreeRoot()
	if err != nil {
		return nil, err
	}

	if err := st.SetLatestBlockHeader(&types.BeaconBlockHeader{
		BodyRoot: bodyRoot,
	}); err != nil {
		return nil, err
	}

	for i := range sp.cs.EpochsPerHistoricalVector() {
		if err := st.UpdateRandaoMixAtIndex(uint64(i), bytes.B32(eth1BlockHash)); err != nil {
			return nil, err
		}
	}

	// Prime the db so that processDeposit doesn't fail.
	if err := st.SetGenesisValidatorsRoot(primitives.Root{}); err != nil {
		return nil, err
	}

	for _, deposit := range deposits {
		// TODO: process deposits into eth1 data.
		if err := sp.processDeposit(st, deposit); err != nil {
			return nil, err
		}
	}

	// TODO: process activations.
	validators, err := st.GetValidators()
	if err != nil {
		return nil, err
	}

	validatorsRoot, err := ssz.MerkleizeListComposite[
		common.ChainSpec, math.U64, [32]byte,
	](validators, uint64(len(validators)))
	if err != nil {
		return nil, err
	}

	if err = st.SetGenesisValidatorsRoot(validatorsRoot); err != nil {
		return nil, err
	}

	if err = st.SetLatestExecutionPayloadHeader(executionPayloadHeader); err != nil {
		return nil, err
	}
	// Setup a bunch of 0s to prime the DB.
	for i := range sp.cs.HistoricalRootsLimit() {
		//#nosec:G701 // won't overflow in practice.
		if err := st.UpdateBlockRootAtIndex(i, primitives.Root{}); err != nil {
			return nil, err
		}
		if err := st.UpdateStateRootAtIndex(i, primitives.Root{}); err != nil {
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

	var updates []*transition.ValidatorUpdate
	updates, err = sp.processSyncCommitteeUpdates(st)
	if err != nil {
		return nil, err
	}
	st.Save()
	return updates, nil
}
