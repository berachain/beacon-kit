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

package backend

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	sszTypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-api/backend/mocks"
	response "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/stretchr/testify/mock"
)

func NewMockBackend() *Backend {
	sdb := &mocks.StateDB{}
	bdb := &mocks.BlockDB{}
	ns := &mocks.NodeState{}
	setStateDBMockValues(sdb)
	setBlockDBMockValues(bdb)
	setNodeStateMockValues(ns)

	b := New(
		func(context.Context, string) StateDB {
			return sdb
		},
		func(context.Context, string) BlockDB {
			return bdb
		},
		func(context.Context) NodeState {
			return ns
		},
	)
	return b
}

func setNodeStateMockValues(ns *mocks.NodeState) {
	ns.EXPECT().GetSpecParams().Return(&response.SpecParamsResponse{}, nil)
	ns.EXPECT().GetBlsToExecutionChanges().Return([]*response.MessageSignature{
		{
			Message: response.BtsToExecutionChangeData{
				ValidatorIndex:     0,
				FromBlsPubkey:      crypto.BLSPubkey{0x01},
				ToExecutionAddress: common.ExecutionAddress{0x01},
			},
			Signature: crypto.BLSSignature{0x01},
		},
		{
			Message: response.BtsToExecutionChangeData{
				ValidatorIndex:     1,
				FromBlsPubkey:      crypto.BLSPubkey{0x01},
				ToExecutionAddress: common.ExecutionAddress{0x01},
			},
			Signature: crypto.BLSSignature{0x01},
		},
	}, nil)
	ns.EXPECT().GetVoluntaryExits().Return(
		[]*response.MessageSignature{
			{
				Message: response.VoluntaryExitData{
					ValidatorIndex: 0,
					Epoch:          0,
				},
				Signature: crypto.BLSSignature{0x01},
			},
			{
				Message: response.VoluntaryExitData{
					ValidatorIndex: 1,
					Epoch:          1,
				},
				Signature: crypto.BLSSignature{0x01},
			},
		},
		nil,
	)
}

func setBlockDBMockValues(bdb *mocks.BlockDB) {
	bdb.EXPECT().
		GetBlockBlobSidecars(mock.Anything).
		Return([]*sszTypes.BlobSidecar{
			sszTypes.BuildBlobSidecar(0,
				types.NewBeaconBlockHeader(0, 0,
					primitives.Root{0x01}, primitives.Root{0x01}, primitives.Root{0x01}),
				&eip4844.Blob{},
				eip4844.KZGCommitment{},
				eip4844.KZGProof{},
				nil),
		}, nil)
	block := types.BeaconBlock{}
	bdb.EXPECT().GetBlock().Return(block.Empty(version.Deneb), nil)
	blockHeader := &response.BlockHeaderData{
		Root:      primitives.Root{0x01},
		Canonical: true,
		Header: response.MessageResponse{
			Message: types.NewBeaconBlockHeader(
				0,
				0,
				primitives.Root{0x01},
				primitives.Root{0x01},
				primitives.Root{0x01},
			),
		},
		Signature: crypto.BLSSignature{0x01},
	}
	bdb.EXPECT().GetBlockHeader().Return(blockHeader, nil)
	bdb.EXPECT().
		GetBlockHeaders(mock.Anything, mock.Anything).
		Return([]*response.BlockHeaderData{blockHeader, blockHeader}, nil)
	bdb.EXPECT().
		GetBlockPropserDuties(mock.Anything).
		Return([]*response.ProposerDutiesData{
			{
				Pubkey:         crypto.BLSPubkey{0x01},
				ValidatorIndex: 0,
				Slot:           0,
			},
			{
				Pubkey:         crypto.BLSPubkey{0x01},
				ValidatorIndex: 1,
				Slot:           1,
			},
		}, nil)
	bdb.EXPECT().GetBlockRewards().Return(&response.BlockRewardsData{
		ProposerIndex:     0,
		Total:             0,
		Attestations:      0,
		SyncAggregate:     0,
		ProposerSlashings: 0,
		AttesterSlashings: 0,
	}, nil)
}

