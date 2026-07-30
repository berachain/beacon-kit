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

package blockchain

import (
	"bytes"
	"context"
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
)

// fetchProposalSidecars retrieves and fully verifies the blob sidecars of a block at the tip of the chain when blob consensus is enabled
// (sidecars no longer ride as a consensus tx). It walks the delivery lanes described on blobreactor.BlobReactor, fastest first:
//
//  1. sidecars pushed over the blob reactor channel (the proposer pushes at
//     proposal time, verifying validators re-gossip),
//  2. the local execution client (blob txs gossip through the EL mempool with
//     their blobs attached, so this is usually a local hit), and
//  3. a by-root request to peers on the blob reactor channel.
//
// Every tier's candidate runs the same full verification against the block (count vs commitments, signature binding, inclusion proofs,
// KZG), so an empty or wrong set from one tier simply falls through to the next. The overall deadline is sized to fit within a consensus
// round.
func (s *Service) fetchProposalSidecars(
	ctx context.Context,
	signedBlk *ctypes.SignedBeaconBlock,
) (datypes.BlobSidecars, error) {
	blk := signedBlk.GetBeaconBlock()
	commitments := blk.GetBody().GetBlobKzgCommitments()
	if len(commitments) == 0 {
		return datypes.BlobSidecars{}, nil
	}

	fetchCtx, cancel := context.WithTimeout(ctx, s.blobRequester.FetchTimeout())
	defer cancel()

	var (
		blockRoot = blk.HashTreeRoot()
		slot      = blk.GetSlot()
		verify    = func(sidecars datypes.BlobSidecars) error {
			return verifySidecarsBinding(
				fetchCtx, s.blobProcessor, blk.GetHeader(), commitments, signedBlk.GetSignature(), sidecars)
		}
	)

	// Tier 1: sidecars pushed to us (or our own, if we are the proposer).
	if sidecars := s.blobRequester.GetPushedSidecars(blockRoot); sidecars != nil {
		if err := verify(sidecars); err == nil {
			s.logger.Info("Blob sidecars retrieved from push cache",
				"slot", slot.Unwrap(), "count", len(sidecars))
			// The set is now bound to a verified proposal; this marks it as such in the push cache and
			// re-gossips it (pushes are never forwarded before this binding).
			s.blobRequester.NotifySidecarsObtained(blockRoot, sidecars)
			return sidecars, nil
		}
		// The cached push does not match the real proposal (e.g. a poisoned copy that arrived first). Evict
		// it so it stops shadowing an honest re-push and is not served to peers, then fall through.
		s.logger.Warn("Pushed blob sidecars failed verification against proposal, discarding",
			"slot", slot.Unwrap())
		s.blobRequester.DiscardPushedSidecars(blockRoot)
	}

	// Tier 2: rebuild from the local execution client's blob pool.
	sidecars, err := s.blobReconstructor.ReconstructSidecars(fetchCtx, signedBlk)
	switch {
	case err != nil:
		s.logger.Debug("Could not reconstruct blob sidecars from execution client",
			"slot", slot.Unwrap(), "error", err)
	case verify(sidecars) == nil:
		s.logger.Info("Blob sidecars reconstructed from execution client",
			"slot", slot.Unwrap(), "count", len(sidecars))
		// Make them available to peers that are still missing them.
		s.blobRequester.NotifySidecarsObtained(blockRoot, sidecars)
		return sidecars, nil
	default:
		s.logger.Warn("EL-reconstructed blob sidecars failed verification",
			"slot", slot.Unwrap())
	}

	// Tier 3: ask peers, those known to have the data first.
	sidecars, err = s.blobRequester.RequestSidecarsByRoot(fetchCtx, slot, blockRoot, verify)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve blob sidecars for slot %d: %w", slot.Unwrap(), err)
	}
	s.blobRequester.NotifySidecarsObtained(blockRoot, sidecars)
	return sidecars, nil
}

// verifySidecarsBinding runs the full acceptance check binding a candidate sidecar set to a block's header, commitments and signature:
// exact count vs commitments (an empty or short set is a failure, never a success), block-signature equality on every sidecar, then
// header equality, inclusion proofs and KZG proofs via the processor. Sidecars must already be sorted by index. It serves both the
// tip-of-chain fetch (against the live proposal) and the background fetcher (against the request recorded at queue time).
func verifySidecarsBinding(
	ctx context.Context,
	processor BlobProcessor,
	header *ctypes.BeaconBlockHeader,
	commitments eip4844.KZGCommitments[common.ExecutionHash],
	signature crypto.BLSSignature,
	sidecars datypes.BlobSidecars,
) error {
	if len(sidecars) != len(commitments) {
		return fmt.Errorf("expected %d sidecars for slot %d, got %d: %w",
			len(commitments), header.GetSlot().Unwrap(), len(sidecars),
			ErrSidecarCommitmentMismatch,
		)
	}
	for i, sidecar := range sidecars {
		sidecarSignature := sidecar.GetSignature()
		if !bytes.Equal(signature[:], sidecarSignature[:]) {
			return fmt.Errorf("%w, idx: %d", ErrSidecarSignatureMismatch, i)
		}
	}
	return processor.VerifySidecars(ctx, sidecars, header, commitments)
}

// sidecarsAlreadyStored reports whether the availability store already holds the block's complete sidecar set: exact count vs
// commitments and every stored sidecar bound to the block's header. A count or availability check alone is not enough, the store is
// slot-indexed and could hold sidecars from a different proposal at the same slot with identical commitments, whose embedded header
// and signature would be wrong.
func sidecarsAlreadyStored(
	store *dastore.Store,
	header *ctypes.BeaconBlockHeader,
	commitments eip4844.KZGCommitments[common.ExecutionHash],
) bool {
	sidecars, err := store.GetBlobSidecars(header.GetSlot())
	if err != nil || len(sidecars) != len(commitments) {
		return false
	}
	root := header.HashTreeRoot()
	for _, sidecar := range sidecars {
		if sidecar.GetBeaconBlockHeader().HashTreeRoot() != root {
			return false
		}
	}
	return true
}
