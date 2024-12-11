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

package merkle_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle"
	mlib "github.com/berachain/beacon-kit/primitives/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz/schema"
	"github.com/stretchr/testify/require"
)

var (
	// beaconStateSchema is the schema for the BeaconState struct defined in
	// beacon-kit/mod/consensus-types/types/state.go.
	beaconStateSchema = schema.DefineContainer(
		schema.NewField("GenesisValidatorsRoot", schema.B32()),
		schema.NewField("Slot", schema.U64()),
		schema.NewField("Fork", schema.DefineContainer(
			schema.NewField("PreviousVersion", schema.B4()),
			schema.NewField("CurrentVersion", schema.B4()),
			schema.NewField("Epoch", schema.U64()),
		)),
		schema.NewField("LatestBlockHeader", schema.DefineContainer(
			schema.NewField("Slot", schema.U64()),
			schema.NewField("ProposerIndex", schema.U64()),
			schema.NewField("ParentBlockRoot", schema.B32()),
			schema.NewField("StateRoot", schema.B32()),
			schema.NewField("BodyRoot", schema.B32()),
		)),
		schema.NewField("BlockRoots", schema.DefineList(schema.B32(), 8192)),
		schema.NewField("StateRoots", schema.DefineList(schema.B32(), 8192)),
		schema.NewField("Eth1Data", schema.DefineContainer(
			schema.NewField("DepositRoot", schema.B32()),
			schema.NewField("DepositCount", schema.U64()),
			schema.NewField("BlockHash", schema.B32()),
		)),
		schema.NewField("Eth1DepositIndex", schema.U64()),
		schema.NewField("LatestExecutionPayloadHeader", schema.DefineContainer(
			schema.NewField("ParentHash", schema.B32()),
			schema.NewField("FeeRecipient", schema.B20()),
			schema.NewField("StateRoot", schema.B32()),
			schema.NewField("ReceiptsRoot", schema.B32()),
			schema.NewField("LogsBloom", schema.B256()),
			schema.NewField("Random", schema.U64()),
			schema.NewField("Number", schema.U64()),
			schema.NewField("GasLimit", schema.U64()),
			schema.NewField("GasUsed", schema.U64()),
			schema.NewField("Timestamp", schema.U64()),
			schema.NewField("ExtraData", schema.DefineByteList(32)),
			schema.NewField("BaseFeePerGas", schema.B32()),
			schema.NewField("BlockHash", schema.B32()),
			schema.NewField("TransactionsRoot", schema.B32()),
			schema.NewField("WithdrawalsRoot", schema.B32()),
			schema.NewField("BlobGasUsed", schema.U64()),
			schema.NewField("ExcessBlobGas", schema.U64()),
		)),
		schema.NewField("Validators", schema.DefineList(schema.DefineContainer(
			schema.NewField("Pubkey", schema.B48()),
			schema.NewField("WithdrawalCredentials", schema.B32()),
			schema.NewField("EffectiveBalance", schema.U64()),
			schema.NewField("Slashed", schema.Bool()),
			schema.NewField("ActivationEligibilityEpoch", schema.U64()),
			schema.NewField("ActivationEpoch", schema.U64()),
			schema.NewField("ExitEpoch", schema.U64()),
			schema.NewField("WithdrawableEpoch", schema.U64()),
		), types.MaxValidators)),
		schema.NewField(
			"Balances", schema.DefineList(schema.U64(), types.MaxValidators),
		),
		schema.NewField("RandaoMixes", schema.DefineList(schema.B32(), 65536)),
		schema.NewField("NextWithdrawalIndex", schema.U64()),
		schema.NewField("NextWithdrawalValidatorIndex", schema.U64()),
		schema.NewField(
			"Slashings", schema.DefineList(schema.U64(), types.MaxValidators),
		),
		schema.NewField("TotalSlashing", schema.U64()),
	)

	// beaconHeaderSchema is the schema for the BeaconBlockHeader struct defined
	// in beacon-kit/mod/consensus-types/types/header.go, with the SSZ
	// expansion of StateRoot to use the BeaconState.
	beaconHeaderSchema = schema.DefineContainer(
		schema.NewField("Slot", schema.U64()),
		schema.NewField("ProposerIndex", schema.U64()),
		schema.NewField("ParentRoot", schema.B32()),
		schema.NewField("State", beaconStateSchema),
		schema.NewField("BodyRoot", schema.B32()),
	)
)

// TestGIndexProposerIndexDeneb tests the generalized index of the proposer
// index in the beacon block on the Deneb fork.
func TestGIndexProposerIndexDeneb(t *testing.T) {
	// GIndex of the proposer index in the beacon block.
	_, proposerIndexGIndexDenebBlock, _, err := mlib.ObjectPath[
		mlib.GeneralizedIndex, [32]byte,
	]("ProposerIndex").GetGeneralizedIndex(beaconHeaderSchema)
	require.NoError(t, err)
	require.Equal(
		t,
		merkle.ProposerIndexGIndexDenebBlock,
		int(proposerIndexGIndexDenebBlock),
	)
}

