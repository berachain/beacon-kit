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

package merkle_test

import (
	"fmt"
	"testing"

	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz/schema"
	mlib "github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/stretchr/testify/require"
)

var (
	beaconStateFieldsDeneb = []*schema.Field{
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
			schema.NewField("Random", schema.B32()),
			schema.NewField("Number", schema.U64()),
			schema.NewField("GasLimit", schema.U64()),
			schema.NewField("GasUsed", schema.U64()),
			schema.NewField("Timestamp", schema.U64()),
			schema.NewField("ExtraData", schema.DefineByteList(32)),
			schema.NewField("BaseFeePerGas", schema.U256()),
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
		), constants.ValidatorsRegistryLimit)),
		schema.NewField(
			"Balances", schema.DefineList(schema.U64(), constants.ValidatorsRegistryLimit),
		),
		schema.NewField("RandaoMixes", schema.DefineList(schema.B32(), 65536)),
		schema.NewField("NextWithdrawalIndex", schema.U64()),
		schema.NewField("NextWithdrawalValidatorIndex", schema.U64()),
		schema.NewField(
			"Slashings", schema.DefineList(schema.U64(), constants.ValidatorsRegistryLimit),
		),
		schema.NewField("TotalSlashing", schema.U64()),
	}

	additionalBeaconStateFieldsElectra = []*schema.Field{
		schema.NewField("PendingPartialWithdrawals", schema.DefineList(schema.DefineContainer(
			schema.NewField("ValidatorIndex", schema.U64()),
			schema.NewField("Amount", schema.U64()),
			schema.NewField("WithdrawableEpoch", schema.U64()),
		), constants.PendingPartialWithdrawalsLimit)),
	}

	// beaconStateSchemaDeneb is the schema for the BeaconState struct in the Deneb forks.
	beaconStateSchemaDeneb = schema.DefineContainer(beaconStateFieldsDeneb...)

	// beaconStateSchemaElectra is the schema for the BeaconState struct in the Electra forks.
	beaconStateSchemaElectra = schema.DefineContainer(
		append(beaconStateFieldsDeneb, additionalBeaconStateFieldsElectra...)...,
	)
)

var (
	// beaconHeaderSchemaDeneb is the schema for the BeaconBlockHeader in the Deneb forks, with the
	// SSZ expansion of StateRoot to use the BeaconState.
	beaconHeaderSchemaDeneb = schema.DefineContainer(
		schema.NewField("Slot", schema.U64()),
		schema.NewField("ProposerIndex", schema.U64()),
		schema.NewField("ParentRoot", schema.B32()),
		schema.NewField("State", beaconStateSchemaDeneb),
		schema.NewField("BodyRoot", schema.B32()),
	)

	// beaconHeaderSchemaElectra is the schema for the BeaconBlockHeader in the Electra forks, with
	// the SSZ expansion of StateRoot to use the BeaconState.
	beaconHeaderSchemaElectra = schema.DefineContainer(
		schema.NewField("Slot", schema.U64()),
		schema.NewField("ProposerIndex", schema.U64()),
		schema.NewField("ParentRoot", schema.B32()),
		schema.NewField("State", beaconStateSchemaElectra),
		schema.NewField("BodyRoot", schema.B32()),
	)
)

// TestGIndexProposerIndex tests the generalized index of the proposer
// index in the beacon block.
func TestGIndexProposerIndex(t *testing.T) {
	t.Parallel()

	// Deneb forks.
	_, proposerIndexGIndexBlock, _, err := mlib.ObjectPath(
		"ProposerIndex",
	).GetGeneralizedIndex(beaconHeaderSchemaDeneb)
	require.NoError(t, err)
	require.Equal(
		t,
		merkle.ProposerIndexGIndexBlock,
		int(proposerIndexGIndexBlock),
	)

	// Electra forks.
	_, proposerIndexGIndexBlockElectra, _, err := mlib.ObjectPath(
		"ProposerIndex",
	).GetGeneralizedIndex(beaconHeaderSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ProposerIndexGIndexBlock,
		int(proposerIndexGIndexBlockElectra),
	)
}

