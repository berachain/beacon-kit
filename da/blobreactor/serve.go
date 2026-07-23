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

package blobreactor

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cometbft/cometbft/p2p"
)

// handleByRootRequest serves the sidecars of one block, from the push cache (tip proposals, including the proposer's own current
// proposal) or the availability store. An empty response signals a miss; the requester never treats it as success.
//
// Only entries bound to a real proposal are served. An unverified push is internally consistent but its block signature is unchecked,
// so a poisoned push raced into the cache would fail the requester's verification and cost us their score. Entries flip to verified
// as soon as our own ProcessProposal binds them, so declining in the meantime just sends the requester to the next peer.
func (br *BlobReactor) handleByRootRequest(src p2p.Peer, req *SidecarsByRootRequest) {
	var chunks [][]byte

	br.pushMu.RLock()
	entry := br.pushCache.get(req.BlockRoot)
	br.pushMu.RUnlock()

	switch {
	case entry != nil && entry.verified:
		chunks = [][]byte{entry.raw}
	default:
		// The store is slot-indexed. If the stored sidecars belong to a different root than requested, the requester's verification will reject
		// them and move on; recomputing the root here is not worth it.
		raw, err := br.blobStore.GetByIndex(req.Slot.Unwrap())
		if err != nil {
			br.logger.Error("Failed to read sidecars from store",
				"slot", req.Slot.Unwrap(), "request_id", req.RequestID, "error", err)
		} else if len(raw) > 0 {
			chunks = [][]byte{encodeBlobSidecarsSSZ(raw)}
		}
	}

	br.metrics.observeServed("by_root", len(chunks) > 0)
	br.respond(src, req.RequestID, chunks)
}

// handleByRangeRequest serves complete per-slot sidecar sets for the requested slot range, up to the response byte budget. Slots with no
// data are skipped; the requester knows which slots it still needs.
func (br *BlobReactor) handleByRangeRequest(src p2p.Peer, req *SidecarsByRangeRequest) {
	var (
		chunks   [][]byte
		headSlot = br.getHeadSlot()
		budget   = responseByteBudget
	)

	endSlot := req.StartSlot + math.Slot(req.Count)
	for slot := req.StartSlot; slot < endSlot; slot++ {
		if headSlot > 0 && slot > headSlot {
			break
		}
		raw, err := br.blobStore.GetByIndex(slot.Unwrap())
		if err != nil {
			br.logger.Error("Failed to read sidecars from store",
				"slot", slot.Unwrap(), "request_id", req.RequestID, "error", err)
			continue
		}
		if len(raw) == 0 {
			continue
		}
		chunk := encodeBlobSidecarsSSZ(raw)
		if len(chunk) > budget {
			// Only complete slots are served; the requester re-requests the rest.
			break
		}
		budget -= len(chunk)
		chunks = append(chunks, chunk)
		if len(chunks) == maxChunksPerResponse {
			break
		}
	}

	br.metrics.observeServed("by_range", len(chunks) > 0)
	br.respond(src, req.RequestID, chunks)
}

// respond sends best-effort: a blocking send here would let a requester that never drains its recv queue pin
// a shared worker for the full p2p send timeout, starving inbound pushes. A dropped response just means the
// requester times out and retries another peer.
func (br *BlobReactor) respond(src p2p.Peer, requestID uint64, chunks [][]byte) {
	resp := &SidecarsResponse{
		RequestID:     requestID,
		SidecarChunks: chunks,
	}
	if !br.trySendToPeer(src, MessageTypeResponse, resp) {
		br.logger.Warn("Failed to send sidecars response", "peer", src.ID(), "request_id", requestID)
	}
}

// encodeBlobSidecarsSSZ combines individually SSZ-encoded BlobSidecar values into a single SSZ-encoded BlobSidecars list: a 4-byte offset
// (always 4) followed by the concatenated fixed-size sidecars.
//
//nolint:mnd // SSZ offset size
func encodeBlobSidecarsSSZ(sidecarBzs [][]byte) []byte {
	totalSize := 4
	for _, data := range sidecarBzs {
		totalSize += len(data)
	}

	result := make([]byte, totalSize)
	binary.LittleEndian.PutUint32(result[0:4], 4)

	pos := 4
	for _, data := range sidecarBzs {
		pos += copy(result[pos:], data)
	}
	return result
}
