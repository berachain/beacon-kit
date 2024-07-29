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

package types_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/stretchr/testify/require"
)

func generateExecutionPayload() *types.ExecutionPayload {
	transactions := make([][]byte, 1)
	transactions[0] = []byte{0x07}
	withdrawals := make([]*engineprimitives.Withdrawal, 1)
	withdrawals[0] = &engineprimitives.Withdrawal{
		Index:     0,
		Validator: 0,
		Address:   common.ExecutionAddress{},
		Amount:    0,
	}
	return &types.ExecutionPayload{
		ParentHash:    gethprimitives.ExecutionHash{},
		FeeRecipient:  gethprimitives.ExecutionAddress{},
		StateRoot:     bytes.B32{},
		ReceiptsRoot:  bytes.B32{},
		LogsBloom:     bytes.B256{},
		Random:        bytes.B32{},
		Number:        math.U64(0),
		GasLimit:      math.U64(0),
		GasUsed:       math.U64(0),
		Timestamp:     math.U64(0),
		ExtraData:     []byte{0x01},
		BaseFeePerGas: math.Wei{},
		BlockHash:     gethprimitives.ExecutionHash{},
		Transactions:  transactions,
		Withdrawals:   withdrawals,
		BlobGasUsed:   math.U64(0),
		ExcessBlobGas: math.U64(0),
	}
}
func TestExecutionPayload_Serialization(t *testing.T) {
	original := generateExecutionPayload()

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.ExecutionPayload
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)

	require.Equal(t, original, &unmarshalled)
}

func TestExecutionPayload_SizeSSZ(t *testing.T) {
	payload := generateExecutionPayload()
	size := payload.SizeSSZ(false)
	require.Equal(t, uint32(578), size)

	state := &types.ExecutionPayload{}
	err := state.UnmarshalSSZ([]byte{0x01, 0x02, 0x03}) // Invalid data
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestExecutionPayload_HashTreeRoot(t *testing.T) {
	payload := generateExecutionPayload()
	_, err := payload.HashTreeRoot()
	require.NoError(t, err)
}

func TestExecutionPayload_GetTree(t *testing.T) {
	payload := generateExecutionPayload()
	tree, err := payload.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestExecutionPayload_Getters(t *testing.T) {
	payload := generateExecutionPayload()
	require.Equal(t, gethprimitives.ExecutionHash{}, payload.GetParentHash())
	require.Equal(
		t,
		gethprimitives.ExecutionAddress{},
		payload.GetFeeRecipient(),
	)

	transactions := make([][]byte, 1)
	transactions[0] = []byte{0x07}
	withdrawals := make([]*engineprimitives.Withdrawal, 1)
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
	require.Equal(t, math.Wei{}, payload.GetBaseFeePerGas())
	require.Equal(t, gethprimitives.ExecutionHash{}, payload.GetBlockHash())
	require.Equal(t, transactions, payload.GetTransactions())
	require.Equal(t, withdrawals, payload.GetWithdrawals())
	require.Equal(t, math.U64(0), payload.GetBlobGasUsed())
	require.Equal(t, math.U64(0), payload.GetExcessBlobGas())
}

func TestExecutionPayload_MarshalJSON(t *testing.T) {
	payload := generateExecutionPayload()

	data, err := payload.MarshalJSON()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.ExecutionPayload
	err = unmarshalled.UnmarshalJSON(data)
	require.NoError(t, err)
	require.Equal(t, payload, &unmarshalled)
}

func TestExecutionPayload_IsNil(t *testing.T) {
	var payload *types.ExecutionPayload
	require.True(t, payload.IsNil())

	payload = generateExecutionPayload()
	require.False(t, payload.IsNil())
}

func TestExecutionPayload_IsBlinded(t *testing.T) {
	payload := generateExecutionPayload()
	require.False(t, payload.IsBlinded())
}

func TestExecutionPayload_Version(t *testing.T) {
	payload := generateExecutionPayload()
	require.Equal(t, version.Deneb, payload.Version())
}

func TestExecutionPayload_Empty(t *testing.T) {
	payload := new(types.ExecutionPayload)
	emptyPayload := payload.Empty(version.Deneb)

	require.NotNil(t, emptyPayload)
	require.Equal(t, version.Deneb, emptyPayload.Version())
}

func TestExecutionPayload_ToHeader(t *testing.T) {
	payload := &types.ExecutionPayload{
		ParentHash:    gethprimitives.ExecutionHash{},
		FeeRecipient:  gethprimitives.ExecutionAddress{},
		StateRoot:     bytes.B32{},
		ReceiptsRoot:  bytes.B32{},
		LogsBloom:     bytes.B256{},
		Random:        bytes.B32{},
		Number:        math.U64(0),
		GasLimit:      math.U64(0),
		GasUsed:       math.U64(0),
		Timestamp:     math.U64(0),
		ExtraData:     []byte{},
		BaseFeePerGas: math.Wei{},
		BlockHash:     gethprimitives.ExecutionHash{},
		Transactions:  [][]byte{[]byte{0x01}},
		Withdrawals:   []*engineprimitives.Withdrawal{},
		BlobGasUsed:   math.U64(0),
		ExcessBlobGas: math.U64(0),
	}

	header, err := payload.ToHeader(
		merkle.NewMerkleizer[[32]byte, *ssz.List[ssz.Byte]](),
		uint64(16), uint64(80087),
	)
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

	// TODO: FIX LATER
	// htrHeader, err := header.HashTreeRoot()
	// require.NoError(t, err)
	// htrPayload, err := payload.HashTreeRoot()
	// require.NoError(t, err)
	// require.Equal(t, htrPayload, htrHeader)
}

func TestExecutionPayload_UnmarshalJSON_Error(t *testing.T) {
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
