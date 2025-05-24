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
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	sszutil "github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// generateExecutionPayloadHeader generates an ExecutionPayloadHeader.
func generateExecutionPayloadHeader() *types.ExecutionPayloadHeader {
	return &types.ExecutionPayloadHeader{
		ParentHash:       common.ExecutionHash{},
		FeeRecipient:     common.ExecutionAddress{},
		StateRoot:        bytes.B32{},
		ReceiptsRoot:     bytes.B32{},
		LogsBloom:        bytes.B256{},
		Random:           bytes.B32{},
		Number:           math.U64(0),
		GasLimit:         math.U64(0),
		GasUsed:          math.U64(0),
		Timestamp:        math.U64(0),
		ExtraData:        nil,
		BaseFeePerGas:    &math.U256{},
		BlockHash:        common.ExecutionHash{},
		TransactionsRoot: common.Root{},
		WithdrawalsRoot:  common.Root{},
		BlobGasUsed:      math.U64(0),
		ExcessBlobGas:    math.U64(0),
	}
}

func TestExecutionPayloadHeader_Getters(t *testing.T) {
	t.Parallel()
	header := generateExecutionPayloadHeader()
	require.NotNil(t, header)

	require.Equal(t, common.ExecutionHash{}, header.GetParentHash())
	require.Equal(
		t,
		common.ExecutionAddress{},
		header.GetFeeRecipient(),
	)
	require.Equal(t, bytes.B32{}, header.GetStateRoot())
	require.Equal(t, bytes.B32{}, header.GetReceiptsRoot())
	require.Equal(t, bytes.B256{}, header.GetLogsBloom())
	require.Equal(t, bytes.B32{}, header.GetPrevRandao())
	require.Equal(t, math.U64(0), header.GetNumber())
	require.Equal(t, math.U64(0), header.GetGasLimit())
	require.Equal(t, math.U64(0), header.GetGasUsed())
	require.Equal(t, math.U64(0), header.GetTimestamp())
	require.Equal(t, []byte(nil), header.GetExtraData())
	require.Equal(t, math.NewU256(0), header.GetBaseFeePerGas())
	require.Equal(t, common.ExecutionHash{}, header.GetBlockHash())
	require.Equal(t, common.Root{}, header.GetTransactionsRoot())
	require.Equal(t, common.Root{}, header.GetWithdrawalsRoot())
	require.Equal(t, math.U64(0), header.GetBlobGasUsed())
	require.Equal(t, math.U64(0), header.GetExcessBlobGas())
}

func TestExecutionPayloadHeader_IsNil(t *testing.T) {
	t.Parallel()
	header := generateExecutionPayloadHeader()
	require.NotNil(t, header)
}

func TestExecutionPayloadHeader_MarshalUnmarshalJSON(t *testing.T) {
	t.Parallel()
	originalHeader := generateExecutionPayloadHeader()

	data, err := originalHeader.MarshalJSON()
	require.NoError(t, err)
	require.NotNil(t, data)

	var header types.ExecutionPayloadHeader
	err = header.UnmarshalJSON(data)
	require.NoError(t, err)
	require.Equal(t, *originalHeader, header)
}

func TestExecutionPayloadHeader_Serialization(t *testing.T) {
	t.Parallel()
	original := generateExecutionPayloadHeader()

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	unmarshalled := types.NewEmptyExecutionPayloadHeader()
	err = sszutil.Unmarshal(data, unmarshalled)
	require.NoError(t, err)
	require.Equal(t, *original, *unmarshalled)
}

