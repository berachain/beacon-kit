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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blobreactor

import (
	datypes "github.com/berachain/beacon-kit/da/types"
)

// BlobRequester is the interface that BlobReactor implements for BeaconKit to request blobs.
type BlobRequester interface {
	// RequestBlobs fetches all blobs for a given slot from peers.
	// Returns all blob sidecars for the slot, or an error if none could be retrieved.
	RequestBlobs(slot uint64) ([]*datypes.BlobSidecar, error)

	// SetHeadSlot updates the reactor's view of the current blockchain head slot.
	// Called by the blockchain service after processing each block.
	SetHeadSlot(slot uint64)
}

// BlobStore is a minimal interface for the BlobReactor to check and serve blobs.
// This matches the IndexDB interface from the AvailabilityStore.
type BlobStore interface {
	// Has checks if a blob exists for the given index and key.
	Has(index uint64, key []byte) (bool, error)

	// GetByIndex retrieves all raw blob data for a given index (slot).
	GetByIndex(index uint64) ([][]byte, error)
}
