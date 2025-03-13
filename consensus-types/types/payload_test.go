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

package types_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	karalabessz "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

func generateExecutionPayload() *types.ExecutionPayload {
	var (
		transactions = [][]byte{{0x07}}
		withdrawals  = []*engineprimitives.Withdrawal{
			{
				Index:     0,
				Validator: 0,
				Address:   common.ExecutionAddress{},
				Amount:    0,
			},
		}
	)

	ep := &types.ExecutionPayload{
		Versionable:   (&types.BeaconBlock{}).WithForkVersion(version.Deneb1()),
		ParentHash:    common.ExecutionHash{},
		FeeRecipient:  common.ExecutionAddress{},
		StateRoot:     bytes.B32{},
		ReceiptsRoot:  bytes.B32{},
		LogsBloom:     bytes.B256{},
		Random:        bytes.B32{},
		Number:        math.U64(0),
		GasLimit:      math.U64(0),
		GasUsed:       math.U64(0),
		Timestamp:     math.U64(0),
		ExtraData:     []byte{0x01},
		BaseFeePerGas: &math.U256{},
		BlockHash:     common.ExecutionHash{},
		Transactions:  transactions,
		Withdrawals:   withdrawals,
		BlobGasUsed:   math.U64(0),
		ExcessBlobGas: math.U64(0),
	}
	return ep
}

func TestExecutionPayload_Serialization(t *testing.T) {
	t.Parallel()
	original := generateExecutionPayload()

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	unmarshalled, err := (&types.ExecutionPayload{}).NewFromSSZ(data, original.GetForkVersion())
	require.NoError(t, err)
	require.Equal(t, original, unmarshalled)

	var buf []byte
	buf, err = original.MarshalSSZTo(buf)
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}