// TestGIndicesValidatorPubkeyDeneb tests the generalized indices used by
// beacon state proofs for validator pubkeys on the Deneb fork.
func TestGIndicesValidatorPubkeyDeneb(t *testing.T) {
	// GIndex of state in the block.
	_, stateGIndexDenebBlock, _, err := mlib.ObjectPath[
		mlib.GeneralizedIndex, [32]byte,
	]("State").GetGeneralizedIndex(beaconHeaderSchema)
	require.NoError(t, err)
	require.Equal(t, merkle.StateGIndexDenebBlock, int(stateGIndexDenebBlock))

	// GIndex of the 0 validator's pubkey in the state.
	_, zeroValidatorPubkeyGIndexDenebState, _, err := mlib.ObjectPath[
		mlib.GeneralizedIndex, [32]byte,
	]("Validators/0/Pubkey").GetGeneralizedIndex(beaconStateSchema)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorPubkeyGIndexDenebState,
		int(zeroValidatorPubkeyGIndexDenebState),
	)

	// GIndex of the 0 validator's pubkey in the block.
	_, zeroValidatorPubkeyGIndexDenebBlock, _, err := mlib.ObjectPath[
		mlib.GeneralizedIndex, [32]byte,
	]("State/Validators/0/Pubkey").GetGeneralizedIndex(beaconHeaderSchema)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorPubkeyGIndexDenebBlock,
		int(zeroValidatorPubkeyGIndexDenebBlock),
	)

	// Concatenation is consistent.
	concatValidatorPubkeyStateToBlock := mlib.GeneralizedIndices{
		stateGIndexDenebBlock,
		zeroValidatorPubkeyGIndexDenebState,
	}.Concat()
	require.Equal(t,
		zeroValidatorPubkeyGIndexDenebBlock,
		concatValidatorPubkeyStateToBlock,
	)

	// GIndex offset of the next validator's pubkey.
	_, oneValidatorPubkeyGIndexDenebState, _, err := mlib.ObjectPath[
		mlib.GeneralizedIndex, [32]byte,
	]("Validators/1/Pubkey").GetGeneralizedIndex(beaconStateSchema)
	require.NoError(t, err)
	require.Equal(t,
		mlib.GeneralizedIndex(merkle.ValidatorPubkeyGIndexOffset),
		oneValidatorPubkeyGIndexDenebState-zeroValidatorPubkeyGIndexDenebState,
	)
}

// TestGInidicesExecutionDeneb tests the generalized indices used by
// beacon state proofs from the execution payload header on the Deneb fork.
func TestGInidicesExecutionDeneb(t *testing.T) {
	// GIndex of the execution number in the state.
	_, executionNumberGIndexDenebState, _, err := mlib.ObjectPath[
		mlib.GeneralizedIndex, [32]byte,
	]("LatestExecutionPayloadHeader/Number").GetGeneralizedIndex(
		beaconStateSchema,
	)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ExecutionNumberGIndexDenebState,
		int(executionNumberGIndexDenebState),
	)

	// GIndex of the execution number in the block.
	_, executionNumberGIndexDenebBlock, _, err := mlib.ObjectPath[
		mlib.GeneralizedIndex, [32]byte,
	]("State/LatestExecutionPayloadHeader/Number").GetGeneralizedIndex(
		beaconHeaderSchema,
	)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ExecutionNumberGIndexDenebBlock,
		int(executionNumberGIndexDenebBlock),
	)

	// Concatenation is consistent.
	concatExecutionNumberStateToBlock := mlib.GeneralizedIndices{
		merkle.StateGIndexDenebBlock,
		executionNumberGIndexDenebState,
	}.Concat()
	require.Equal(t,
		executionNumberGIndexDenebBlock,
		concatExecutionNumberStateToBlock,
	)

	// GIndex of the execution fee recipient in the state.
	_, executionFeeRecipientGIndexDenebState, _, err := mlib.ObjectPath[
		mlib.GeneralizedIndex, [32]byte,
	]("LatestExecutionPayloadHeader/FeeRecipient").GetGeneralizedIndex(
		beaconStateSchema,
	)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ExecutionFeeRecipientGIndexDenebState,
		int(executionFeeRecipientGIndexDenebState),
	)

	// GIndex of the execution fee recipient in the block.
	_, executionFeeRecipientGIndexDenebBlock, _, err := mlib.ObjectPath[
		mlib.GeneralizedIndex, [32]byte,
	]("State/LatestExecutionPayloadHeader/FeeRecipient").GetGeneralizedIndex(
		beaconHeaderSchema,
	)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ExecutionFeeRecipientGIndexDenebBlock,
		int(executionFeeRecipientGIndexDenebBlock),
	)

	// Concatenation is consistent.
	concatExecutionFeeRecipientStateToBlock := mlib.GeneralizedIndices{
		merkle.StateGIndexDenebBlock,
		executionFeeRecipientGIndexDenebState,
	}.Concat()
	require.Equal(t,
		executionFeeRecipientGIndexDenebBlock,
		concatExecutionFeeRecipientStateToBlock,
	)
}