func setStateDBMockValues(sdb *mocks.StateDB) {
	sdb.EXPECT().GetGenesisDetails().Return(&response.GenesisData{
		GenesisTime:           0,
		GenesisValidatorsRoot: primitives.Root{0x01},
		GenesisForkVersion:    primitives.Version{0x01},
	}, nil)
	sdb.EXPECT().GetSlot().Return(1, nil)
	sdb.EXPECT().GetLatestExecutionPayloadHeader().Return(nil, nil)
	sdb.EXPECT().SetLatestExecutionPayloadHeader(mock.Anything).Return(nil)
	sdb.EXPECT().GetEth1DepositIndex().Return(0, nil)
	sdb.EXPECT().SetEth1DepositIndex(mock.Anything).Return(nil)
	sdb.EXPECT().GetBalance(mock.Anything).Return(1, nil)
	sdb.EXPECT().SetBalance(mock.Anything, mock.Anything).Return(nil)
	sdb.EXPECT().SetSlot(mock.Anything).Return(nil)
	sdb.EXPECT().GetFork().Return(nil, nil)
	sdb.EXPECT().SetFork(mock.Anything).Return(nil)
	sdb.EXPECT().GetLatestBlockHeader().Return(nil, nil)
	sdb.EXPECT().SetLatestBlockHeader(mock.Anything).Return(nil)
	sdb.EXPECT().
		GetBlockRootAtIndex(mock.Anything).
		Return(primitives.Root{0x01}, nil)
	sdb.EXPECT().
		StateRootAtIndex(mock.Anything).
		Return(primitives.Root{0x01}, nil)
	sdb.EXPECT().GetEth1Data().Return(nil, nil)
	sdb.EXPECT().SetEth1Data(mock.Anything).Return(nil)
	sdb.EXPECT().GetValidators().Return(nil, nil)
	sdb.EXPECT().GetBalances().Return(nil, nil)
	sdb.EXPECT().GetNextWithdrawalIndex().Return(0, nil)
	sdb.EXPECT().SetNextWithdrawalIndex(mock.Anything).Return(nil)
	sdb.EXPECT().GetNextWithdrawalValidatorIndex().Return(0, nil)
	sdb.EXPECT().SetNextWithdrawalValidatorIndex(mock.Anything).Return(nil)
	sdb.EXPECT().GetTotalSlashing().Return(0, nil)
	sdb.EXPECT().SetTotalSlashing(mock.Anything).Return(nil)
	sdb.EXPECT().
		GetRandaoMixAtIndex(mock.Anything).
		Return(primitives.Bytes32{0x01}, nil)
	sdb.EXPECT().GetSlashings().Return(nil, nil)
	sdb.EXPECT().SetSlashingAtIndex(mock.Anything, mock.Anything).Return(nil)
	sdb.EXPECT().GetSlashingAtIndex(mock.Anything).Return(0, nil)
	sdb.EXPECT().GetTotalValidators().Return(0, nil)
	sdb.EXPECT().GetTotalActiveBalances(mock.Anything).Return(0, nil)
	sdb.EXPECT().ValidatorByIndex(mock.Anything).Return(&types.Validator{
		Pubkey:                     crypto.BLSPubkey{0x01},
		WithdrawalCredentials:      types.WithdrawalCredentials{0x01},
		EffectiveBalance:           0,
		Slashed:                    false,
		ActivationEligibilityEpoch: 0,
		ActivationEpoch:            0,
		ExitEpoch:                  0,
		WithdrawableEpoch:          0,
	}, nil)
	sdb.EXPECT().
		UpdateBlockRootAtIndex(mock.Anything, mock.Anything).
		Return(nil)
	sdb.EXPECT().
		UpdateStateRootAtIndex(mock.Anything, mock.Anything).
		Return(nil)
	sdb.EXPECT().
		UpdateRandaoMixAtIndex(mock.Anything, mock.Anything).
		Return(nil)
	sdb.EXPECT().
		UpdateValidatorAtIndex(mock.Anything, mock.Anything).
		Return(nil)
	sdb.EXPECT().ValidatorIndexByPubkey(mock.Anything).Return(0, nil)
	sdb.EXPECT().AddValidator(mock.Anything).Return(nil)
	sdb.EXPECT().GetValidatorsByEffectiveBalance().Return(nil, nil)
	sdb.EXPECT().
		GetStateCommittees(mock.Anything).
		Return([]*response.CommitteeData{
			{Index: 0, Slot: 0, Validators: []uint64{1, 2}},
			{Index: 1, Slot: 1, Validators: []uint64{1, 2}},
		}, nil)
	sdb.EXPECT().
		GetStateSyncCommittees(mock.Anything).
		Return(&response.SyncCommitteeData{
			Validators:          []uint64{1, 2},
			ValidatorAggregates: [][]uint64{{1, 2}, {1, 2}},
		}, nil)
	sdb.EXPECT().GetEpoch().Return(0, nil)
}
