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
	"context"
	"testing"

	"cosmossdk.io/log"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/blobreactor"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage/block"
	"github.com/berachain/beacon-kit/storage/deposit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockBlobProcessorForExecutor struct{ mock.Mock }

func (m *mockBlobProcessorForExecutor) VerifySidecars(
	ctx context.Context,
	sidecars datypes.BlobSidecars,
	blkHeader *ctypes.BeaconBlockHeader,
	kzgCommitments eip4844.KZGCommitments[common.ExecutionHash],
) error {
	return m.Called(ctx, sidecars, blkHeader, kzgCommitments).Error(0)
}

func (m *mockBlobProcessorForExecutor) ProcessSidecars(avs *dastore.Store, sidecars datypes.BlobSidecars) error {
	return m.Called(avs, sidecars).Error(0)
}

type mockBlobRequesterForExecutor struct{ mock.Mock }

func (m *mockBlobRequesterForExecutor) RequestBlobs(
	ctx context.Context,
	slot math.Slot,
	verifier func(datypes.BlobSidecars) error,
) ([]*datypes.BlobSidecar, error) {
	args := m.Called(ctx, slot, verifier)
	if args.Get(0) != nil {
		sidecars, ok := args.Get(0).([]*datypes.BlobSidecar)
		if !ok {
			return nil, args.Error(1)
		}
		if verifier != nil {
			if err := verifier(sidecars); err != nil {
				return nil, err // Verifier rejected
			}
		}
		return sidecars, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockBlobRequesterForExecutor) SetHeadSlot(_ math.Slot) {}

type mockStorageBackendForExecutor struct {
	availStore *dastore.Store
}

func (m *mockStorageBackendForExecutor) AvailabilityStore() *dastore.Store { return m.availStore }
func (m *mockStorageBackendForExecutor) StateFromContext(_ context.Context) *statedb.StateDB {
	return nil
}
func (m *mockStorageBackendForExecutor) DepositStore() deposit.StoreManager { return nil }
func (m *mockStorageBackendForExecutor) BlockStore() *block.KVStore[*ctypes.BeaconBlock] {
	return nil
}

// Test when peer sends invalid blobs (verification should reject them)
func TestBlobFetchExecutor_ByzantineBlobs_Rejected(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	mockProcessor := &mockBlobProcessorForExecutor{}
	mockRequester := &mockBlobRequesterForExecutor{}
	mockStorage := &mockStorageBackendForExecutor{availStore: &dastore.Store{}}

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
	mockProcessor := &mockBlobProcessorForExecutor{}
	mockRequester := &mockBlobRequesterForExecutor{}
	mockStorage := &mockStorageBackendForExecutor{availStore: &dastore.Store{}}

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
	mockProcessor.On("ProcessSidecars", mockStorage.availStore, mock.Anything).Return(nil)

	err := executor.FetchBlobsAndVerify(ctx, request)
	require.NoError(t, err)
	require.True(t, verifierCalled, "Verifier must be called for Byzantine protection")
}

// Test when all peers fail - no valid blobs available
func TestBlobFetchExecutor_AllPeersFailed(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	mockProcessor := &mockBlobProcessorForExecutor{}
	mockRequester := &mockBlobRequesterForExecutor{}
	mockStorage := &mockStorageBackendForExecutor{availStore: &dastore.Store{}}

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
