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

package store

import (
	"context"
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
)

// Store is the default implementation of the AvailabilityStore.
type Store struct {
	// IndexDB is a basic database interface.
	IndexDB
	// logger is used for logging.
	logger log.Logger
}

// New creates a new instance of the AvailabilityStore.
func New(
	db IndexDB,
	logger log.Logger,
) *Store {
	return &Store{
		IndexDB: db,
		logger:  logger,
	}
}

// IsDataAvailable ensures that all blobs referenced in the block are
// stored before it returns without an error.
func (s *Store) IsDataAvailable(
	_ context.Context,
	slot math.Slot,
	body *ctypes.BeaconBlockBody,
) bool {
	// Commitments can be duplicated within a block, so the storage key
	// includes the blob index alongside the commitment.
	for i, commitment := range body.GetBlobKzgCommitments() {
		if i > maxBlobIndex {
			s.logger.Error("Blob index exceeds maximum storable value", "index", i)
			return false
		}
		blockData, err := s.IndexDB.Has(slot.Unwrap(), sidecarKey(commitment, uint64(i))) //#nosec:G115 // i>=0
		if err != nil || !blockData {
			return false
		}
	}
	return true
}

// maxBlobIndex is the largest blob index encodable in the storage key.
const maxBlobIndex = 255

// sidecarKey builds the storage key for one sidecar: the KZG commitment with
// the blob index appended, so duplicated commitments within a block do not
// overwrite each other.
func sidecarKey(commitment eip4844.KZGCommitment, index uint64) []byte {
	return append(commitment[:], byte(index)) // #nosec G115 -- callers enforce index <= maxBlobIndex (255)
}

// GetBlobSidecars fetches the sidecars for a specific slot.
func (s *Store) GetBlobSidecars(slot math.Slot) (types.BlobSidecars, error) {
	sidecarBzs, err := s.IndexDB.GetByIndex(slot.Unwrap())
	if err != nil {
		return nil, err
	}

	sidecars := make(types.BlobSidecars, 0, len(sidecarBzs))
	for _, sidecarBz := range sidecarBzs {
		sidecar := new(types.BlobSidecar)
		if err = ssz.Unmarshal(sidecarBz, sidecar); err != nil {
			return sidecars, err
		}
		sidecars = append(sidecars, sidecar)
	}

	return sidecars, nil
}

// Persist ensures the sidecar data remains accessible, utilizing parallel
// processing for efficiency. A block's sidecars are always persisted as a
// complete set, so the slot's existing entries are cleared first: this keeps
// Persist idempotent and, across the storage-key change (commitment ->
// commitment||index), prevents a re-persisted block from leaving both the old-
// and new-key entries and returning duplicate indices.
func (s *Store) Persist(sidecars types.BlobSidecars) error {
	if len(sidecars) == 0 {
		return nil
	}
	if sidecars[0] == nil {
		return ErrAttemptedToStoreNilSidecar
	}
	slot := sidecars[0].GetBeaconBlockHeader().GetSlot()
	if err := s.IndexDB.DeleteByIndex(slot.Unwrap()); err != nil {
		return err
	}

	// Store each sidecar sequentially. The store's underlying RangeDB is not
	// built to handle concurrent writes.
	for _, sidecar := range sidecars {
		if sidecar == nil {
			return ErrAttemptedToStoreNilSidecar
		}
		// Every sidecar must belong to the slot whose entries were just
		// deleted; a mixed-slot slice would delete one slot's data and then
		// write under another.
		if sidecarSlot := sidecar.GetBeaconBlockHeader().GetSlot(); sidecarSlot != slot {
			return fmt.Errorf("%w: expected %d, got %d",
				ErrMixedSlotSidecars, slot.Unwrap(), sidecarSlot.Unwrap())
		}
		index := sidecar.GetIndex()
		if index > maxBlobIndex {
			return fmt.Errorf("blob index %d exceeds maximum storable value %d", index, maxBlobIndex)
		}
		bz, err := sidecar.MarshalSSZ()
		if err != nil {
			return err
		}
		err = s.IndexDB.Set(slot.Unwrap(), sidecarKey(sidecar.KzgCommitment, index), bz)
		if err != nil {
			return err
		}
	}

	s.logger.Info("Successfully stored all blob sidecars 🚗",
		"slot", slot.Base10(), "num_sidecars", len(sidecars),
	)
	return nil
}

// DeleteBlobSidecars removes all blob sidecars for the specified slot.
func (s *Store) DeleteBlobSidecars(slot math.Slot) error {
	return s.IndexDB.DeleteByIndex(slot.Unwrap())
}
