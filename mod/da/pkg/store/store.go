// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package store

import (
	"context"
	"sync"

	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// SLOT_COMMITMENTS_KEY is the key used to store the commitments for a slot
// in the DB. We use this key to avoid conflicts with the slot index.
const SLOT_COMMITMENTS_KEY = "slot_commitments"

// Store is the default implementation of the AvailabilityStore.
type Store[BeaconBlockBodyT BeaconBlockBody] struct {
	// IndexDB is a basic database interface.
	IndexDB
	// logger is used for logging.
	logger log.Logger
	// chainSpec contains the chain specification.
	chainSpec common.ChainSpec
}

// New creates a new instance of the AvailabilityStore.
func New[BeaconBlockT BeaconBlockBody](
	db IndexDB,
	logger log.Logger,
	chainSpec common.ChainSpec,
) *Store[BeaconBlockT] {
	return &Store[BeaconBlockT]{
		IndexDB:   db,
		chainSpec: chainSpec,
		logger:    logger,
	}
}

// IsDataAvailable ensures that all blobs referenced in the block are
// stored before it returns without an error.
func (s *Store[BeaconBlockBodyT]) IsDataAvailable(
	_ context.Context,
	slot math.Slot,
	body BeaconBlockBodyT,
) bool {
	for _, commitment := range body.GetBlobKzgCommitments() {
		// Check if the block data is available in the IndexDB
		blockData, err := s.IndexDB.Has(slot.Unwrap(), commitment[:])
		if err != nil || !blockData {
			return false
		}
	}
	return true
}

// Persist ensures the sidecar data remains accessible, utilizing parallel
// processing for efficiency.
func (s *Store[BeaconBlockT]) Persist(
	slot math.Slot,
	sidecars *types.BlobSidecars,
) error {
	// Exit early if there are no sidecars to store.
	if sidecars.IsNil() || sidecars.Len() == 0 {
		return nil
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

	// Create error channel and wait group for parallel processing
	errChan := make(chan error, len(sidecars.Sidecars))
	var wg sync.WaitGroup

	// Create a list of commitments for this slot. We need to store the
	// commitments for this slot because we key each sidecar by its
	// commitment in the DB, and so this is necessary to retrieve the
	// sidecars later in GetBlobsFromStore.
	commitments := make([][]byte, len(sidecars.Sidecars))

	// Process and store sidecars in parallel, and collect commitments
	for i, sidecar := range sidecars.Sidecars {
		if sidecar == nil {
			return ErrAttemptedToStoreNilSidecar
		}

		wg.Add(1)
		go func(index int, sc *types.BlobSidecar) {
			defer wg.Done()

			bz, err := sc.MarshalSSZ()
			if err != nil {
				errChan <- err
				return
			}

			// Store the sidecar
			if err := s.IndexDB.Set(slot.Unwrap(), sc.KzgCommitment[:], bz); err != nil {
				errChan <- err
				return
			}

			// Store the commitment for the slot index. This is thread-safe
			// since every goroutine writes to a different index in the
			// commitments slice.
			commitments[index] = sc.KzgCommitment[:]
		}(i, sidecar)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	slotCommitments := &types.SlotCommitments{
		Commitments: commitments,
	}

	serializedCommitments, err := slotCommitments.MarshalSSZ()
	if err != nil {
		return err
	}

	// Store the commitments.
	if err := s.IndexDB.Set(slot.Unwrap(), []byte(SLOT_COMMITMENTS_KEY), serializedCommitments); err != nil {
		return err
	}

	s.logger.Info("Successfully stored all blob sidecars 🚗",
		"slot", slot.Base10(), "num_sidecars", sidecars.Len(),
	)
	return nil
}

// GetBlobsFromStore returns all blob sidecars for a given slot.
func (s *Store[BeaconBlockT]) GetBlobsFromStore(
	slot math.Slot,
) (*types.BlobSidecars, error) {
	// Get the commitment list for this slot
	serializedCommitments, err := s.IndexDB.Get(slot.Unwrap(), []byte(SLOT_COMMITMENTS_KEY))
	if err != nil {
		return &types.BlobSidecars{Sidecars: make([]*types.BlobSidecar, 0)}, nil // Return empty if not found
	}

	slotCommitments := &types.SlotCommitments{}
	if err := slotCommitments.UnmarshalSSZ(serializedCommitments); err != nil {
		return nil, err
	}
	commitments := slotCommitments.Commitments

	// Create error channel and wait group for parallel processing
	errChan := make(chan error, len(commitments))
	var wg sync.WaitGroup

	// Create slice to hold all sidecars
	sidecars := make([]*types.BlobSidecar, len(commitments))

	// Retrieve and unmarshal sidecars in parallel
	for i, commitment := range commitments {
		wg.Add(1)
		go func(index int, comm []byte) {
			defer wg.Done()

			// Get the sidecar bytes from the db
			bz, err := s.IndexDB.Get(slot.Unwrap(), comm)
			if err != nil {
				errChan <- err
				return
			}

			// Unmarshal the sidecar
			sidecar := new(types.BlobSidecar)
			if err := sidecar.UnmarshalSSZ(bz); err != nil {
				errChan <- err
				return
			}

			// Safely store the sidecar in the slice. This is thread-safe
			// since every goroutine writes to a different index in the
			// sidecars slice.
			sidecars[index] = sidecar
		}(i, commitment)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return &types.BlobSidecars{Sidecars: sidecars}, nil
}
