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
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/sourcegraph/conc/iter"
)

// Constants for key prefixes and encoding.
const (
	blobKeyPrefix       = "blob-"
	headerKeyPrefix     = "header-"
	kzgKeyPrefix        = "kzg-"
	proofKeyPrefix      = "proof-"
	commitmentKeyPrefix = "commitment-"

	// Size constants for key components.
	slotSize  = 8 // size of uint64 in bytes
	indexSize = 4 // size of uint32 in bytes
)

// Store is the default implementation of the AvailabilityStore.
type Store[BeaconBlockBodyT BeaconBlockBody, BeaconBlockHeaderT any] struct {
	// IndexDB is a basic database interface.
	IndexDB
	// logger is used for logging.
	logger log.Logger
	// chainSpec contains the chain specification.
	chainSpec common.ChainSpec
}

// New creates a new instance of the AvailabilityStore.
func New[BeaconBlockT BeaconBlockBody, BeaconBlockHeaderT any](
	db IndexDB,
	logger log.Logger,
	chainSpec common.ChainSpec,
) *Store[BeaconBlockT, BeaconBlockHeaderT] {
	return &Store[BeaconBlockT, BeaconBlockHeaderT]{
		IndexDB:   db,
		chainSpec: chainSpec,
		logger:    logger,
	}
}

// IsDataAvailable ensures that all blobs referenced in the block are
// stored before it returns without an error.
func (s *Store[BeaconBlockBodyT, _]) IsDataAvailable(
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
func (s *Store[BeaconBlockT, _]) Persist(
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

	// Store each sidecar in parallel.
	if err := errors.Join(iter.Map(
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
			return s.Set(slot.Unwrap(), sc.KzgCommitment[:], bz)
		},
	)...); err != nil {
		return err
	}

	s.logger.Info("Successfully stored all blob sidecars 🚗",
		"slot", slot.Base10(), "num_sidecars", sidecars.Len(),
	)
	return nil
}

// GetBlobSideCars retrieves blob sidecars for a given slot.
func (s *Store[BeaconBlockBodyT, BeaconBlockHeaderT]) GetBlobSideCars(
	slot math.U64,
) (*[]backend.BlobSideCar[BeaconBlockHeaderT], error) {
	// Implementation to fetch blob sidecars for the given slot
	// 1. Check if we have data for this slot
	// 2. Convert the stored data to BlobSidecarData format
	// 3. Return the result

	blobSidecars := make([]backend.BlobSideCar[BeaconBlockHeaderT], 0)
	// TODO: Implement actual data retrieval logic here
	// - Query the stored blobs for the given slot
	// - Convert them to BlobSidecarData format
	// - Handle any potential errors

	// First, check if we have any data for this slot

	hasData, err := s.Has(uint64(slot), []byte(blobKeyPrefix))
	if err != nil {
		return nil, fmt.Errorf("error checking slot data existence: %w", err)
	}
	s.logger.Info("hasData", "hasData", hasData)
	if !hasData {
		// No data for this slot
		return &blobSidecars, nil
	}

	// Retrieve the block header for this slot.
	headerBytes, err := s.Get(uint64(slot), []byte(headerKeyPrefix))
	if err != nil {
		return nil, fmt.Errorf("error retrieving block header: %w", err)
	}
	s.logger.Info("headerBytes", "headerBytes", headerBytes)

	// header, err := s.decodeBlockHeader(headerBytes)
	// if err != nil {
	//	return nil, fmt.Errorf("error decoding block header: %w", err)
	// }
	//

	// Retrieve the number of blobs for this slot
	countBytes, errInGet := s.Get(uint64(slot), []byte("count"))
	if errInGet != nil {
		return nil, fmt.Errorf("error retrieving blob count: %w", errInGet)
	}
	count := binary.BigEndian.Uint32(countBytes)

	// Retrieve each blob and its associated data.
	for i := range count {
		// Retrieve blob data
		var blobBytes []byte
		blobBytes, err = s.Get(uint64(slot), buildKey(blobKeyPrefix, slot, int(i)))
		if err != nil {
			return nil, fmt.Errorf("error retrieving blob %d: %w", i, err)
		}
		var proofBytes []byte
		// Retrieve KZG proof
		proofBytes, err = s.Get(uint64(slot), buildKey(proofKeyPrefix, slot, int(i)))
		if err != nil {
			return nil, fmt.Errorf("error retrieving KZG proof %d: %w", i, err)
		}
		var commitmentBytes []byte

		// Retrieve KZG commitment
		commitmentBytes, err = s.Get(
			uint64(slot),
			buildKey(commitmentKeyPrefix, slot, int(i)),
		)
		if err != nil {
			return nil, fmt.Errorf("error retrieving KZG commitment %d: %w", i, err)
		}

		// Create blob sidecar data
		var blob eip4844.Blob
		copy(blob[:], blobBytes)

		var proof eip4844.KZGProof
		copy(proof[:], proofBytes)

		var commitment eip4844.KZGCommitment
		copy(commitment[:], commitmentBytes)

		// sidecar := backend.BlobSideCar{
		//	Index:                       uint64(i),
		//	Blob:                        blob,
		//	KzgProof:                    proof,
		//	KzgCommitment:               commitment,
		//	BeaconBlockHeader:           beacontypes.BlockHeader[]{},
		//	KzgCommitmentInclusionProof: headerBytes,
		// }
		//
		// blobSidecars = append(blobSidecars, sidecar)
	}
	return &blobSidecars, nil
}

// Helper methods for encoding/decoding block headers.
// func (s *Store[BeaconBlockBodyT]) encodeBlockHeader(
// header beacontypes.BeaconBlockHeader,
// ) ([]byte, error) {
//	// Implement header encoding logic
//	return header.MarshalSSZ()
// }

// func (s *Store[BeaconBlockBodyT]) decodeBlockHeader(data []byte) (
// beacontypes.BeaconBlockHeader,
// error,
// ) {
//	// Implement header decoding logic
//	var header beacontypes.BeaconBlockHeader
//	err := header.UnmarshalSSZ(data)
//	return header, err
// }

// buildKey creates a composite key for storing blob-related data.
func buildKey(prefix string, slot math.U64, index int) []byte {
	// Calculate total key size using named constants
	keySize := len(prefix) + slotSize + indexSize
	key := make([]byte, keySize)

	// Copy prefix into key
	copy(key, prefix)

	// Write slot number using constant offset
	binary.BigEndian.PutUint64(key[len(prefix):], uint64(slot))

	// Write index using constant offset
	binary.BigEndian.PutUint32(key[len(prefix)+slotSize:], uint32(index))

	return key
}
