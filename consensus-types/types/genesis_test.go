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

//nolint:lll // long strings.
package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenesis(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		g := types.DefaultGenesis(v)
		if !version.Equals(g.ForkVersion, v) {
			t.Errorf(
				"Expected fork version %v, but got %v",
				v,
				g.ForkVersion,
			)
		}

		if len(g.Deposits) != 0 {
			t.Errorf("Expected no deposits, but got %v", len(g.Deposits))
		}
		// add assertions for ExecutionPayloadHeader
		require.NotNil(t, g.ExecutionPayloadHeader,
			"Expected ExecutionPayloadHeader to be non-nil")
		require.Equal(t, common.ExecutionHash{},
			g.ExecutionPayloadHeader.GetParentHash(),
			"Unexpected ParentHash")
		require.Equal(t, common.ExecutionAddress{},
			g.ExecutionPayloadHeader.GetFeeRecipient(),
			"Unexpected FeeRecipient")
		require.Equal(t, math.U64(30000000),
			g.ExecutionPayloadHeader.GetGasLimit(),
			"Unexpected GasLimit")
		require.Equal(t, math.U64(0),
			g.ExecutionPayloadHeader.GetGasUsed(),
			"Unexpected GasUsed")
		require.Equal(t, math.U64(0),
			g.ExecutionPayloadHeader.GetTimestamp(),
			"Unexpected Timestamp")
	})
}

func TestDefaultGenesisExecutionPayloadHeader(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		header, err := types.DefaultGenesisExecutionPayloadHeader(v)
		require.NoError(t, err)
		require.NotNil(t, header)
	})
}

func TestGenesisGetForkVersion(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		g := types.DefaultGenesis(v)
		forkVersion := g.GetForkVersion()
		require.Equal(t, v, forkVersion)
	})
}

func TestGenesisGetDeposits(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		g := types.DefaultGenesis(v)
		deposits := g.GetDeposits()
		require.Empty(t, deposits)
	})
}

func TestGenesisGetExecutionPayloadHeader(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		g := types.DefaultGenesis(v)
		header := g.GetExecutionPayloadHeader()
		require.NotNil(t, header)
	})
}

func TestDefaultGenesisPanics(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		require.NotPanics(t, func() {
			types.DefaultGenesis(v)
		})
	})
}

func TestGenesisUnmarshalJSON(t *testing.T) {
	t.Parallel()
	t.Helper()
	testCases := []struct {
		name                string
		jsonInput           string
		expectedError       bool
		expectedFork        bytes.B4
		expectedDepositsLen int
	}{
		{
			name: "Valid JSON with empty deposits",
			jsonInput: `{
				  "fork_version": "0x04000000",
				  "deposits": [],
				  "execution_payload_header": {
					"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"feeRecipient": "0x0000000000000000000000000000000000000000",
					"stateRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"receiptsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
					"prevRandao": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"blockNumber": "0x0",
					"gasLimit": "0x0",
					"gasUsed": "0x0",
					"timestamp": "0x0",
					"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"baseFeePerGas": "0x0",
					"blockHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"withdrawalsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"blobGasUsed": "0x0",
					"excessBlobGas": "0x0"
				  }
				}`,
			expectedError:       false,
			expectedFork:        bytes.B4{0x4, 0x0, 0x0, 0x0},
			expectedDepositsLen: 0,
		},
		{
			name: "Valid JSON with non-empty deposits",
			jsonInput: `{
				  "fork_version": "0x04000000",
				  "deposits": [{"key": "value"}],
				  "execution_payload_header": {
					"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"feeRecipient": "0x0000000000000000000000000000000000000000",
					"stateRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"receiptsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
					"prevRandao": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"blockNumber": "0x0",
					"gasLimit": "0x0",
					"gasUsed": "0x0",
					"timestamp": "0x0",
					"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"baseFeePerGas": "0x0",
					"blockHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"withdrawalsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"blobGasUsed": "0x0",
					"excessBlobGas": "0x0"
				  }
				}`,
			expectedError:       false,
			expectedFork:        bytes.B4{0x4, 0x0, 0x0, 0x0},
			expectedDepositsLen: 1,
		},
		{
			name: "Invalid JSON input",
			jsonInput: `{
				"fork_version": 12345,
				"deposits": [],
				"execution_payload_header": {
				}
			}`,
			expectedError: true,
		},
		{
			name: "Missing fields in JSON input",
			jsonInput: `{
				  "fork_version": "0x04000000",
				  "deposits": [{"key": "value"}],
				  "execution_payload_header": {
					"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"feeRecipient": "0x0000000000000000000000000000000000000000",
					"stateRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"receiptsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
					"prevRandao": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"gasLimit": "0x1c9c380",
					"gasUsed": "0x0",
					"timestamp": "0x0",
					"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"baseFeePerGas": "0x3b9aca",
					"blockHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"withdrawalsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
					"blobGasUsed": "0x0",
					"excessBlobGas": "0x0"
				  }
				}`,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				g := types.DefaultGenesis(v)
				err := g.UnmarshalJSON([]byte(tc.jsonInput))
				if tc.expectedError {
					require.Error(t, err, "Expected error but got none")
				} else {
					require.NoError(t, err, "Unexpected error")
					require.Equal(t, tc.expectedFork, g.ForkVersion, "Unexpected ForkVersion")
					require.Len(t, g.Deposits, tc.expectedDepositsLen, "Unexpected number of deposits")
					require.NotNil(t, g.ExecutionPayloadHeader, "Expected ExecutionPayloadHeader to be non-nil")
				}
			})
		})
	}
}
