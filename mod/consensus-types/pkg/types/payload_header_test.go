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
	"encoding/json"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/stretchr/testify/require"
)

func generateExecutionPayloadHeaderDeneb() *types.ExecutionPayloadHeaderDeneb {
	return &types.ExecutionPayloadHeaderDeneb{
		ParentHash:       common.ExecutionHash{},
		FeeRecipient:     common.ExecutionAddress{},
		StateRoot:        bytes.B32{},
		ReceiptsRoot:     bytes.B32{},
		LogsBloom:        make([]byte, 256),
		Random:           bytes.B32{},
		Number:           math.U64(0),
		GasLimit:         math.U64(0),
		GasUsed:          math.U64(0),
		Timestamp:        math.U64(0),
		ExtraData:        []byte{},
		BaseFeePerGas:    math.Wei{},
		BlockHash:        common.ExecutionHash{},
		TransactionsRoot: bytes.B32{},
		WithdrawalsRoot:  bytes.B32{},
		BlobGasUsed:      math.U64(0),
		ExcessBlobGas:    math.U64(0),
	}
}

func TestExecutionPayloadHeaderDeneb_Getters(t *testing.T) {
	header := generateExecutionPayloadHeaderDeneb()

	require.NotNil(t, header)

	require.Equal(t, common.ExecutionHash{}, header.GetParentHash())
	require.Equal(t, common.ExecutionAddress{}, header.GetFeeRecipient())
	require.Equal(t, bytes.B32{}, header.GetStateRoot())
	require.Equal(t, bytes.B32{}, header.GetReceiptsRoot())
	require.Equal(t, make([]byte, 256), header.GetLogsBloom())
	require.Equal(t, bytes.B32{}, header.GetPrevRandao())
	require.Equal(t, math.U64(0), header.GetNumber())
	require.Equal(t, math.U64(0), header.GetGasLimit())
	require.Equal(t, math.U64(0), header.GetGasUsed())
	require.Equal(t, math.U64(0), header.GetTimestamp())
	require.Equal(t, []byte{}, header.GetExtraData())
	require.Equal(t, math.Wei{}, header.GetBaseFeePerGas())
	require.Equal(t, common.ExecutionHash{}, header.GetBlockHash())
	require.Equal(t, bytes.B32{}, header.GetTransactionsRoot())
	require.Equal(t, bytes.B32{}, header.GetWithdrawalsRoot())
	require.Equal(t, math.U64(0), header.GetBlobGasUsed())
	require.Equal(t, math.U64(0), header.GetExcessBlobGas())
}

func TestExecutionPayloadHeaderDeneb_IsBlinded(t *testing.T) {
	header := generateExecutionPayloadHeaderDeneb()
	require.False(t, header.IsBlinded())
}

func TestExecutionPayloadHeaderDeneb_IsNil(t *testing.T) {
	header := generateExecutionPayloadHeaderDeneb()
	require.False(t, header.IsNil())
}

func TestExecutionPayloadHeaderDeneb_Version(t *testing.T) {
	header := generateExecutionPayloadHeaderDeneb()
	require.Equal(t, version.Deneb, header.Version())
}

func TestExecutionPayloadHeaderDeneb_MarshalUnmarshalJSON(t *testing.T) {
	originalHeader := generateExecutionPayloadHeaderDeneb()

	data, err := originalHeader.MarshalJSON()
	require.NoError(t, err)
	require.NotNil(t, data)

	var header types.ExecutionPayloadHeaderDeneb
	err = header.UnmarshalJSON(data)
	require.NoError(t, err)

	require.Equal(t, originalHeader, &header)
}

func TestExecutionPayloadHeaderDeneb_Serialization(t *testing.T) {
	original := generateExecutionPayloadHeaderDeneb()

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.ExecutionPayloadHeaderDeneb
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)

	require.Equal(t, original, &unmarshalled)
}

func TestExecutionPayloadHeaderDeneb_MarshalSSZTo(t *testing.T) {
	header := generateExecutionPayloadHeaderDeneb()

	buf := make([]byte, 0)
	data, err := header.MarshalSSZTo(buf)
	require.NoError(t, err)
	require.NotNil(t, data)
}

func TestExecutionPayloadHeaderDeneb_SizeSSZ(t *testing.T) {
	header := generateExecutionPayloadHeaderDeneb()
	size := header.SizeSSZ()
	require.Equal(t, 584, size)
}

func TestExecutionPayloadHeaderDeneb_HashTreeRoot(t *testing.T) {
	header := generateExecutionPayloadHeaderDeneb()
	_, err := header.HashTreeRoot()
	require.NoError(t, err)
}

func TestExecutionPayloadHeaderDeneb_GetTree(t *testing.T) {
	header := generateExecutionPayloadHeaderDeneb()
	_, err := header.GetTree()
	require.NoError(t, err)
}

func TestExecutionPayloadHeaderDeneb_Empty(t *testing.T) {
	header := new(types.ExecutionPayloadHeader)
	emptyHeader := header.Empty(version.Deneb)

	require.NotNil(t, emptyHeader)
	require.Equal(t, version.Deneb, emptyHeader.Version())
}

//nolint:lll
func TestExecutablePayloadHeaderDeneb_UnmarshalJSON_Error(t *testing.T) {
	original := generateExecutionPayloadHeaderDeneb()
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
			expectedError: "missing required field 'parentHash' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'feeRecipient'",
			removeField:   "feeRecipient",
			expectedError: "missing required field 'feeRecipient' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'stateRoot'",
			removeField:   "stateRoot",
			expectedError: "missing required field 'stateRoot' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'receiptsRoot'",
			removeField:   "receiptsRoot",
			expectedError: "missing required field 'receiptsRoot' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'logsBloom'",
			removeField:   "logsBloom",
			expectedError: "missing required field 'logsBloom' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'prevRandao'",
			removeField:   "prevRandao",
			expectedError: "missing required field 'prevRandao' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'blockNumber'",
			removeField:   "blockNumber",
			expectedError: "missing required field 'blockNumber' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'gasLimit'",
			removeField:   "gasLimit",
			expectedError: "missing required field 'gasLimit' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'gasUsed'",
			removeField:   "gasUsed",
			expectedError: "missing required field 'gasUsed' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'timestamp'",
			removeField:   "timestamp",
			expectedError: "missing required field 'timestamp' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'extraData'",
			removeField:   "extraData",
			expectedError: "missing required field 'extraData' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'baseFeePerGas'",
			removeField:   "baseFeePerGas",
			expectedError: "missing required field 'baseFeePerGas' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'blockHash'",
			removeField:   "blockHash",
			expectedError: "missing required field 'blockHash' for ExecutionPayloadHeaderDeneb",
		},
		{
			name:          "missing required field 'transactionsRoot'",
			removeField:   "transactionsRoot",
			expectedError: "missing required field 'transactionsRoot' for ExecutionPayloadHeaderDeneb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var payload types.ExecutionPayloadHeaderDeneb
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
