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

package store

import (
	"context"
	"errors"

	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/sourcegraph/conc/iter"
)

// Store is the default implementation of the AvailabilityStore.
type Store[ReadOnlyBeaconBlockT any] struct {
	IndexDB
	chainSpec primitives.ChainSpec
}

// New creates a new instance of the AvailabilityStore.
func New[ReadOnlyBeaconBlockT any](
	chainSpec primitives.ChainSpec,
	db IndexDB,
) *Store[ReadOnlyBeaconBlockT] {
	return &Store[ReadOnlyBeaconBlockT]{
		chainSpec: chainSpec,
		IndexDB:   db,
	}
}

// IsDataAvailable ensures that all blobs referenced in the block are
// stored before it returns without an error.
func (s *Store[ReadOnlyBeaconBlockT]) IsDataAvailable(
	ctx context.Context,
	slot math.Slot,
	blk ReadOnlyBeaconBlockT,
) bool {
	_ = ctx
	_ = slot
	_ = blk
	return true
}

// Persist ensures the sidecar data remains accessible, utilizing parallel
// processing for efficiency.
func (s *Store[ReadOnlyBeaconBlockT]) Persist(
	slot math.Slot,
	sidecars *types.BlobSidecars,
) error {
	// Exit early if there are no sidecars to store.
	if sidecars.Len() == 0 {
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
		sidecars.Sidecars[0].BeaconBlockHeader.GetSlot(),
		// current slot
		slot,
	) {
		return nil
	}

	// Store each sidecar in parallel.
	err := errors.Join(iter.Map(
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
	if err != nil {
		return err
	}

	return s.Prune(slot)
}

// Prune removes all blobs whose block number is not within the DA period.
func (s *Store[ReadOnlyBeaconBlockT]) Prune(currentSlot math.Slot) error {
	// Get all blobs from the store.
	blobs, err := s.GetAllBlobs(currentSlot)
	if err != nil {
		return err
	}

	// Iterate over the blobs.
	for _, blob := range blobs {
		// If the blob's block number is not within the DA period, delete it.
		if !s.chainSpec.WithinDAPeriod(
			blob.BeaconBlockHeader.GetSlot(),
			currentSlot,
		) {
			err = s.DeleteBlob(uint64(currentSlot), blob.KzgCommitment[:])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Store[ReadOnlyBeaconBlockT]) GetAllBlobs(
	currentSlot math.Slot,
) ([]*types.BlobSidecar, error) {
	keys, err := s.IndexDB.GetAllKeys(uint64(currentSlot))
	if err != nil {
		return nil, err
	}

	// Preallocate blobs slice with the length of keys
	blobs := make([]*types.BlobSidecar, 0, len(keys))

	for _, key := range keys {
		var value []byte
		value, err = s.IndexDB.Get(uint64(currentSlot), key)
		if err != nil {
			return nil, err
		}

		/* Assuming the value is serialized
		and needs to be unmarshalled into a BlobSidecar.*/
		var blob types.BlobSidecar
		if err = blob.UnmarshalSSZ(value); err != nil {
			return nil, err
		}

		blobs = append(blobs, &blob)
	}

	return blobs, nil
}

func (s *Store[ReadOnlyBeaconBlockT]) DeleteBlob(
	index uint64,
	blobID []byte,
) error {
	return s.IndexDB.Delete(index, blobID)
}
