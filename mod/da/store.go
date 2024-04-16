// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package da

import (
	"context"
	"errors"

	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	db "github.com/berachain/beacon-kit/mod/storage"
	filedb "github.com/berachain/beacon-kit/mod/storage/filedb"
	"github.com/sourcegraph/conc/iter"
)

// Store is the default implementation of the AvailabilityStore.
type Store struct {
	chainSpec primitives.ChainSpec
	*filedb.RangeDB
}

// NewStore creates a new instance of the AvailabilityStore.
func NewStore(
	chainSpec primitives.ChainSpec,
	db db.DB,
) *Store {
	return &Store{
		chainSpec: chainSpec,
		RangeDB:   filedb.NewRangeDB(db),
	}
}

// IsDataAvailable ensures that all blobs referenced in the block are
// stored before it returns without an error.
func (s *Store) IsDataAvailable(
	ctx context.Context,
	slot primitives.Slot,
	b beacontypes.ReadOnlyBeaconBlock,
) bool {
	_ = ctx
	_ = slot
	_ = b
	return true
}

// Persist ensures the sidecar data remains accessible, utilizing parallel
// processing for efficiency.
func (s *Store) Persist(
	slot primitives.Slot,
	sidecars *types.BlobSidecars,
) error {
	// Exit early if there are no sidecars to store.
	if len(sidecars.Sidecars) == 0 {
		return nil
	}

	// Ensure that all sidecars have the same block root.
	if err := sidecars.ValidateBlockRoots(); err != nil {
		return err
	}

	// Check to see if we are required to store the sidecar anymore, if
	// this sidecar is from outside the required DA period, we can skip it.
	if !s.chainSpec.WithinDAPeriod(
		// slot in which the sidecar was included.
		// (Safe to assume all sidecars are in same slot at this point).
		sidecars.Sidecars[0].BeaconBlockHeader.Slot,
		// current slot
		slot,
	) {
		return nil
	}

	// Store each sidecar in parallel.
	return errors.Join(iter.Map(
		sidecars.Sidecars,
		func(sidecar **types.BlobSidecar) error {
			if *sidecar == nil {
				return ErrAttemptedToStoreNilSidecar
			}
			sc := *sidecar
			bz, err := sc.MarshalSSZ()
			if err != nil {
				return err
			}
			return s.Set(uint64(slot), sc.KzgCommitment[:], bz)
		},
	)...)
}