func TestExecutionPayload_SizeSSZ(t *testing.T) {
	t.Parallel()
	payload := generateExecutionPayload()
	size := karalabessz.Size(payload)
	require.Equal(t, uint32(578), size)

	_, err := (&types.ExecutionPayload{}).NewFromSSZ(
		[]byte{0x01, 0x02, 0x03}, // Invalid data
		version.Deneb1(),
	)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestExecutionPayload_HashTreeRoot(t *testing.T) {
	t.Parallel()
	payload := generateExecutionPayload()
	require.NotPanics(t, func() {
		_ = payload.HashTreeRoot()
	})
}

func TestExecutionPayload_GetTree(t *testing.T) {
	t.Parallel()
	payload := generateExecutionPayload()
	tree, err := payload.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestExecutionPayload_Getters(t *testing.T) {
	t.Parallel()
	payload := generateExecutionPayload()
	require.Equal(t, common.ExecutionHash{}, payload.GetParentHash())
	require.Equal(
		t,
		common.ExecutionAddress{},
		payload.GetFeeRecipient(),
	)

	transactions := make(engineprimitives.Transactions, 1)
	transactions[0] = []byte{0x07}
	withdrawals := make(engineprimitives.Withdrawals, 1)
	withdrawals[0] = &engineprimitives.Withdrawal{
		Index:     0,
		Validator: 0,
		Address:   common.ExecutionAddress{},
		Amount:    0,
	}
	require.Equal(t, common.ExecutionHash{}, payload.GetParentHash())
	require.Equal(t, common.ExecutionAddress{}, payload.GetFeeRecipient())
	require.Equal(t, bytes.B32{}, payload.GetStateRoot())
	require.Equal(t, bytes.B32{}, payload.GetReceiptsRoot())
	require.Equal(t, bytes.B256{}, payload.GetLogsBloom())
	require.Equal(t, bytes.B32{}, payload.GetPrevRandao())
	require.Equal(t, math.U64(0), payload.GetNumber())
	require.Equal(t, math.U64(0), payload.GetGasLimit())
	require.Equal(t, math.U64(0), payload.GetGasUsed())
	require.Equal(t, math.U64(0), payload.GetTimestamp())
	require.Equal(t, []byte{0x01}, payload.GetExtraData())
	require.Equal(t, &math.U256{}, payload.GetBaseFeePerGas())
	require.Equal(t, common.ExecutionHash{}, payload.GetBlockHash())
	require.Equal(t, transactions, payload.GetTransactions())
	require.Equal(t, withdrawals, payload.GetWithdrawals())
	require.Equal(t, math.U64(0), payload.GetBlobGasUsed())
	require.Equal(t, math.U64(0), payload.GetExcessBlobGas())
}

func TestExecutionPayload_MarshalJSON(t *testing.T) {
	t.Parallel()
	payload := generateExecutionPayload()

	data, err := payload.MarshalJSON()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.ExecutionPayload
	err = unmarshalled.UnmarshalJSON(data)
	require.NoError(t, err)

	unmarshalled.Versionable = payload.Versionable
	require.Equal(t, payload, &unmarshalled)
}

func TestExecutionPayload_MarshalJSON_ValueAndPointer(t *testing.T) {
	t.Parallel()
	val := types.ExecutionPayload{}

	// Marshal on raw val uses default json marshal
	valSerialized, err := json.Marshal(val)
	require.NoError(t, err)

	// Marshal on ptr val uses implemented MarshalJSON
	ptrSerialized, err := json.Marshal(&val)
	require.NoError(t, err)

	require.Equal(t, valSerialized, ptrSerialized)
}

func TestExecutionPayload_IsNil(t *testing.T) {
	t.Parallel()
	var payload *types.ExecutionPayload
	require.True(t, payload.IsNil())

	payload = generateExecutionPayload()
	require.False(t, payload.IsNil())
}

func TestExecutionPayload_IsBlinded(t *testing.T) {
	t.Parallel()
	payload := generateExecutionPayload()
	require.False(t, payload.IsBlinded())
}

func TestExecutionPayload_Version(t *testing.T) {
	t.Parallel()
	payload := generateExecutionPayload()
	require.Equal(t, version.Deneb1(), payload.GetForkVersion())
}

func TestExecutionPayload_ToHeader(t *testing.T) {
	t.Parallel()
	payload := &types.ExecutionPayload{
		Versionable:   (&types.BeaconBlock{}).WithForkVersion(version.Deneb1()),
		ParentHash:    common.ExecutionHash{},
		FeeRecipient:  common.ExecutionAddress{},
		StateRoot:     bytes.B32{},
		ReceiptsRoot:  bytes.B32{},
		LogsBloom:     bytes.B256{},
		Random:        bytes.B32{},
		Number:        math.U64(0),
		GasLimit:      math.U64(0),
		GasUsed:       math.U64(0),
		Timestamp:     math.U64(0),
		ExtraData:     []byte{},
		BaseFeePerGas: &math.U256{},
		BlockHash:     common.ExecutionHash{},
		Transactions:  [][]byte{{0x01}},
		Withdrawals:   engineprimitives.Withdrawals{},
		BlobGasUsed:   math.U64(0),
		ExcessBlobGas: math.U64(0),
	}
	header, err := payload.ToHeader()
	require.NoError(t, err)
	require.NotNil(t, header)

	require.Equal(t, payload.GetParentHash(), header.GetParentHash())
	require.Equal(t, payload.GetFeeRecipient(), header.GetFeeRecipient())
	require.Equal(t, payload.GetStateRoot(), header.GetStateRoot())
	require.Equal(t, payload.GetReceiptsRoot(), header.GetReceiptsRoot())
	require.Equal(t, payload.GetLogsBloom(), header.GetLogsBloom())
	require.Equal(t, payload.GetPrevRandao(), header.GetPrevRandao())
	require.Equal(t, payload.GetNumber(), header.GetNumber())
	require.Equal(t, payload.GetGasLimit(), header.GetGasLimit())
	require.Equal(t, payload.GetGasUsed(), header.GetGasUsed())
	require.Equal(t, payload.GetTimestamp(), header.GetTimestamp())
	require.Equal(t, payload.GetExtraData(), header.GetExtraData())
	require.Equal(t, payload.GetBaseFeePerGas(), header.GetBaseFeePerGas())
	require.Equal(t, payload.GetBlockHash(), header.GetBlockHash())
	require.Equal(t, payload.GetBlobGasUsed(), header.GetBlobGasUsed())
	require.Equal(t, payload.GetExcessBlobGas(), header.GetExcessBlobGas())
	require.Equal(t, payload.GetForkVersion(), header.GetForkVersion())

	require.Equal(t, payload.HashTreeRoot(), header.HashTreeRoot())
}

func TestExecutionPayload_UnmarshalJSON_Error(t *testing.T) {
	t.Parallel()
	original := generateExecutionPayload()
	validJSON, err := original.MarshalJSON()
	require.NoError(t, err)

	testCases := []struct {
		name          string
		removeField   string
		expectedError string
	}{
		{
			name:          "missing required field 'parentHash'",
			removeField:   "parentHash",
			expectedError: "missing required field 'parentHash' for ExecutionPayload",
		},
		{
			name:          "missing required field 'feeRecipient'",
			removeField:   "feeRecipient",
			expectedError: "missing required field 'feeRecipient' for ExecutionPayload",
		},
		{
			name:          "missing required field 'stateRoot'",
			removeField:   "stateRoot",
			expectedError: "missing required field 'stateRoot' for ExecutionPayload",
		},
		{
			name:          "missing required field 'receiptsRoot'",
			removeField:   "receiptsRoot",
			expectedError: "missing required field 'receiptsRoot' for ExecutionPayload",
		},
		{
			name:          "missing required field 'logsBloom'",
			removeField:   "logsBloom",
			expectedError: "missing required field 'logsBloom' for ExecutionPayload",
		},
		{
			name:          "missing required field 'prevRandao'",
			removeField:   "prevRandao",
			expectedError: "missing required field 'prevRandao' for ExecutionPayload",
		},
		{
			name:          "missing required field 'blockNumber'",
			removeField:   "blockNumber",
			expectedError: "missing required field 'blockNumber' for ExecutionPayload",
		},
		{
			name:          "missing required field 'gasLimit'",
			removeField:   "gasLimit",
			expectedError: "missing required field 'gasLimit' for ExecutionPayload",
		},
		{
			name:          "missing required field 'gasUsed'",
			removeField:   "gasUsed",
			expectedError: "missing required field 'gasUsed' for ExecutionPayload",
		},
		{
			name:          "missing required field 'timestamp'",
			removeField:   "timestamp",
			expectedError: "missing required field 'timestamp' for ExecutionPayload",
		},
		{
			name:          "missing required field 'extraData'",
			removeField:   "extraData",
			expectedError: "missing required field 'extraData' for ExecutionPayload",
		},
		{
			name:          "missing required field 'baseFeePerGas'",
			removeField:   "baseFeePerGas",
			expectedError: "missing required field 'baseFeePerGas' for ExecutionPayload",
		},
		{
			name:          "missing required field 'blockHash'",
			removeField:   "blockHash",
			expectedError: "missing required field 'blockHash' for ExecutionPayload",
		},
		{
			name:          "missing required field 'transactions'",
			removeField:   "transactions",
			expectedError: "missing required field 'transactions' for ExecutionPayload",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var payload types.ExecutionPayload
			var jsonMap map[string]interface{}

			errUnmarshal := json.Unmarshal(validJSON, &jsonMap)
			require.NoError(t, errUnmarshal)

			delete(jsonMap, tc.removeField)

			malformedJSON, errMarshal := json.Marshal(jsonMap)
			require.NoError(t, errMarshal)

			err = payload.UnmarshalJSON(malformedJSON)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedError)
		})
	}
}
