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

package blockchain

import (
	"context"
	"fmt"

	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/log"
)

// blobFetchExecutor handles the Byzantine-critical blob fetch and verification logic.
// This is the core component that ensures we only accept valid blobs from peers.
type blobFetchExecutor struct {
	blobProcessor  BlobProcessor
	blobRequester  BlobRequester
	storageBackend StorageBackend
	logger         log.Logger
}

// FetchBlobsAndVerify fetches, verifies, and stores blobs for a single request.
// It creates a verifier function that the BlobRequester uses to validate blobs.
// If verification fails, the BlobRequester will automatically try the next peer.
func (e *blobFetchExecutor) FetchBlobsAndVerify(ctx context.Context, req BlobFetchRequest) error {
	e.logger.Info("Fetching blobs from peers", "slot", req.Header.Slot.Unwrap(), "expected_blobs", len(req.Commitments))

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Create a verifier function that validates blobs against the stored header and commitments.
	// This is the Byzantine fault tolerance mechanism - if a peer sends invalid blobs,
	// verification will fail and BlobRequester will try the next peer.
	verifier := func(sidecars datypes.BlobSidecars) error {
		return e.blobProcessor.VerifySidecars(ctx, sidecars, req.Header, req.Commitments)
	}

	// Request blobs with verification - will try multiple peers if verification fails
	fetchedBlobs, err := e.blobRequester.RequestBlobs(ctx, req.Header.Slot, verifier)
	if err != nil {
		return fmt.Errorf("failed to request valid blobs for slot %d: %w", req.Header.Slot.Unwrap(), err)
	}

	// Process and store the validated blobs
	err = e.blobProcessor.ProcessSidecars(e.storageBackend.AvailabilityStore(), fetchedBlobs)
	if err != nil {
		return fmt.Errorf("failed to process blobs for slot %d: %w", req.Header.Slot.Unwrap(), err)
	}

	e.logger.Info("Successfully fetched and stored blobs", "slot", req.Header.Slot.Unwrap(), "count", len(fetchedBlobs))
	return nil
}
