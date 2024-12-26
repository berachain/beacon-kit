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
	"github.com/berachain/beacon-kit/chain-spec/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/math"
)

// Store is the default implementation of the AvailabilityStore.
type Store struct {
	// IndexDB is a basic database interface.
	IndexDB
	// logger is used for logging.
	logger log.Logger
	// chainSpec contains the chain specification.
	chainSpec chain.ChainSpec
}

// New creates a new instance of the AvailabilityStore.
func New(
	db IndexDB,
	logger log.Logger,
	chainSpec chain.ChainSpec,
) *Store {
	return &Store{
		IndexDB:   db,
		chainSpec: chainSpec,
		logger:    logger,
	}
}

// IsDataAvailable ensures that all blobs referenced in the block are
// stored before it returns without an error.
func (s *Store) IsDataAvailable(
	_ context.Context,
	slot math.Slot,
	body *ctypes.BeaconBlockBody,
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

// GetBlobSidecars fetches the sidecars for a specific slot.
func (s *Store) GetBlobSidecars(slot math.Slot) (types.BlobSidecars, error) {
	sidecarBzs, err := s.IndexDB.GetByIndex(slot.Unwrap())
	if err != nil {
		return nil, err
	}

	sidecars := make(types.BlobSidecars, 0, len(sidecarBzs))
	for _, sidecarBz := range sidecarBzs {
		sidecar := types.BlobSidecar{}
		err = sidecar.UnmarshalSSZ(sidecarBz)
		if err != nil {
			return sidecars, err
		}
		sidecars = append(sidecars, &sidecar)
	}

	return sidecars, nil
}

// Persist ensures the sidecar data remains accessible, utilizing parallel
// processing for efficiency.
func (s *Store) Persist(
	slot math.Slot,
	sidecars types.BlobSidecars,
) error {
	// Exit early if there are no sidecars to store.
	if sidecars.IsNil() || len(sidecars) == 0 {
		return nil
	}

	// Check to see if we are required to store the sidecar anymore, if
	// this sidecar is from outside the required DA period, we can skip it.
	if !s.chainSpec.WithinDAPeriod(
		// slot in which the sidecar was included.
		// (Safe to assume all sidecars are in same slot at this point).
		sidecars[0].SignedBeaconBlockHeader.Header.GetSlot(),
		// current slot
		slot,
	) {
		return nil
	}

	// Store each sidecar sequentially. The store's underlying RangeDB is not
	// built to handle concurrent writes.
	for _, sidecar := range sidecars {
		sc := sidecar
		if sc == nil {
			return ErrAttemptedToStoreNilSidecar
		}
		bz, err := sc.MarshalSSZ()
		if err != nil {
			return err
		}
		err = s.IndexDB.Set(slot.Unwrap(), sc.KzgCommitment[:], bz)

		if err != nil {
			return err
		}
	}

	s.logger.Info("Successfully stored all blob sidecars 🚗",
		"slot", slot.Base10(), "num_sidecars", len(sidecars),
	)
	return nil
}
