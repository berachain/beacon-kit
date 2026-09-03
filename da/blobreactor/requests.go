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
	"cmp"
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"

	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cometbft/cometbft/p2p"
)

// newRequestID returns a cryptographically random request ID, so responses cannot be spoofed by guessing a sequential counter.
func newRequestID() uint64 {
	var buf [8]byte
	if _, err := crand.Read(buf[:]); err != nil {
		// crypto/rand never fails on supported platforms; fall back to time.
		return uint64(time.Now().UnixNano()) // #nosec G115 -- fallback only
	}
	return binary.LittleEndian.Uint64(buf[:])
}

// registerPending allocates a request ID bound to the given peer and returns it with the channel the response will be delivered on.
func (br *BlobReactor) registerPending(peerID p2p.ID) (uint64, chan *SidecarsResponse) {
	ch := make(chan *SidecarsResponse, 1)
	br.pendingMu.Lock()
	defer br.pendingMu.Unlock()
	for {
		id := newRequestID()
		if _, exists := br.pending[id]; exists {
			continue
		}
		br.pending[id] = pendingRequest{peerID: peerID, ch: ch}
		br.metrics.setActiveRequests(len(br.pending))
		return id, ch
	}
}

func (br *BlobReactor) unregisterPending(id uint64) {
	br.pendingMu.Lock()
	delete(br.pending, id)
	br.metrics.setActiveRequests(len(br.pending))
	br.pendingMu.Unlock()
}

// hasPendingFrom reports whether any in-flight request is bound to the given peer. The pending map holds at
// most a handful of entries, so a scan is fine.
func (br *BlobReactor) hasPendingFrom(peerID p2p.ID) bool {
	br.pendingMu.Lock()
	defer br.pendingMu.Unlock()
	for _, pending := range br.pending {
		if pending.peerID == peerID {
			return true
		}
	}
	return false
}

// handleResponse correlates a response with its in-flight request. Responses are only accepted from the exact peer the request was sent
// to; anything else is treated as a spoofing attempt and penalized without disturbing the still-pending request.
func (br *BlobReactor) handleResponse(src p2p.Peer, resp *SidecarsResponse) {
	br.pendingMu.Lock()
	pending, ok := br.pending[resp.RequestID]
	if ok && pending.peerID != src.ID() {
		br.pendingMu.Unlock()
		br.logger.Warn("Dropping response from wrong peer",
			"request_id", resp.RequestID, "expected_peer", pending.peerID, "actual_peer", src.ID())
		br.adjustScore(src.ID(), scoreSpoofedResp)
		return
	}
	if ok {
		delete(br.pending, resp.RequestID)
		br.metrics.setActiveRequests(len(br.pending))
	}
	br.pendingMu.Unlock()

	if !ok {
		// Late response after timeout; benign.
		br.logger.Debug("No waiting request for response", "request_id", resp.RequestID, "peer", src.ID())
		return
	}

	select {
	case pending.ch <- resp:
	default:
	}
}

// roundTrip sends one request to one peer and waits for the matching response, bounded by the per-peer request timeout and the caller's
// context.
func (br *BlobReactor) roundTrip(
	ctx context.Context,
	peerID p2p.ID,
	msgType MessageType,
	buildMsg func(requestID uint64) sszMarshaler,
) (*SidecarsResponse, error) {
	peer := br.getPeer(peerID)
	if peer == nil {
		return nil, newFetchError(fmt.Errorf("peer %s not available", peerID), statusPeerNotFound)
	}

	requestID, respChan := br.registerPending(peerID)
	defer br.unregisterPending(requestID)

	timeoutCtx, cancel := context.WithTimeout(ctx, br.config.RequestTimeout)
	defer cancel()

	// peer.Send blocks up to CometBFT's 10s send-queue timeout when the channel is congested, which is not
	// bounded by our context. Run it off the critical path so a congested peer cannot hold ProcessProposal or
	// FinalizeBlock past the round budget; the send goroutine drains on its own.
	msg := buildMsg(requestID)
	sent := make(chan bool, 1)
	go func() { sent <- br.sendToPeer(peer, msgType, msg) }()

	select {
	case ok := <-sent:
		if !ok {
			return nil, newFetchError(fmt.Errorf("failed to send %s to peer %s", msgType, peerID), statusSendFailed)
		}
	case <-timeoutCtx.Done():
		return nil, newFetchError(fmt.Errorf("send to peer %s exceeded %v", peerID, br.config.RequestTimeout), statusTimeout)
	}

	select {
	case resp := <-respChan:
		return resp, nil
	case <-timeoutCtx.Done():
		if ctx.Err() != nil {
			return nil, newFetchError(fmt.Errorf("request to peer %s cancelled: %w", peerID, ctx.Err()), statusTimeout)
		}
		return nil, newFetchError(
			fmt.Errorf("request to peer %s timed out after %v", peerID, br.config.RequestTimeout), statusTimeout)
	}
}

