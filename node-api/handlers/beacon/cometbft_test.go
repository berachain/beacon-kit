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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package beacon_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCometBFTBlock(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	// Create a test block
	testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)
	testBlock := &cmttypes.Block{
		Header: cmttypes.Header{
			Version: cmttypes.Consensus{Block: 11, App: 0},
			ChainID: "test-chain",
			Height:  100,
			Time:    testTime,
			LastBlockID: cmttypes.BlockID{
				Hash: []byte("prev-block-hash"),
				PartSetHeader: cmttypes.PartSetHeader{
					Total: 1,
					Hash:  []byte("part-set-hash"),
				},
			},
			LastCommitHash:     []byte("last-commit-hash"),
			DataHash:           []byte("data-hash"),
			ValidatorsHash:     []byte("validators-hash"),
			NextValidatorsHash: []byte("next-validators-hash"),
			ConsensusHash:      []byte("consensus-hash"),
			AppHash:            []byte("app-hash"),
			LastResultsHash:    []byte("last-results-hash"),
			EvidenceHash:       []byte("evidence-hash"),
			ProposerAddress:    []byte("proposer-address"),
		},
		Data: cmttypes.Data{
			Txs: []cmttypes.Tx{
				[]byte("transaction-1"),
				[]byte("transaction-2"),
			},
		},
		Evidence: cmttypes.EvidenceData{
			Evidence: []cmttypes.Evidence{},
		},
		LastCommit: &cmttypes.Commit{
			Height: 99,
			Round:  0,
			BlockID: cmttypes.BlockID{
				Hash: []byte("prev-block-hash"),
				PartSetHeader: cmttypes.PartSetHeader{
					Total: 1,
					Hash:  []byte("part-set-hash"),
				},
			},
			Signatures: []cmttypes.CommitSig{
				{
					BlockIDFlag:      cmttypes.BlockIDFlagCommit,
					ValidatorAddress: []byte("validator-1"),
					Timestamp:        testTime,
					Signature:        []byte("signature-1"),
				},
			},
		},
	}

	mockBackend.EXPECT().GetCometBFTBlock(int64(100)).Return(testBlock).Once()

	req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/block/100", nil)
	req.SetPathValue("height", "100")
	rec := httptest.NewRecorder()

	result, err := handler.GetCometBFTBlock(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Marshal and unmarshal to check response structure
	respJSON, err := json.Marshal(result)
	require.NoError(t, err)

	var response beacontypes.GenericResponse
	err = json.Unmarshal(respJSON, &response)
	require.NoError(t, err)

	require.True(t, response.Finalized)
	require.False(t, response.ExecutionOptimistic)
	require.NotNil(t, response.Data)

	// Verify response contains block data
	dataJSON, err := json.Marshal(response.Data)
	require.NoError(t, err)

	var blockData cmttypes.Block
	err = json.Unmarshal(dataJSON, &blockData)
	require.NoError(t, err)

	require.Equal(t, "test-chain", blockData.Header.ChainID)
	require.Equal(t, int64(100), blockData.Header.Height)
	require.Len(t, blockData.Data.Txs, 2)

	_ = rec
}

