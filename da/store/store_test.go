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

package store_test

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/storage/filedb"
	"github.com/stretchr/testify/require"
)

func setSlot(scs datypes.BlobSidecars, slot math.Slot) {
	for _, sc := range scs {
		hdr := sc.GetBeaconBlockHeader()
		hdr.SetSlot(slot)
	}
}

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	logger := log.NewNopLogger()
	return store.New(
		filedb.NewRangeDB(
			filedb.NewDB(filedb.WithRootDirectory(t.TempDir()),
				filedb.WithFileExtension("ssz"),
				filedb.WithDirectoryPermissions(0700),
				filedb.WithLogger(logger),
			),
		),
		logger.With("service", "da-store"),
	)
}

func newTestSidecar(slot math.Slot, index uint64, commitment eip4844.KZGCommitment) *datypes.BlobSidecar {
	return &datypes.BlobSidecar{
		Index:         index,
		KzgCommitment: commitment,
		SignedBeaconBlockHeader: &types.SignedBeaconBlockHeader{
			Header: &types.BeaconBlockHeader{Slot: slot},
		},
		InclusionProof: make([]common.Root, types.KZGInclusionProofDepth),
	}
}

func TestStore_PersistRace(t *testing.T) {
	t.Parallel()
	// This test case needs to be run with the '-race' flag
	s := newTestStore(t)

	// This many blobs is not currently possible, but it doesn't hurt eh
	sc := make([]*datypes.BlobSidecar, 20)
	for i := range sc {
		sc[i] = newTestSidecar(0, uint64(i), eip4844.KZGCommitment{})
	}
	var sidecars datypes.BlobSidecars = sc

	// Multiple writes to DB
	setSlot(sidecars, 0)
	err := s.Persist(sidecars)
	require.NoError(t, err)
	setSlot(sidecars, 1)
	err = s.Persist(sidecars)
	require.NoError(t, err)

	// Pruning here primes the race condition for db.firstNonNilIndex
	err = s.Prune(0, 1)
	require.NoError(t, err)

	// Persisting slot-0 data after pruning past it must be rejected, not silently written under the pruning
	// lower bound. Clamping would corrupt slot 1's data with slot 0's sidecars.
	setSlot(sidecars, 0)
	err = s.Persist(sidecars)
	require.ErrorIs(t, err, filedb.ErrIndexPruned)
}

// Duplicate KZG commitments within one block must not overwrite each other: the storage key includes the blob
// index, so every sidecar stays retrievable and availability accounts for each index individually.
func TestStore_DuplicateCommitments(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)

	var (
		slot       = math.Slot(7)
		commitment = eip4844.KZGCommitment{0xaa}
	)
	require.NoError(t, s.Persist(datypes.BlobSidecars{
		newTestSidecar(slot, 0, commitment),
		newTestSidecar(slot, 1, commitment),
	}))

	got, err := s.GetBlobSidecars(slot)
	require.NoError(t, err)
	require.Len(t, got, 2, "the second sidecar must not overwrite the first")

	body := &types.BeaconBlockBody{}
	body.SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash]{commitment, commitment})
	require.True(t, s.IsDataAvailable(t.Context(), slot, body))
}

// Persist replaces a slot's existing sidecars rather than accumulating them, so re-persisting a block (e.g. on
// replay across the storage-key change) cannot leave duplicate entries for the same index.
func TestStore_PersistReplacesSlot(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)

	slot := math.Slot(9)
	set := datypes.BlobSidecars{
		newTestSidecar(slot, 0, eip4844.KZGCommitment{0}),
		newTestSidecar(slot, 1, eip4844.KZGCommitment{1}),
	}
	require.NoError(t, s.Persist(set))
	require.NoError(t, s.Persist(set)) // re-persist the same block

	got, err := s.GetBlobSidecars(slot)
	require.NoError(t, err)
	require.Len(t, got, 2, "re-persisting must replace, not duplicate")
}

// Persist deletes the first sidecar's slot before writing, so a mixed-slot slice must be rejected outright
// rather than deleting one slot's data and writing under another.
func TestStore_PersistRejectsMixedSlots(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)

	commitment := eip4844.KZGCommitment{0xaa}
	require.NoError(t, s.Persist(datypes.BlobSidecars{
		newTestSidecar(math.Slot(8), 0, commitment),
	}))

	err := s.Persist(datypes.BlobSidecars{
		newTestSidecar(math.Slot(7), 0, commitment),
		newTestSidecar(math.Slot(8), 1, commitment),
	})
	require.ErrorIs(t, err, store.ErrMixedSlotSidecars)

	// Slot 8's previously persisted sidecar must be untouched by the rejected call.
	got, err := s.GetBlobSidecars(math.Slot(8))
	require.NoError(t, err)
	require.Len(t, got, 1)
}