// TestGIndicesValidatorPubkeyDeneb tests the generalized indices used by
// beacon state proofs for validator pubkeys on the Deneb forks.
func TestGIndicesValidatorPubkeyDeneb(t *testing.T) {
	t.Parallel()

	// GIndex of state in the block.
	_, stateGIndexBlock, _, err := mlib.ObjectPath(
		"State",
	).GetGeneralizedIndex(beaconHeaderSchemaDeneb)
	require.NoError(t, err)
	require.Equal(t, merkle.StateGIndexBlock, int(stateGIndexBlock))

	// GIndex of the 0 validator's pubkey in the state.
	_, zeroValidatorPubkeyGIndexState, _, err := mlib.ObjectPath(
		"Validators/0/Pubkey",
	).GetGeneralizedIndex(beaconStateSchemaDeneb)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorPubkeyGIndexDenebState,
		int(zeroValidatorPubkeyGIndexState),
	)

	// GIndex of the 0 validator's pubkey in the block.
	_, zeroValidatorPubkeyGIndexBlock, _, err := mlib.ObjectPath(
		"State/Validators/0/Pubkey",
	).GetGeneralizedIndex(beaconHeaderSchemaDeneb)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorPubkeyGIndexDenebBlock,
		int(zeroValidatorPubkeyGIndexBlock),
	)

	// Concatenation is consistent.
	concatValidatorPubkeyStateToBlock := mlib.GeneralizedIndices{
		mlib.GeneralizedIndex(stateGIndexBlock),
		mlib.GeneralizedIndex(zeroValidatorPubkeyGIndexState),
	}.Concat()
	require.Equal(t,
		zeroValidatorPubkeyGIndexBlock,
		uint64(concatValidatorPubkeyStateToBlock),
	)

	// GIndex offset of the next validator's pubkey.
	_, oneValidatorPubkeyGIndexState, _, err := mlib.ObjectPath(
		"Validators/1/Pubkey",
	).GetGeneralizedIndex(beaconStateSchemaDeneb)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ValidatorGIndexOffset,
		int(oneValidatorPubkeyGIndexState-zeroValidatorPubkeyGIndexState),
	)
}

// TestGIndicesValidatorPubkeyElectra tests the generalized indices used by
// beacon state proofs for validator pubkeys on the Electra forks.
func TestGIndicesValidatorPubkeyElectra(t *testing.T) {
	t.Parallel()

	// GIndex of state in the block.
	_, stateGIndexBlock, _, err := mlib.ObjectPath(
		"State",
	).GetGeneralizedIndex(beaconHeaderSchemaElectra)
	require.NoError(t, err)
	require.Equal(t, merkle.StateGIndexBlock, int(stateGIndexBlock))

	// GIndex of the 0 validator's pubkey in the state.
	_, zeroValidatorPubkeyGIndexState, _, err := mlib.ObjectPath(
		"Validators/0/Pubkey",
	).GetGeneralizedIndex(beaconStateSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorPubkeyGIndexElectraState,
		int(zeroValidatorPubkeyGIndexState),
	)

	// GIndex of the 0 validator's pubkey in the block.
	_, zeroValidatorPubkeyGIndexBlock, _, err := mlib.ObjectPath(
		"State/Validators/0/Pubkey",
	).GetGeneralizedIndex(beaconHeaderSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorPubkeyGIndexElectraBlock,
		int(zeroValidatorPubkeyGIndexBlock),
	)

	// Concatenation is consistent.
	concatValidatorPubkeyStateToBlock := mlib.GeneralizedIndices{
		mlib.GeneralizedIndex(stateGIndexBlock),
		mlib.GeneralizedIndex(zeroValidatorPubkeyGIndexState),
	}.Concat()
	require.Equal(t,
		zeroValidatorPubkeyGIndexBlock,
		uint64(concatValidatorPubkeyStateToBlock),
	)

	// GIndex offset of the next validator's pubkey.
	_, oneValidatorPubkeyGIndexState, _, err := mlib.ObjectPath(
		"Validators/1/Pubkey",
	).GetGeneralizedIndex(beaconStateSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ValidatorGIndexOffset,
		int(oneValidatorPubkeyGIndexState-zeroValidatorPubkeyGIndexState),
	)
}