func TestGetCometBFTBlock_NotFound(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	mockBackend.EXPECT().GetCometBFTBlock(int64(999999)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/block/999999", nil)
	req.SetPathValue("height", "999999")

	_, err = handler.GetCometBFTBlock(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "block not found")
}

func TestGetCometBFTSignedHeader(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)
	testSignedHeader := &cmttypes.SignedHeader{
		Header: &cmttypes.Header{
			Version:         cmttypes.Consensus{Block: 11, App: 0},
			ChainID:         "test-chain",
			Height:          100,
			Time:            testTime,
			ProposerAddress: []byte("proposer-address"),
		},
		Commit: &cmttypes.Commit{
			Height: 100,
			Round:  0,
			BlockID: cmttypes.BlockID{
				Hash: []byte("block-hash"),
				PartSetHeader: cmttypes.PartSetHeader{
					Total: 1,
					Hash:  []byte("part-set-hash"),
				},
			},
			Signatures: []cmttypes.CommitSig{
				{
					BlockIDFlag:      cmttypes.BlockIDFlagCommit,
					ValidatorAddress: []byte("validator-1"),
					Timestamp:        testTime,
					Signature:        []byte("signature-1"),
				},
				{
					BlockIDFlag:      cmttypes.BlockIDFlagCommit,
					ValidatorAddress: []byte("validator-2"),
					Timestamp:        testTime.Add(time.Millisecond * 100),
					Signature:        []byte("signature-2"),
				},
			},
		},
	}

	mockBackend.EXPECT().GetCometBFTSignedHeader(int64(100)).Return(testSignedHeader).Once()

	req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/signed_header/100", nil)
	req.SetPathValue("height", "100")
	rec := httptest.NewRecorder()

	result, err := handler.GetCometBFTSignedHeader(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	respJSON, err := json.Marshal(result)
	require.NoError(t, err)

	var response beacontypes.GenericResponse
	err = json.Unmarshal(respJSON, &response)
	require.NoError(t, err)

	require.True(t, response.Finalized)
	require.False(t, response.ExecutionOptimistic)

	dataJSON, err := json.Marshal(response.Data)
	require.NoError(t, err)

	var signedHeaderData cmttypes.SignedHeader
	err = json.Unmarshal(dataJSON, &signedHeaderData)
	require.NoError(t, err)

	require.Equal(t, "test-chain", signedHeaderData.Header.ChainID)
	require.Equal(t, int64(100), signedHeaderData.Header.Height)
	require.Equal(t, int64(100), signedHeaderData.Commit.Height)
	require.Equal(t, int32(0), signedHeaderData.Commit.Round)
	require.Len(t, signedHeaderData.Commit.Signatures, 2)

	_ = rec
}

func TestGetCometBFTSignedHeader_NotFound(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	mockBackend.EXPECT().GetCometBFTSignedHeader(int64(999999)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/signed_header/999999", nil)
	req.SetPathValue("height", "999999")

	_, err = handler.GetCometBFTSignedHeader(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "signed header not found")
}



func TestCometBFTConversionFunctions(t *testing.T) {
	t.Run("Block conversion includes all fields", func(t *testing.T) {
		testTime := time.Date(2024, 1, 26, 12, 0, 0, 123456789, time.UTC)
		block := &cmttypes.Block{
			Header: cmttypes.Header{
				Version:         cmttypes.Consensus{Block: 11, App: 1},
				ChainID:         "test-chain-123",
				Height:          12345,
				Time:            testTime,
				ProposerAddress: []byte("proposer"),
			},
			Data: cmttypes.Data{
				Txs: []cmttypes.Tx{[]byte("tx1"), []byte("tx2"), []byte("tx3")},
			},
			Evidence: cmttypes.EvidenceData{
				Evidence: []cmttypes.Evidence{},
			},
		}

		mockBackend := mocks.NewBackend(t)
		cs, err := spec.DevnetChainSpec()
		require.NoError(t, err)
		handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

		mockBackend.EXPECT().GetCometBFTBlock(int64(12345)).Return(block).Once()

		req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/block/12345", nil)
		req.SetPathValue("height", "12345")

		result, err := handler.GetCometBFTBlock(req)
		require.NoError(t, err)

		respJSON, err := json.Marshal(result)
		require.NoError(t, err)

		var response beacontypes.GenericResponse
		err = json.Unmarshal(respJSON, &response)
		require.NoError(t, err)

		dataJSON, err := json.Marshal(response.Data)
		require.NoError(t, err)

		var blockData cmttypes.Block
		err = json.Unmarshal(dataJSON, &blockData)
		require.NoError(t, err)

		require.Equal(t, "test-chain-123", blockData.Header.ChainID)
		require.Equal(t, int64(12345), blockData.Header.Height)
		require.Equal(t, uint64(11), blockData.Header.Version.Block)
		require.Equal(t, uint64(1), blockData.Header.Version.App)
		require.Len(t, blockData.Data.Txs, 3)
	})

	t.Run("SignedHeader conversion handles multiple signatures", func(t *testing.T) {
		testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)
		signedHeader := &cmttypes.SignedHeader{
			Header: &cmttypes.Header{
				Version:         cmttypes.Consensus{Block: 11, App: 0},
				ChainID:         "test-chain",
				Height:          100,
				Time:            testTime,
				ProposerAddress: []byte("proposer"),
			},
			Commit: &cmttypes.Commit{
				Height: 100,
				Round:  2,
				BlockID: cmttypes.BlockID{
					Hash: []byte("block-hash"),
				},
				Signatures: []cmttypes.CommitSig{
					{
						BlockIDFlag:      cmttypes.BlockIDFlagCommit,
						ValidatorAddress: []byte("val1"),
						Timestamp:        testTime,
						Signature:        []byte("sig1"),
					},
					{
						BlockIDFlag:      cmttypes.BlockIDFlagAbsent,
						ValidatorAddress: []byte("val2"),
						Timestamp:        testTime,
						Signature:        nil,
					},
					{
						BlockIDFlag:      cmttypes.BlockIDFlagNil,
						ValidatorAddress: []byte("val3"),
						Timestamp:        testTime,
						Signature:        []byte("sig3"),
					},
				},
			},
		}

		mockBackend := mocks.NewBackend(t)
		cs, err := spec.DevnetChainSpec()
		require.NoError(t, err)
		handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

		mockBackend.EXPECT().GetCometBFTSignedHeader(int64(100)).Return(signedHeader).Once()

		req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/signed_header/100", nil)
		req.SetPathValue("height", "100")

		result, err := handler.GetCometBFTSignedHeader(req)
		require.NoError(t, err)

		respJSON, err := json.Marshal(result)
		require.NoError(t, err)

		var response beacontypes.GenericResponse
		err = json.Unmarshal(respJSON, &response)
		require.NoError(t, err)

		dataJSON, err := json.Marshal(response.Data)
		require.NoError(t, err)

		var signedHeaderData cmttypes.SignedHeader
		err = json.Unmarshal(dataJSON, &signedHeaderData)
		require.NoError(t, err)

		require.Equal(t, "test-chain", signedHeaderData.Header.ChainID)
		require.Equal(t, int64(100), signedHeaderData.Header.Height)
		require.Equal(t, int32(2), signedHeaderData.Commit.Round)
		require.Len(t, signedHeaderData.Commit.Signatures, 3)

		// Verify different block ID flags are represented
		require.Equal(t, cmttypes.BlockIDFlagCommit, signedHeaderData.Commit.Signatures[0].BlockIDFlag)
		require.Equal(t, cmttypes.BlockIDFlagAbsent, signedHeaderData.Commit.Signatures[1].BlockIDFlag)
		require.Equal(t, cmttypes.BlockIDFlagNil, signedHeaderData.Commit.Signatures[2].BlockIDFlag)
	})
}

func TestCometBFTEndpoints_Integration(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)

	// Create related test data
	commit := &cmttypes.Commit{
		Height: 99,
		Round:  0,
		BlockID: cmttypes.BlockID{
			Hash: []byte("block-99"),
		},
		Signatures: []cmttypes.CommitSig{
			{
				BlockIDFlag:      cmttypes.BlockIDFlagCommit,
				ValidatorAddress: validator.Address,
				Timestamp:        testTime,
				Signature:        []byte("signature"),
			},
		},
	}

	block := &cmttypes.Block{
		Header: cmttypes.Header{
			Version:         cmttypes.Consensus{Block: 11, App: 0},
			ChainID:         "integration-test",
			Height:          100,
			Time:            testTime,
			ProposerAddress: []byte("test-proposer"),
		},
		LastCommit: commit,
	}

	signedHeader := &cmttypes.SignedHeader{
		Header: &block.Header,
		Commit: commit,
	}

	// Setup expectations
	mockBackend.EXPECT().GetCometBFTBlock(int64(100)).Return(block).Once()
	mockBackend.EXPECT().GetCometBFTSignedHeader(int64(100)).Return(signedHeader).Once()

	// Test block endpoint
	blockReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/block/100", nil)
	blockReq.SetPathValue("height", "100")
	blockResult, err := handler.GetCometBFTBlock(blockReq)
	require.NoError(t, err)
	require.NotNil(t, blockResult)

	// Test signed header endpoint
	signedHeaderReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/signed_header/100", nil)
	signedHeaderReq.SetPathValue("height", "100")
	signedHeaderResult, err := handler.GetCometBFTSignedHeader(signedHeaderReq)
	require.NoError(t, err)
	require.NotNil(t, signedHeaderResult)

	// Verify all results are properly formatted
	for _, result := range []any{blockResult, signedHeaderResult} {
		respJSON, err := json.Marshal(result)
		require.NoError(t, err)

		var response beacontypes.GenericResponse
		err = json.Unmarshal(respJSON, &response)
		require.NoError(t, err)

		require.True(t, response.Finalized)
		require.False(t, response.ExecutionOptimistic)
	}
}

