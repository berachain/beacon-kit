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

package cometbft_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	cometbftapi "github.com/berachain/beacon-kit/node-api/handlers/cometbft"
	"github.com/berachain/beacon-kit/node-api/handlers/cometbft/mocks"
	"github.com/berachain/beacon-kit/node-api/middleware"
	cmtversion "github.com/cometbft/cometbft/api/cometbft/version/v1"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGetBlock(t *testing.T) {
	t.Parallel()

	testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)
	testBlock := &cmttypes.Block{
		Header: cmttypes.Header{
			Version: cmtversion.Consensus{Block: 11, App: 0},
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

	testCases := []struct {
		name                string
		height              string
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name:   "success",
			height: "100",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GetCometBFTBlock(int64(100)).Return(testBlock).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()
				require.NoError(t, err)
				require.NotNil(t, res)

				respJSON, err := json.Marshal(res)
				require.NoError(t, err)

				var response cometbftapi.Response
				err = json.Unmarshal(respJSON, &response)
				require.NoError(t, err)
				require.NotNil(t, response.Data)

				dataJSON, err := json.Marshal(response.Data)
				require.NoError(t, err)

				var blockData cmttypes.Block
				err = json.Unmarshal(dataJSON, &blockData)
				require.NoError(t, err)

				require.Equal(t, "test-chain", blockData.Header.ChainID)
				require.Equal(t, int64(100), blockData.Header.Height)
				require.Len(t, blockData.Data.Txs, 2)
			},
		},
		{
			name:   "not found",
			height: "999999",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GetCometBFTBlock(int64(999999)).Return(nil).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "block not found")
				require.Nil(t, res)
			},
		},
		{
			name:                "invalid height - non-numeric",
			height:              "abc",
			setMockExpectations: func(_ *mocks.Backend) {},
			check: func(t *testing.T, res any, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to bind request")
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			backend := mocks.NewBackend(t)
			h := cometbftapi.NewHandler(backend, noop.NewLogger[log.Logger]())
			e := echo.New()
			e.Validator = &middleware.CustomValidator{
				Validator: middleware.ConstructValidator(),
			}

			tc.setMockExpectations(backend)

			req := httptest.NewRequest(http.MethodGet, "/cometbft/v1/block/"+tc.height, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("height")
			c.SetParamValues(tc.height)

			result, err := h.GetBlock(c)

			tc.check(t, result, err)
		})
	}
}

func TestGetSignedHeader(t *testing.T) {
	t.Parallel()

	testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)
	testSignedHeader := &cmttypes.SignedHeader{
		Header: &cmttypes.Header{
			Version:         cmtversion.Consensus{Block: 11, App: 0},
			ChainID:         "test-chain",
			Height:          100,
			Time:            testTime,
			ProposerAddress: []byte("proposer-address"),
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

	testCases := []struct {
		name                string
		height              string
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name:   "success",
			height: "100",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GetCometBFTSignedHeader(int64(100)).Return(testSignedHeader).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()
				require.NoError(t, err)
				require.NotNil(t, res)

				respJSON, err := json.Marshal(res)
				require.NoError(t, err)

				var response cometbftapi.Response
				err = json.Unmarshal(respJSON, &response)
				require.NoError(t, err)

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
			},
		},
		{
			name:   "not found",
			height: "999999",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GetCometBFTSignedHeader(int64(999999)).Return(nil).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "signed header not found")
				require.Nil(t, res)
			},
		},
		{
			name:                "invalid height - non-numeric",
			height:              "abc",
			setMockExpectations: func(_ *mocks.Backend) {},
			check: func(t *testing.T, res any, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to bind request")
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			backend := mocks.NewBackend(t)
			h := cometbftapi.NewHandler(backend, noop.NewLogger[log.Logger]())
			e := echo.New()
			e.Validator = &middleware.CustomValidator{
				Validator: middleware.ConstructValidator(),
			}

			tc.setMockExpectations(backend)

			req := httptest.NewRequest(http.MethodGet, "/cometbft/v1/signed_header/"+tc.height, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("height")
			c.SetParamValues(tc.height)

			result, err := h.GetSignedHeader(c)

			tc.check(t, result, err)
		})
	}
}

func TestGetBlock_FieldConversion(t *testing.T) {
	t.Parallel()

	testTime := time.Date(2024, 1, 26, 12, 0, 0, 123456789, time.UTC)
	block := &cmttypes.Block{
		Header: cmttypes.Header{
			Version:         cmtversion.Consensus{Block: 11, App: 1},
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

	backend := mocks.NewBackend(t)
	h := cometbftapi.NewHandler(backend, noop.NewLogger[log.Logger]())
	e := echo.New()
	e.Validator = &middleware.CustomValidator{
		Validator: middleware.ConstructValidator(),
	}

	backend.EXPECT().GetCometBFTBlock(int64(12345)).Return(block).Once()

	req := httptest.NewRequest(http.MethodGet, "/cometbft/v1/block/12345", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("height")
	c.SetParamValues("12345")

	result, err := h.GetBlock(c)
	require.NoError(t, err)

	respJSON, err := json.Marshal(result)
	require.NoError(t, err)

	var response cometbftapi.Response
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
}

func TestGetSignedHeader_MultipleSignatures(t *testing.T) {
	t.Parallel()

	testTime := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)
	signedHeader := &cmttypes.SignedHeader{
		Header: &cmttypes.Header{
			Version:         cmtversion.Consensus{Block: 11, App: 0},
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

	backend := mocks.NewBackend(t)
	h := cometbftapi.NewHandler(backend, noop.NewLogger[log.Logger]())
	e := echo.New()
	e.Validator = &middleware.CustomValidator{
		Validator: middleware.ConstructValidator(),
	}

	backend.EXPECT().GetCometBFTSignedHeader(int64(100)).Return(signedHeader).Once()

	req := httptest.NewRequest(http.MethodGet, "/cometbft/v1/signed_header/100", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("height")
	c.SetParamValues("100")

	result, err := h.GetSignedHeader(c)
	require.NoError(t, err)

	respJSON, err := json.Marshal(result)
	require.NoError(t, err)

	var response cometbftapi.Response
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

	require.Equal(t, cmttypes.BlockIDFlagCommit, signedHeaderData.Commit.Signatures[0].BlockIDFlag)
	require.Equal(t, cmttypes.BlockIDFlagAbsent, signedHeaderData.Commit.Signatures[1].BlockIDFlag)
	require.Equal(t, cmttypes.BlockIDFlagNil, signedHeaderData.Commit.Signatures[2].BlockIDFlag)
}