func TestValidatorWithdrawalCredentialsGIndexElectra(t *testing.T) {
	t.Parallel()

	// GIndex of the 0 validator's withdrawal credentials in the state.
	_, zeroValidatorWithdrawalCredentialsGIndexState, _, err := mlib.ObjectPath(
		"Validators/0/WithdrawalCredentials",
	).GetGeneralizedIndex(beaconStateSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorCredentialsGIndexElectraState,
		int(zeroValidatorWithdrawalCredentialsGIndexState),
	)

	// GIndex of the 0 validator's withdrawal credentials in the block.
	_, zeroValidatorWithdrawalCredentialsGIndexBlock, _, err := mlib.ObjectPath(
		"State/Validators/0/WithdrawalCredentials",
	).GetGeneralizedIndex(beaconHeaderSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorCredentialsGIndexElectraBlock,
		int(zeroValidatorWithdrawalCredentialsGIndexBlock),
	)

	// Concatenation is consistent.
	concatValidatorWithdrawalCredentialsStateToBlock := mlib.GeneralizedIndices{
		mlib.GeneralizedIndex(merkle.StateGIndexBlock),
		mlib.GeneralizedIndex(zeroValidatorWithdrawalCredentialsGIndexState),
	}.Concat()
	require.Equal(t,
		zeroValidatorWithdrawalCredentialsGIndexBlock,
		uint64(concatValidatorWithdrawalCredentialsStateToBlock),
	)

	// GIndex offset of the next validator's withdrawal credentials.
	_, oneValidatorWithdrawalCredentialsGIndexState, _, err := mlib.ObjectPath(
		"Validators/1/WithdrawalCredentials",
	).GetGeneralizedIndex(beaconStateSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ValidatorGIndexOffset,
		int(oneValidatorWithdrawalCredentialsGIndexState-zeroValidatorWithdrawalCredentialsGIndexState),
	)
}

func TestValidatorBalanceGIndexElectra(t *testing.T) {
	t.Parallel()

	// GIndex of the 0 validator's balance in the state.
	_, zeroValidatorBalanceGIndexState, _, err := mlib.ObjectPath(
		"Balances/0",
	).GetGeneralizedIndex(beaconStateSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorBalanceGIndexElectraState,
		int(zeroValidatorBalanceGIndexState),
	)

	// GIndex of the 0 validator's balance in the block.
	_, zeroValidatorBalanceGIndexBlock, _, err := mlib.ObjectPath(
		"State/Balances/0",
	).GetGeneralizedIndex(beaconHeaderSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.ZeroValidatorBalanceGIndexElectraBlock,
		int(zeroValidatorBalanceGIndexBlock),
	)

	// Concatenation is consistent.
	concatValidatorBalanceStateToBlock := mlib.GeneralizedIndices{
		mlib.GeneralizedIndex(merkle.StateGIndexBlock),
		mlib.GeneralizedIndex(zeroValidatorBalanceGIndexState),
	}.Concat()
	require.Equal(t,
		zeroValidatorBalanceGIndexBlock,
		uint64(concatValidatorBalanceStateToBlock),
	)

	// Verify that balances 0-3 share the same GIndex (packed in same leaf)
	for i := 0; i < 4; i++ {
		path := fmt.Sprintf("Balances/%d", i)
		_, balanceGIndex, _, err := mlib.ObjectPath(path).GetGeneralizedIndex(beaconStateSchemaElectra)
		require.NoError(t, err)
		require.Equal(t, zeroValidatorBalanceGIndexState, balanceGIndex)
	}
	
	// GIndex offset of the next validator's balance.
	// Balances are packed 4 per leaf (uint64 values, 32 bytes per leaf)
	_, balanceGIndex4, _, err := mlib.ObjectPath("Balances/4").GetGeneralizedIndex(beaconStateSchemaElectra)
	require.NoError(t, err)
	require.Equal(t,
		merkle.BalanceGIndexOffset,
		int(balanceGIndex4-zeroValidatorBalanceGIndexState),
	)
}