func TestCometBFTEndpoints_MockExpectations(t *testing.T) {
	t.Run("All expectations are met", func(t *testing.T) {
		mockBackend := mocks.NewBackend(t)
		cs, err := spec.DevnetChainSpec()
		require.NoError(t, err)
		handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

		testBlock := &cmttypes.Block{
			Header: cmttypes.Header{
				ChainID: "test",
				Height:  1,
			},
		}

		testSignedHeader := &cmttypes.SignedHeader{
			Header: &testBlock.Header,
			Commit: &cmttypes.Commit{
				Height: 1,
			},
		}

		// Set expectations
		mockBackend.EXPECT().GetCometBFTBlock(int64(1)).Return(testBlock).Once()
		mockBackend.EXPECT().GetCometBFTSignedHeader(int64(1)).Return(testSignedHeader).Once()

		// Execute
		blockReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/block/1", nil)
		blockReq.SetPathValue("height", "1")
		_, _ = handler.GetCometBFTBlock(blockReq)

		signedHeaderReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/signed_header/1", nil)
		signedHeaderReq.SetPathValue("height", "1")
		_, _ = handler.GetCometBFTSignedHeader(signedHeaderReq)

		// Verify all expectations met
		mock.AssertExpectationsForObjects(t, mockBackend)
	})
}