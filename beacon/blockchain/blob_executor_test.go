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
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//nolint:testpackage // Testing internal components
package blockchain

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/beacon/blockchain/testhelpers"
	"github.com/berachain/beacon-kit/da/blobreactor"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test when peer sends invalid blobs (verification should reject them)
func TestBlobFetchExecutor_ByzantineBlobs_Rejected(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	mockProcessor := &testhelpers.SimpleBlobProcessor{}
	mockRequester := &testhelpers.SimpleBlobRequester{}
	mockStorage := testhelpers.NewSimpleStorageBackend(&dastore.Store{})

	executor := &blobFetchExecutor{
		blobProcessor:  mockProcessor,
		blobRequester:  mockRequester,
		storageBackend: mockStorage,
		logger:         log.NewNopLogger(),
	}

	request := createTestBlobRequest(math.Slot(100), 2)
	invalidBlobs := []*datypes.BlobSidecar{{Index: 0}, {Index: 1}}

	// Byzantine peer returns invalid blobs - KZG proof verification fails
	verifyErr := errors.New("KZG proof verification failed")
	mockProcessor.On("VerifySidecars", ctx, mock.Anything, request.Header, request.Commitments).Return(verifyErr)
	mockRequester.On("RequestBlobs", ctx, math.Slot(100), mock.Anything).Return(invalidBlobs, verifyErr)

	err := executor.FetchBlobsAndVerify(ctx, request)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to request valid blobs")

	// ProcessSidecars must NOT be called with invalid blobs
	mockProcessor.AssertNotCalled(t, "ProcessSidecars")
}

// Test Verifier function is called (Byzantine protection mechanism)
func TestBlobFetchExecutor_VerifierCalled(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	mockProcessor := &testhelpers.SimpleBlobProcessor{}
	mockRequester := &testhelpers.SimpleBlobRequester{}
	mockStorage := testhelpers.NewSimpleStorageBackend(&dastore.Store{})

	executor := &blobFetchExecutor{
		blobProcessor:  mockProcessor,
		blobRequester:  mockRequester,
		storageBackend: mockStorage,
		logger:         log.NewNopLogger(),
	}

	request := createTestBlobRequest(math.Slot(100), 1)
	validBlobs := []*datypes.BlobSidecar{{Index: 0}}

	verifierCalled := false
	mockProcessor.On("VerifySidecars", ctx, mock.Anything, request.Header, request.Commitments).
		Run(func(_ mock.Arguments) { verifierCalled = true }).
		Return(nil)
	mockRequester.On("RequestBlobs", ctx, math.Slot(100), mock.Anything).Return(validBlobs, nil)
	mockProcessor.On("ProcessSidecars", mockStorage.AvailabilityStore(), mock.Anything).Return(nil)

	err := executor.FetchBlobsAndVerify(ctx, request)
	require.NoError(t, err)
	require.True(t, verifierCalled, "Verifier must be called for Byzantine protection")
}

// Test when all peers fail - no valid blobs available
func TestBlobFetchExecutor_AllPeersFailed(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	mockProcessor := &testhelpers.SimpleBlobProcessor{}
	mockRequester := &testhelpers.SimpleBlobRequester{}
	mockStorage := testhelpers.NewSimpleStorageBackend(&dastore.Store{})

	executor := &blobFetchExecutor{
		blobProcessor:  mockProcessor,
		blobRequester:  mockRequester,
		storageBackend: mockStorage,
		logger:         log.NewNopLogger(),
	}

	request := createTestBlobRequest(math.Slot(100), 2)

	// All peers failed (timeout, byzantine, offline, etc.)
	mockRequester.On("RequestBlobs", ctx, math.Slot(100), mock.Anything).Return(nil, blobreactor.ErrAllPeersFailed)

	err := executor.FetchBlobsAndVerify(ctx, request)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to request valid blobs")

	// No blobs should be stored
	mockProcessor.AssertNotCalled(t, "ProcessSidecars")
}