func TestExecutionPayloadHeader_MarshalSSZTo(t *testing.T) {
	t.Parallel()
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
		{
			name: "invalid extra data passes marshalling",
			malleate: func() *types.ExecutionPayloadHeader {
				header := generateExecutionPayloadHeader()
				header.ExtraData = make([]byte, 100)
				return header
			},
			expErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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

func TestExecutionPayloadHeader_NewFromSSZ_EmptyBuf(t *testing.T) {
	t.Parallel()
	buf := make([]byte, 0)
	header := types.NewEmptyExecutionPayloadHeader()
	err := sszutil.Unmarshal(buf, header)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestExecutionPayloadHeader_NewFromSSZ_Invalid(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name     string
		malleate func() []byte
		expErr   error
	}{
		{
			name: "offset exceeds length",
			malleate: func() []byte {
				header := generateExecutionPayloadHeader()
				buf, err := header.MarshalSSZ()
				require.NoError(t, err)

				buf[436] = 10
				buf[437] = 10
				buf[438] = 10
				buf[439] = 10
				return buf
			},
			expErr: ssz.ErrOffsetBeyondCapacity,
		},
		{
			name: "invalid extra data: extra data too large",
			malleate: func() []byte {
				header := generateExecutionPayloadHeader()
				buf, err := header.MarshalSSZ()

				// add dummy extra data to exceed the 32 limit
				dummyExtra := make([]byte, 100)
				buf = append(buf, dummyExtra...)
				require.NoError(t, err)
				return buf
			},
			expErr: ssz.ErrMaxLengthExceeded,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			buf := tc.malleate()
			dest := types.NewEmptyExecutionPayloadHeader()
			err := sszutil.Unmarshal(buf, dest)
			require.ErrorIs(t, err, tc.expErr)
		})
	}
}

func TestExecutionPayloadHeader_NewFromSSZ_Invalid_TooSmall(t *testing.T) {
	t.Parallel()
	header := generateExecutionPayloadHeader()
	buf, err := header.MarshalSSZ()
	require.NoError(t, err)

	buf[436] = 1
	buf[437] = 0
	buf[438] = 0
	buf[439] = 0

	dest := types.NewEmptyExecutionPayloadHeader()
	err = sszutil.Unmarshal(buf, dest)
	require.Error(t, err)

	// Can be either ErrFirstOffsetMismatch or ErrBadOffsetProgression due to reused Decoder in
	// SSZ lib. If the SSZ lib happens to grab a reused Decoder from the decoderPool, the decoder's
	// `offsets` field is already initialized to an empty slice instead of nil. This passes the nil
	// check in the ErrFirstOffsetMismatch error condition resulting in no error. Immediately after,
	// it will still fail the ErrBadOffsetProgression error condition. This flakiness depends upon
	// retrieving a used Decoder from the decoderPool as well as the intentional misuse of the
	// marshaled data. In the case that an actor intentionally tries to induce this behavior, the
	// unmarshaling of the data correctly results in error, just a different error.
	// In this unit test, we simply expect the error to be one of the two possible errors.
	isExpectedError := errors.IsAny(err, ssz.ErrFirstOffsetMismatch, ssz.ErrBadOffsetProgression)
	require.True(
		t, isExpectedError, "expected %w or %w, got %w",
		ssz.ErrFirstOffsetMismatch, ssz.ErrBadOffsetProgression, err,
	)
}

func TestExecutionPayloadHeader_SizeSSZ(t *testing.T) {
	t.Parallel()
	header := generateExecutionPayloadHeader()
	size := ssz.Size(header)
	require.Equal(t, types.ExecutionPayloadHeaderStaticSize, size)
}

func TestExecutionPayloadHeader_HashTreeRoot(t *testing.T) {
	t.Parallel()
	header := generateExecutionPayloadHeader()
	require.NotPanics(t, func() {
		header.HashTreeRoot()
	})
}

func TestExecutionPayloadHeader_GetTree(t *testing.T) {
	t.Parallel()
	header := generateExecutionPayloadHeader()
	_, err := header.GetTree()
	require.NoError(t, err)
}

func TestExecutablePayloadHeader_UnmarshalJSON_Error(t *testing.T) {
	t.Parallel()
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

func TestExecutablePayloadHeader_UnmarshalJSON_Empty(t *testing.T) {
	t.Parallel()
	var payload types.ExecutionPayloadHeader
	err := payload.UnmarshalJSON([]byte{})
	require.Error(t, err)
}

func TestExecutablePayloadHeader_HashTreeRootWith(t *testing.T) {
	t.Parallel()
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
			expErr: fastssz.ErrIncorrectListSize,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hh := fastssz.DefaultHasherPool.Get()
			header := tc.malleate()
			err := header.HashTreeRootWith(hh)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestExecutionPayloadHeader_NewFromSSZ(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		data           []byte
		expErr         error
		expectedHeader *types.ExecutionPayloadHeader
	}{
		{
			name: "Valid SSZ data",
			data: func() []byte {
				data, _ := generateExecutionPayloadHeader().MarshalSSZ()
				return data
			}(),
			expErr:         nil,
			expectedHeader: generateExecutionPayloadHeader(),
		},
		{
			name:           "Invalid SSZ data",
			data:           []byte{0x01, 0x02},
			expErr:         io.ErrUnexpectedEOF,
			expectedHeader: nil,
		},
		{
			name:           "Empty SSZ data",
			data:           []byte{},
			expErr:         io.ErrUnexpectedEOF,
			expectedHeader: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			header := types.NewEmptyExecutionPayloadHeader()
			if tc.name == "Different fork version" {
				require.Panics(t, func() {
					_ = sszutil.Unmarshal(tc.data, header)
				}, "Expected panic for different fork version")
			} else {
				err := sszutil.Unmarshal(tc.data, header)
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
	t.Parallel()
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
			t.Parallel()
			header := types.NewEmptyExecutionPayloadHeader()
			err := json.Unmarshal(tc.data, header)

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