// selectPeers returns fetch candidates: peers known to hold the wanted root first, then the rest, ordered by descending score with random
// tiebreaks.
func (br *BlobReactor) selectPeers(root *common.Root) []p2p.ID {
	type candidate struct {
		id      p2p.ID
		score   int
		hasRoot bool
	}

	br.stateMu.RLock()
	candidates := make([]candidate, 0, len(br.peers))
	for id, ps := range br.peers {
		// hasRoot comes from unvalidated Have/push gossip, so it only earns the fast lane for peers that are
		// not already in the penalty box: a junk peer cannot announce Have for everything and keep first pick
		// regardless of its accumulated failures.
		candidates = append(candidates, candidate{
			id:      id,
			score:   ps.score,
			hasRoot: root != nil && ps.score >= 0 && ps.knownRoots.has(*root),
		})
	}
	br.stateMu.RUnlock()

	// Shuffle first so equal (hasRoot, score) peers are tried in random order.
	rand.Shuffle(len(candidates), func(i, j int) { // #nosec G404 -- load balancing only
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	slices.SortStableFunc(candidates, func(a, b candidate) int {
		if a.hasRoot != b.hasRoot {
			if a.hasRoot {
				return -1
			}
			return 1
		}
		return cmp.Compare(b.score, a.score)
	})

	peerIDs := make([]p2p.ID, len(candidates))
	for i, c := range candidates {
		peerIDs[i] = c.id
	}
	return peerIDs
}

// RequestSidecarsByRoot fetches the complete, verified sidecar set of one block from peers. The verify callback must implement the
// caller's full acceptance check, including that the sidecar count matches the block's commitment count, so an empty or short response
// can never pass as success. Peers are tried in order until one returns a verifying set, the context expires, or all peers were tried.
func (br *BlobReactor) RequestSidecarsByRoot(
	ctx context.Context,
	slot math.Slot,
	blockRoot common.Root,
	verify func(datypes.BlobSidecars) error,
) (datypes.BlobSidecars, error) {
	start := time.Now()
	peerIDs := br.selectPeers(&blockRoot)
	if len(peerIDs) == 0 {
		return nil, ErrNoPeersAvailable
	}

	for _, peerID := range peerIDs {
		if ctx.Err() != nil {
			br.metrics.recordFetchDone("by_root", statusTimeout, start)
			return nil, fmt.Errorf("by-root request cancelled: %w", ctx.Err())
		}

		sidecars, err := br.requestByRootFromPeer(ctx, peerID, slot, blockRoot, verify)
		if err != nil {
			br.metrics.recordFetchAttempt("by_root", fetchErrStatus(err))
			br.adjustScore(peerID, scoreFailure)
			br.logger.Debug("By-root fetch attempt failed", "peer", peerID, "slot", slot.Unwrap(), "error", err)
			continue
		}

		br.metrics.recordFetchAttempt("by_root", statusSuccess)
		br.metrics.recordFetchDone("by_root", statusSuccess, start)
		br.adjustScore(peerID, scoreSuccess)
		br.markPeerHasRoot(peerID, blockRoot)
		br.logger.Info("Fetched blob sidecars by root",
			"peer", peerID, "slot", slot.Unwrap(), "block_root", blockRoot, "count", len(sidecars))
		return sidecars, nil
	}

	br.metrics.recordFetchDone("by_root", statusAllPeersFailed, start)
	return nil, fmt.Errorf("%w: slot %d root %s (%d peers tried)",
		ErrAllPeersFailed, slot.Unwrap(), blockRoot, len(peerIDs))
}

func (br *BlobReactor) requestByRootFromPeer(
	ctx context.Context,
	peerID p2p.ID,
	slot math.Slot,
	blockRoot common.Root,
	verify func(datypes.BlobSidecars) error,
) (datypes.BlobSidecars, error) {
	resp, err := br.roundTrip(ctx, peerID, MessageTypeByRootRequest, func(requestID uint64) sszMarshaler {
		return &SidecarsByRootRequest{RequestID: requestID, Slot: slot, BlockRoot: blockRoot}
	})
	if err != nil {
		return nil, err
	}

	// An empty response is a miss, never a success: the caller only requests blocks that are known to have blobs.
	if len(resp.SidecarChunks) == 0 {
		return nil, newFetchError(fmt.Errorf("peer %s has no sidecars for slot %d", peerID, slot.Unwrap()), statusMiss)
	}
	if len(resp.SidecarChunks) != 1 {
		return nil, newFetchError(
			fmt.Errorf("peer %s sent %d chunks for by-root request", peerID, len(resp.SidecarChunks)),
			statusInvalidResponse)
	}

	sidecars, err := br.decodeSidecarsChunk(resp.SidecarChunks[0])
	if err != nil {
		return nil, newFetchError(fmt.Errorf("peer %s: %w", peerID, err), statusInvalidResponse)
	}

	if root := sidecars[0].GetBeaconBlockHeader().HashTreeRoot(); root != blockRoot {
		return nil, newFetchError(
			fmt.Errorf("peer %s returned sidecars for root %s, wanted %s", peerID, root, blockRoot),
			statusInvalidResponse)
	}

	if err = verify(sidecars); err != nil {
		return nil, newFetchError(fmt.Errorf("sidecars from peer %s failed verification: %w", peerID, err), statusVerifyFailed)
	}
	return sidecars, nil
}

// RequestSidecarsByRange fetches sidecars for slots [start, start+count) from peers. The verify callback is invoked once per returned
// slot with that slot's complete sidecar set and must enforce the expected count for the slot. The first peer that yields at least one
// verified slot ends the attempt; the caller re-requests whatever is still missing on its next tick, so partial progress is never lost
// and never mistaken for completion.
func (br *BlobReactor) RequestSidecarsByRange(
	ctx context.Context,
	startSlot math.Slot,
	count uint64,
	verify func(math.Slot, datypes.BlobSidecars) error,
) (map[math.Slot]datypes.BlobSidecars, error) {
	if count == 0 {
		return map[math.Slot]datypes.BlobSidecars{}, nil
	}
	if count > MaxRequestedSlots {
		count = MaxRequestedSlots
	}

	start := time.Now()
	peerIDs := br.selectPeers(nil)
	if len(peerIDs) == 0 {
		return nil, ErrNoPeersAvailable
	}

	for _, peerID := range peerIDs {
		if ctx.Err() != nil {
			br.metrics.recordFetchDone("by_range", statusTimeout, start)
			return nil, fmt.Errorf("by-range request cancelled: %w", ctx.Err())
		}

		verified, err := br.requestByRangeFromPeer(ctx, peerID, startSlot, count, verify)
		if err != nil {
			br.metrics.recordFetchAttempt("by_range", fetchErrStatus(err))
			br.adjustScore(peerID, scoreFailure)
			br.logger.Debug("By-range fetch attempt failed",
				"peer", peerID, "start_slot", startSlot.Unwrap(), "count", count, "error", err)
			continue
		}

		br.metrics.recordFetchAttempt("by_range", statusSuccess)
		br.metrics.recordFetchDone("by_range", statusSuccess, start)
		br.adjustScore(peerID, scoreSuccess)
		br.logger.Info("Fetched blob sidecars by range",
			"peer", peerID, "start_slot", startSlot.Unwrap(), "count", count, "verified_slots", len(verified))
		return verified, nil
	}

	br.metrics.recordFetchDone("by_range", statusAllPeersFailed, start)
	return nil, fmt.Errorf("%w: slots [%d, %d) (%d peers tried)",
		ErrAllPeersFailed, startSlot.Unwrap(), startSlot.Unwrap()+count, len(peerIDs))
}

func (br *BlobReactor) requestByRangeFromPeer(
	ctx context.Context,
	peerID p2p.ID,
	startSlot math.Slot,
	count uint64,
	verify func(math.Slot, datypes.BlobSidecars) error,
) (map[math.Slot]datypes.BlobSidecars, error) {
	resp, err := br.roundTrip(ctx, peerID, MessageTypeByRangeRequest, func(requestID uint64) sszMarshaler {
		return &SidecarsByRangeRequest{RequestID: requestID, StartSlot: startSlot, Count: count}
	})
	if err != nil {
		return nil, err
	}

	if len(resp.SidecarChunks) == 0 {
		return nil, newFetchError(
			fmt.Errorf("peer %s has no sidecars in [%d, %d)", peerID, startSlot.Unwrap(), startSlot.Unwrap()+count),
			statusMiss)
	}

	verified := make(map[math.Slot]datypes.BlobSidecars)
	for _, chunk := range resp.SidecarChunks {
		sidecars, decodeErr := br.decodeSidecarsChunk(chunk)
		if decodeErr != nil {
			return nil, newFetchError(fmt.Errorf("peer %s: %w", peerID, decodeErr), statusInvalidResponse)
		}

		slot := sidecars[0].GetBeaconBlockHeader().GetSlot()
		if slot < startSlot || slot >= startSlot+math.Slot(count) {
			return nil, newFetchError(
				fmt.Errorf("peer %s returned slot %d outside requested range", peerID, slot.Unwrap()),
				statusInvalidResponse)
		}
		if _, dup := verified[slot]; dup {
			return nil, newFetchError(
				fmt.Errorf("peer %s returned duplicate chunk for slot %d", peerID, slot.Unwrap()),
				statusInvalidResponse)
		}

		verifyErr := verify(slot, sidecars)
		switch {
		case verifyErr == nil:
			verified[slot] = sidecars
		case errors.Is(verifyErr, ErrSlotNotRequested):
			// The peer cannot know which in-range slots we still need; returning one we did not ask for is not a fault.
			continue
		default:
			// A single bad slot poisons trust in the whole response.
			return nil, newFetchError(
				fmt.Errorf("sidecars for slot %d from peer %s failed verification: %w",
					slot.Unwrap(), peerID, verifyErr),
				statusVerifyFailed)
		}
	}

	if len(verified) == 0 {
		return nil, newFetchError(fmt.Errorf("peer %s returned no usable slots", peerID), statusMiss)
	}
	return verified, nil
}

func sortSidecarsByIndex(sidecars datypes.BlobSidecars) {
	slices.SortFunc(sidecars, func(a, b *datypes.BlobSidecar) int {
		return cmp.Compare(a.GetIndex(), b.GetIndex())
	})
}
