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
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

func generateExecutionPayloadHeader() *types.ExecutionPayloadHeader {
	return &types.ExecutionPayloadHeader{
		ParentHash:       gethprimitives.ExecutionHash{},
		FeeRecipient:     gethprimitives.ExecutionAddress{},
		StateRoot:        bytes.B32{123},
		ReceiptsRoot:     bytes.B32{123},
		LogsBloom:        bytes.B256{123},
		Random:           bytes.B32{123},
		Number:           math.U64(0),
		GasLimit:         math.U64(0),
		GasUsed:          math.U64(0),
		Timestamp:        math.U64(0),
		ExtraData:        nil,
		BaseFeePerGas:    &math.U256{},
		BlockHash:        gethprimitives.ExecutionHash{},
		TransactionsRoot: bytes.B32{},
		WithdrawalsRoot:  bytes.B32{},
		BlobGasUsed:      math.U64(0),
		ExcessBlobGas:    math.U64(0),
	}
}

func TestExecutionPayloadHeader_Getters(t *testing.T) {
	header := generateExecutionPayloadHeader()

	require.NotNil(t, header)

	require.Equal(t, gethprimitives.ExecutionHash{}, header.GetParentHash())
	require.Equal(
		t,
		gethprimitives.ExecutionAddress{},
		header.GetFeeRecipient(),
	)
	require.Equal(t, bytes.B32{123}, header.GetStateRoot())
	require.Equal(t, bytes.B32{123}, header.GetReceiptsRoot())
	require.Equal(t, bytes.B256{123}, header.GetLogsBloom())
	require.Equal(t, bytes.B32{123}, header.GetPrevRandao())
	require.Equal(t, math.U64(0), header.GetNumber())
	require.Equal(t, math.U64(0), header.GetGasLimit())
	require.Equal(t, math.U64(0), header.GetGasUsed())
	require.Equal(t, math.U64(0), header.GetTimestamp())
	require.Equal(t, []byte(nil), header.GetExtraData())
	require.Equal(t, math.NewU256(0), header.GetBaseFeePerGas())
	require.Equal(t, gethprimitives.ExecutionHash{}, header.GetBlockHash())
	require.Equal(t, bytes.B32{}, header.GetTransactionsRoot())
	require.Equal(t, bytes.B32{}, header.GetWithdrawalsRoot())
	require.Equal(t, math.U64(0), header.GetBlobGasUsed())
	require.Equal(t, math.U64(0), header.GetExcessBlobGas())
}

func TestExecutionPayloadHeader_IsNil(t *testing.T) {
	header := generateExecutionPayloadHeader()
	require.False(t, header.IsNil())
}

func TestExecutionPayloadHeader_Version(t *testing.T) {
	header := generateExecutionPayloadHeader()
	require.Equal(t, version.Deneb, header.Version())
}

func TestExecutionPayloadHeader_MarshalUnmarshalJSON(t *testing.T) {
	originalHeader := generateExecutionPayloadHeader()

	data, err := originalHeader.MarshalJSON()
	require.NoError(t, err)
	require.NotNil(t, data)

	var header types.ExecutionPayloadHeader
	err = header.UnmarshalJSON(data)
	require.NoError(t, err)

	require.Equal(t, originalHeader, &header)
}

func TestExecutionPayloadHeader_Serialization(t *testing.T) {
	original := generateExecutionPayloadHeader()

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled = new(types.ExecutionPayloadHeader).Empty(version.Deneb)
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.Equal(t, original, unmarshalled)
}

