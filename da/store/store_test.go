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
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/storage/filedb"
	"github.com/stretchr/testify/require"
)

func setSlot(scs datypes.BlobSidecars, slot math.Slot) {
	for _, sc := range scs {
		hdr := sc.GetSignedBeaconBlockHeader().GetHeader()
		hdr.SetSlot(slot)
	}
}

func TestStore_PersistRace(t *testing.T) {
	t.Parallel()
	// This test case needs to be run with the '-race' flag
	tmpFilePath := t.TempDir()

	logger := log.NewNopLogger()
	chainSpec, err := spec.DevnetChainSpec()
	require.NoError(t, err)

	// Create the DB
	s := store.New(
		filedb.NewRangeDB(
			filedb.NewDB(filedb.WithRootDirectory(tmpFilePath),
				filedb.WithFileExtension("ssz"),
				filedb.WithDirectoryPermissions(0700),
				filedb.WithLogger(logger),
			),
		),
		logger.With("service", "da-store"),
		chainSpec,
	)

	// This many blobs is not currently possible, but it doesn't hurt eh
	sc := make([]*datypes.BlobSidecar, 20)
	for i := range sc {
		sc[i] = &datypes.BlobSidecar{
			Index: uint64(i),
			SignedBeaconBlockHeader: &types.SignedBeaconBlockHeader{
				Header: &types.BeaconBlockHeader{},
			},
			InclusionProof: make([]common.Root, types.KZGInclusionProofDepth),
		}
	}
	var sidecars datypes.BlobSidecars = sc

	// Multiple writes to DB
	setSlot(sidecars, 0)
	err = s.Persist(sidecars)
	require.NoError(t, err)
	setSlot(sidecars, 1)
	err = s.Persist(sidecars)
	require.NoError(t, err)

	// Pruning here primes the race condition for db.firstNonNilIndex
	err = s.Prune(0, 1)
	require.NoError(t, err)
	setSlot(sidecars, 0)
	err = s.Persist(sidecars)
	require.NoError(t, err)
}
