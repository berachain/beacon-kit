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

func TestGetCometBFTCommit(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)
	testCommit := &cmttypes.Commit{
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
	}

	mockBackend.EXPECT().GetCometBFTCommit(int64(100)).Return(testCommit).Once()

	req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/commit/100", nil)
	req.SetPathValue("height", "100")
	rec := httptest.NewRecorder()

	result, err := handler.GetCometBFTCommit(req)
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

	var commitData cmttypes.Commit
	err = json.Unmarshal(dataJSON, &commitData)
	require.NoError(t, err)

	require.Equal(t, int64(100), commitData.Height)
	require.Equal(t, int32(0), commitData.Round)
	require.Len(t, commitData.Signatures, 2)

	_ = rec
}

func TestGetCometBFTCommit_NotFound(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	mockBackend.EXPECT().GetCometBFTCommit(int64(999999)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/commit/999999", nil)
	req.SetPathValue("height", "999999")

	_, err = handler.GetCometBFTCommit(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "commit not found")
}

func TestGetCometBFTValidators(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	validator1 := &cmttypes.Validator{
		Address:          []byte("validator-address-1"),
		PubKey:           cmttypes.NewMockPV().PrivKey.PubKey(),
		VotingPower:      1000000,
		ProposerPriority: -500000,
	}

	validator2 := &cmttypes.Validator{
		Address:          []byte("validator-address-2"),
		PubKey:           cmttypes.NewMockPV().PrivKey.PubKey(),
		VotingPower:      2000000,
		ProposerPriority: 500000,
	}

	testValidatorSet := cmttypes.NewValidatorSet([]*cmttypes.Validator{
		validator1,
		validator2,
	})

	mockBackend.EXPECT().GetCometBFTValidators(int64(100)).Return(testValidatorSet, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/validators/100", nil)
	req.SetPathValue("height", "100")
	rec := httptest.NewRecorder()

	result, err := handler.GetCometBFTValidators(req)
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

	var validatorsData cmttypes.ValidatorSet
	err = json.Unmarshal(dataJSON, &validatorsData)
	require.NoError(t, err)

	require.Len(t, validatorsData.Validators, 2)
	require.NotNil(t, validatorsData.Proposer)
	require.Equal(t, int64(3000000), validatorsData.TotalVotingPower())

	_ = rec
}

func TestGetCometBFTValidators_Error(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	expectedErr := errors.New("state store error")
	mockBackend.EXPECT().GetCometBFTValidators(int64(100)).Return(nil, expectedErr).Once()

	req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/validators/100", nil)
	req.SetPathValue("height", "100")

	_, err = handler.GetCometBFTValidators(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get validators")
}

func TestGetCometBFTValidators_NotFound(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	mockBackend.EXPECT().GetCometBFTValidators(int64(999999)).Return(nil, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/validators/999999", nil)
	req.SetPathValue("height", "999999")

	_, err = handler.GetCometBFTValidators(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "validators not found")
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

	t.Run("Commit conversion handles multiple signatures", func(t *testing.T) {
		testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)
		commit := &cmttypes.Commit{
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
		}

		mockBackend := mocks.NewBackend(t)
		cs, err := spec.DevnetChainSpec()
		require.NoError(t, err)
		handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

		mockBackend.EXPECT().GetCometBFTCommit(int64(100)).Return(commit).Once()

		req := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/commit/100", nil)
		req.SetPathValue("height", "100")

		result, err := handler.GetCometBFTCommit(req)
		require.NoError(t, err)

		respJSON, err := json.Marshal(result)
		require.NoError(t, err)

		var response beacontypes.GenericResponse
		err = json.Unmarshal(respJSON, &response)
		require.NoError(t, err)

		dataJSON, err := json.Marshal(response.Data)
		require.NoError(t, err)

		var commitData cmttypes.Commit
		err = json.Unmarshal(dataJSON, &commitData)
		require.NoError(t, err)

		require.Equal(t, int32(2), commitData.Round)
		require.Len(t, commitData.Signatures, 3)

		// Verify different block ID flags are represented
		require.Equal(t, cmttypes.BlockIDFlagCommit, commitData.Signatures[0].BlockIDFlag)
		require.Equal(t, cmttypes.BlockIDFlagAbsent, commitData.Signatures[1].BlockIDFlag)
		require.Equal(t, cmttypes.BlockIDFlagNil, commitData.Signatures[2].BlockIDFlag)
	})
}

func TestCometBFTEndpoints_Integration(t *testing.T) {
	mockBackend := mocks.NewBackend(t)
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	handler := beacon.NewHandler(mockBackend, cs, noop.NewLogger())

	testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)

	// Create related test data
	validator := &cmttypes.Validator{
		Address:          []byte("test-validator"),
		PubKey:           cmttypes.NewMockPV().PrivKey.PubKey(),
		VotingPower:      1000000,
		ProposerPriority: 0,
	}

	validatorSet := cmttypes.NewValidatorSet([]*cmttypes.Validator{validator})

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
			ProposerAddress: validator.Address,
			ValidatorsHash:  validatorSet.Hash(),
		},
		LastCommit: commit,
	}

	// Setup expectations
	mockBackend.EXPECT().GetCometBFTBlock(int64(100)).Return(block).Once()
	mockBackend.EXPECT().GetCometBFTCommit(int64(100)).Return(commit).Once()
	mockBackend.EXPECT().GetCometBFTValidators(int64(100)).Return(validatorSet, nil).Once()

	// Test block endpoint
	blockReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/block/100", nil)
	blockReq.SetPathValue("height", "100")
	blockResult, err := handler.GetCometBFTBlock(blockReq)
	require.NoError(t, err)
	require.NotNil(t, blockResult)

	// Test commit endpoint
	commitReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/commit/100", nil)
	commitReq.SetPathValue("height", "100")
	commitResult, err := handler.GetCometBFTCommit(commitReq)
	require.NoError(t, err)
	require.NotNil(t, commitResult)

	// Test validators endpoint
	valsReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/validators/100", nil)
	valsReq.SetPathValue("height", "100")
	valsResult, err := handler.GetCometBFTValidators(valsReq)
	require.NoError(t, err)
	require.NotNil(t, valsResult)

	// Verify all results are properly formatted
	for _, result := range []any{blockResult, commitResult, valsResult} {
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

		testCommit := &cmttypes.Commit{
			Height: 1,
		}

		testValidators := cmttypes.NewValidatorSet([]*cmttypes.Validator{
			{
				Address:     []byte("validator"),
				PubKey:      cmttypes.NewMockPV().PrivKey.PubKey(),
				VotingPower: 1000,
			},
		})

		// Set expectations
		mockBackend.EXPECT().GetCometBFTBlock(int64(1)).Return(testBlock).Once()
		mockBackend.EXPECT().GetCometBFTCommit(int64(1)).Return(testCommit).Once()
		mockBackend.EXPECT().GetCometBFTValidators(int64(1)).Return(testValidators, nil).Once()

		// Execute
		blockReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/block/1", nil)
		blockReq.SetPathValue("height", "1")
		_, _ = handler.GetCometBFTBlock(blockReq)

		commitReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/commit/1", nil)
		commitReq.SetPathValue("height", "1")
		_, _ = handler.GetCometBFTCommit(commitReq)

		valsReq := httptest.NewRequest(http.MethodGet, "/eth/v1/beacon/cometbft/validators/1", nil)
		valsReq.SetPathValue("height", "1")
		_, _ = handler.GetCometBFTValidators(valsReq)

		// Verify all expectations met
		mock.AssertExpectationsForObjects(t, mockBackend)
	})
}