func TestExecutionPayloadHeader_MarshalSSZTo(t *testing.T) {
	testcases := []struct {
		name     string
		malleate func() *types.ExecutionPayloadHeader
		expErr   error
	}{
		{
			name:     "valid",
			malleate: generateExecutionPayloadHeader,
			expErr:   nil,
		},
		// TODO: Is this okay?
		// {
		// 	name: "invalid extra data",
		// 	malleate: func() *types.ExecutionPayloadHeader {
		// 		header := generateExecutionPayloadHeader()
		// 		header.ExtraData = make([]byte, 100)
		// 		return header
		// 	},
		// 	expErr: ssz.ErrBytesLengthFn(extraDataField, 100, 32),
		// },
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			header := tc.malleate()
			buf := make([]byte, 64)
			_, err := header.MarshalSSZTo(buf)
			if tc.expErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestExecutionPayloadHeader_UnmarshalSSZ_EmptyBuf(t *testing.T) {
	header := generateExecutionPayloadHeader()
	buf := make([]byte, 0)
	err := header.UnmarshalSSZ(buf)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

// func TestExecutionPayloadHeader_UnmarshalSSZ(t *testing.T) {
// 	testcases := []struct {
// 		name     string
// 		malleate func() []byte
// 		expErr   error
// 	}{
// 		{
// 			name: "offset exceeds length",
// 			malleate: func() []byte {
// 				header := generateExecutionPayloadHeader()
// 				buf, err := header.MarshalSSZ()
// 				require.NoError(t, err)

// 				buf[436] = 10
// 				buf[437] = 10
// 				buf[438] = 10
// 				buf[439] = 10
// 				return buf
// 			},
// 			expErr: ssz.ErrOffset,
// 		},
// 		{
// 			name: "invalid extra data: offset too small",
// 			malleate: func() []byte {
// 				header := generateExecutionPayloadHeader()
// 				buf, err := header.MarshalSSZ()
// 				require.NoError(t, err)

// 				buf[436] = 1
// 				buf[437] = 0
// 				buf[438] = 0
// 				buf[439] = 0
// 				return buf
// 			},
// 			expErr: ssz.ErrInvalidVariableOffset,
// 		},
// 		{
// 			name: "invalid extra data: extra data too large",
// 			malleate: func() []byte {
// 				header := generateExecutionPayloadHeader()
// 				buf, err := header.MarshalSSZ()

// 				// add dummy extra data to exceed the 32 limit
// 				dummyExtra := make([]byte, 100)
// 				buf = append(buf, dummyExtra...)
// 				require.NoError(t, err)
// 				return buf
// 			},
// 			expErr: ssz.ErrBytesLength,
// 		},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var header types.ExecutionPayloadHeader
// 			buf := tc.malleate()
// 			err := header.UnmarshalSSZ(buf)
// 			require.Equal(t, tc.expErr, err)
// 		})
// 	}
// }

func TestExecutionPayloadHeader_SizeSSZ(t *testing.T) {
	header := generateExecutionPayloadHeader()
	size := header.SizeSSZ(true)
	require.Equal(t, uint32(584), size)
}

func TestExecutionPayloadHeader_HashTreeRoot(t *testing.T) {
	header := generateExecutionPayloadHeader()
	_, err := header.HashTreeRoot()
	require.NoError(t, err)
}

func TestExecutionPayloadHeader_GetTree(t *testing.T) {
	header := generateExecutionPayloadHeader()
	_, err := header.GetTree()
	require.NoError(t, err)
}

func TestExecutionPayloadHeader_Empty(t *testing.T) {
	header := new(types.ExecutionPayloadHeader)
	emptyHeader := header.Empty(version.Deneb)

	require.NotNil(t, emptyHeader)
	require.Equal(t, version.Deneb, emptyHeader.Version())
}

//nolint:lll
func TestExecutablePayloadHeaderDeneb_UnmarshalJSON_Error(t *testing.T) {
	original := generateExecutionPayloadHeader()
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
			expectedError: "missing required field 'parentHash' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'feeRecipient'",
			removeField:   "feeRecipient",
			expectedError: "missing required field 'feeRecipient' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'stateRoot'",
			removeField:   "stateRoot",
			expectedError: "missing required field 'stateRoot' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'receiptsRoot'",
			removeField:   "receiptsRoot",
			expectedError: "missing required field 'receiptsRoot' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'logsBloom'",
			removeField:   "logsBloom",
			expectedError: "missing required field 'logsBloom' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'prevRandao'",
			removeField:   "prevRandao",
			expectedError: "missing required field 'prevRandao' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'blockNumber'",
			removeField:   "blockNumber",
			expectedError: "missing required field 'blockNumber' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'gasLimit'",
			removeField:   "gasLimit",
			expectedError: "missing required field 'gasLimit' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'gasUsed'",
			removeField:   "gasUsed",
			expectedError: "missing required field 'gasUsed' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'timestamp'",
			removeField:   "timestamp",
			expectedError: "missing required field 'timestamp' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'extraData'",
			removeField:   "extraData",
			expectedError: "missing required field 'extraData' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'baseFeePerGas'",
			removeField:   "baseFeePerGas",
			expectedError: "missing required field 'baseFeePerGas' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'blockHash'",
			removeField:   "blockHash",
			expectedError: "missing required field 'blockHash' for ExecutionPayloadHeader",
		},
		{
			name:          "missing required field 'transactionsRoot'",
			removeField:   "transactionsRoot",
			expectedError: "missing required field 'transactionsRoot' for ExecutionPayloadHeader",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var payload types.ExecutionPayloadHeader
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

func TestExecutablePayloadHeaderDeneb_UnmarshalJSON_Empty(t *testing.T) {
	var payload types.ExecutionPayloadHeader
	err := payload.UnmarshalJSON([]byte{})
	require.Error(t, err)
}

func TestExecutablePayloadHeaderDeneb_HashTreeRootWith(t *testing.T) {
	testcases := []struct {
		name     string
		malleate func() *types.ExecutionPayloadHeader
		expErr   error
	}{
		{
			name: "invalid ExtraData length",
			malleate: func() *types.ExecutionPayloadHeader {
				var header = generateExecutionPayloadHeader()
				header.ExtraData = make([]byte, 50)
				return header
			},
			expErr: ssz.ErrIncorrectListSize,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			hh := ssz.DefaultHasherPool.Get()
			header := tc.malleate()
			err := header.HashTreeRootWith(hh)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestExecutionPayloadHeader_NewFromSSZ(t *testing.T) {
	t.Helper()
	testCases := []struct {
		name           string
		data           []byte
		forkVersion    uint32
		expErr         error
		expectedHeader *types.ExecutionPayloadHeader
	}{
		{
			name: "Valid SSZ data",
			data: func() []byte {
				data, _ := generateExecutionPayloadHeader().MarshalSSZ()
				return data
			}(),
			forkVersion:    version.Deneb,
			expErr:         nil,
			expectedHeader: generateExecutionPayloadHeader(),
		},
		{
			name:           "Invalid SSZ data",
			data:           []byte{0x01, 0x02},
			forkVersion:    version.Deneb,
			expErr:         io.ErrUnexpectedEOF,
			expectedHeader: nil,
		},
		{
			name:           "Empty SSZ data",
			data:           []byte{},
			forkVersion:    version.Deneb,
			expErr:         io.ErrUnexpectedEOF,
			expectedHeader: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "Different fork version" {
				require.Panics(t, func() {
					_, _ = new(types.ExecutionPayloadHeader).
						NewFromSSZ(tc.data, tc.forkVersion)
				}, "Expected panic for different fork version")
			} else {
				header, err := new(types.ExecutionPayloadHeader).
					NewFromSSZ(tc.data, tc.forkVersion)
				if tc.expErr != nil {
					require.ErrorIs(t, err, tc.expErr)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedHeader, header)
				}
			}
		})
	}
}

func TestExecutionPayloadHeader_NewFromJSON(t *testing.T) {
	t.Helper()
	type testCase struct {
		name          string
		data          []byte
		header        *types.ExecutionPayloadHeader
		expectedError error
	}
	testCases := []testCase{
		func() testCase {
			header := generateExecutionPayloadHeader()
			return testCase{
				name:   "Valid JSON",
				header: header,
				data: func() []byte {
					data, err := json.Marshal(header)
					require.NoError(t, err)
					return data
				}(),
			}
		}(),
		{
			name:          "Invalid JSON",
			data:          []byte{},
			expectedError: errors.New("unexpected end of JSON input"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			header, err := new(types.ExecutionPayloadHeader).NewFromJSON(
				tc.data,
				version.Deneb,
			)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			if tc.header != nil {
				require.Equal(t, tc.header, header)
			}
		})
	}
}
