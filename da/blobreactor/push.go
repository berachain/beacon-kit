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
	"context"
	"fmt"
	"time"

	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cometbft/cometbft/p2p"
)

// pushVerifyTimeout bounds the self-contained verification of a pushed sidecar set (inclusion proofs + KZG batch).
const pushVerifyTimeout = 5 * time.Second

// pushSlotTolerance is how far ahead of our finalized head a pushed sidecar set may be. Proposals live at
// head+1; the slack covers a node that is a few blocks behind processing its backlog.
const pushSlotTolerance = math.Slot(8)

// maxUnverifiedPushesPerPeer caps how many push-cache entries a single peer can occupy before its sidecars
// were bound to a real proposal, so fabricated pushes only compete with their sender's own entries for cache
// space and can never flush the cache.
const maxUnverifiedPushesPerPeer = 2

// pushEntry is one block's sidecars held for the tip of the chain.
type pushEntry struct {
	root     common.Root
	raw      []byte // SSZ-encoded BlobSidecars, ready for forwarding/serving
	sidecars datypes.BlobSidecars
	// src is the peer the entry came from; empty for our own proposals and proposal-bound sets.
	src p2p.ID
	// verified marks entries bound to a real proposal (our own, or fully verified in ProcessProposal).
	// Unverified entries are evicted first and are never re-gossiped.
	verified bool
}

// pushCache holds sidecars for the most recent proposals, keyed by block root, in insertion order. Guarded by
// BlobReactor.pushMu.
type pushCache struct {
	capacity int
	entries  map[common.Root]*pushEntry
	order    []common.Root // insertion order, oldest first
}

func newPushCache(capacity int) *pushCache {
	return &pushCache{
		capacity: capacity,
		entries:  make(map[common.Root]*pushEntry, capacity),
	}
}

func (pc *pushCache) get(root common.Root) *pushEntry {
	return pc.entries[root]
}

func (pc *pushCache) remove(root common.Root) {
	if _, ok := pc.entries[root]; !ok {
		return
	}
	delete(pc.entries, root)
	for i, r := range pc.order {
		if r == root {
			pc.order = append(pc.order[:i], pc.order[i+1:]...)
			break
		}
	}
}

// add inserts an entry. Unverified entries are evicted before verified ones, so a fabricated push can never
// displace our own proposal or a set already bound to one, and a single peer's unverified entries are capped
// at maxUnverifiedPushesPerPeer (its oldest is replaced beyond that).
func (pc *pushCache) add(entry *pushEntry) {
	if _, ok := pc.entries[entry.root]; ok {
		return
	}

	if !entry.verified && entry.src != "" {
		pc.capUnverifiedForPeer(entry.src)
	}
	if len(pc.order) >= pc.capacity {
		pc.evictOne()
	}

	pc.entries[entry.root] = entry
	pc.order = append(pc.order, entry.root)
}

// capUnverifiedForPeer removes the peer's oldest unverified entry once it hits its allowance.
func (pc *pushCache) capUnverifiedForPeer(src p2p.ID) {
	var (
		unverified int
		oldest     common.Root
		found      bool
	)
	for _, r := range pc.order {
		if e := pc.entries[r]; e != nil && !e.verified && e.src == src {
			if !found {
				oldest, found = r, true
			}
			unverified++
		}
	}
	if unverified >= maxUnverifiedPushesPerPeer {
		pc.remove(oldest)
	}
}

// evictOne removes the oldest unverified entry, or the oldest overall if every entry is verified.
func (pc *pushCache) evictOne() {
	evict := pc.order[0]
	for _, r := range pc.order {
		if e := pc.entries[r]; e != nil && !e.verified {
			evict = r
			break
		}
	}
	pc.remove(evict)
}

// sszUnmarshalSidecars decodes SSZ-encoded BlobSidecars.
func sszUnmarshalSidecars(bz []byte, sidecars *datypes.BlobSidecars) error {
	return ssz.Unmarshal(bz, sidecars)
}

// commitmentsOf collects the KZG commitments carried by the sidecars themselves, ordered by sidecar index position.
func commitmentsOf(sidecars datypes.BlobSidecars) eip4844.KZGCommitments[common.ExecutionHash] {
	commitments := make(eip4844.KZGCommitments[common.ExecutionHash], len(sidecars))
	for i, sc := range sidecars {
		commitments[i] = sc.GetKzgCommitment()
	}
	return commitments
}

// BroadcastSidecars is called by the proposer when it hands a block to CometBFT. It stores the sidecars in the push cache (so the
// proposer can serve them and finds them in its own ProcessProposal) and pushes them to all peers that are not known to have them. It
// never blocks on peer I/O.
func (br *BlobReactor) BroadcastSidecars(sidecarsBz []byte) error {
	var sidecars datypes.BlobSidecars
	if err := sszUnmarshalSidecars(sidecarsBz, &sidecars); err != nil {
		return fmt.Errorf("broadcast: failed to unmarshal own sidecars: %w", err)
	}
	if len(sidecars) == 0 {
		return nil
	}

	root := sidecars[0].GetBeaconBlockHeader().HashTreeRoot()
	entry := &pushEntry{root: root, raw: sidecarsBz, sidecars: sidecars, verified: true}

	br.pushMu.Lock()
	br.pushCache.add(entry)
	br.pushMu.Unlock()

	br.logger.Info("Broadcasting blob sidecars",
		"block_root", root, "slot", sidecars[0].GetBeaconBlockHeader().GetSlot().Unwrap(), "count", len(sidecars))

	go br.forwardPush(entry)
	return nil
}

// NotifySidecarsObtained is called when the node obtained and fully verified the sidecars for a block at the
// tip through a non-push path (its own EL or a by-root fetch). It caches them as the authoritative copy for
// serving and re-gossips them to peers that lack them.
//
// It always replaces any existing entry for the root rather than trusting it: a push that was cached but never
// bound to a real proposal (e.g. one carrying a bad signature) may still occupy this root, and we must not
// keep serving or forwarding that unverified data.
func (br *BlobReactor) NotifySidecarsObtained(root common.Root, sidecars datypes.BlobSidecars) {
	if len(sidecars) == 0 {
		return
	}

	raw, err := sidecars.MarshalSSZ()
	if err != nil {
		br.logger.Error("Failed to marshal obtained sidecars", "block_root", root, "error", err)
		return
	}

	entry := &pushEntry{root: root, raw: raw, sidecars: sidecars, verified: true}
	br.pushMu.Lock()
	br.pushCache.remove(root)
	br.pushCache.add(entry)
	br.pushMu.Unlock()

	go br.forwardPush(entry)
}

// DiscardPushedSidecars evicts a cached push for the given root. It is called when a cached push failed
// verification against the actual proposal, so the poisoned entry no longer shadows an honest re-push and is
// no longer served to peers.
func (br *BlobReactor) DiscardPushedSidecars(root common.Root) {
	br.pushMu.Lock()
	br.pushCache.remove(root)
	br.pushMu.Unlock()
}

// GetPushedSidecars returns the cached sidecars for a block root, or nil.
func (br *BlobReactor) GetPushedSidecars(root common.Root) datypes.BlobSidecars {
	br.pushMu.RLock()
	defer br.pushMu.RUnlock()
	if entry := br.pushCache.get(root); entry != nil {
		return entry.sidecars
	}
	return nil
}

// forwardPush sends the entry's sidecars to every peer not known to have them, marking those peers as having the root.
//
// Peers already known to have the root get a Have announcement instead, because knownRoots knowledge is one-directional. A peer that
// pushed or announced this root to us proved that IT holds the data, but it may not know that WE now do, and at its own forward step
// it would push the full ~768 KiB payload back to us. The 36-byte Have closes that gap. For peers that learned the root from us (we
// pushed or announced to them earlier) the extra Have is redundant but harmless.
func (br *BlobReactor) forwardPush(entry *pushEntry) {
	br.stateMu.RLock()
	pushTargets := make([]p2p.ID, 0, len(br.peers))
	haveTargets := make([]p2p.ID, 0, len(br.peers))
	for id, ps := range br.peers {
		if ps.knownRoots.has(entry.root) {
			haveTargets = append(haveTargets, id)
		} else {
			pushTargets = append(pushTargets, id)
		}
	}
	br.stateMu.RUnlock()

	push := &SidecarsPush{BlockRoot: entry.root, SidecarData: entry.raw}
	sent := 0
	for _, id := range pushTargets {
		peer := br.getPeer(id)
		if peer == nil {
			continue
		}
		if br.trySendToPeer(peer, MessageTypePush, push) {
			// The peer now has (or will shortly have) the sidecars; suppress future pushes of this root to it.
			br.markPeerHasRoot(id, entry.root)
			br.metrics.observePushSent()
			sent++
		}
	}

	// Tell peers that already hold the set that we now hold it too, so they do not push it back to us.
	have := &SidecarsHave{BlockRoot: entry.root}
	for _, id := range haveTargets {
		if peer := br.getPeer(id); peer != nil {
			br.trySendToPeer(peer, MessageTypeHave, have)
		}
	}
	br.logger.Debug("Forwarded blob sidecars push",
		"block_root", entry.root, "pushed", sent, "announced", len(haveTargets))
}

// announceHave tells every peer except src that we hold the complete sidecar set for root. Peers about to re-push the payload skip us
// once the announcement lands, so a ~36-byte message replaces a ~768 KiB duplicate, and by-root fetchers learn who to ask.
func (br *BlobReactor) announceHave(root common.Root, src p2p.ID) {
	br.stateMu.RLock()
	targets := make([]p2p.ID, 0, len(br.peers))
	for id := range br.peers {
		if id != src {
			targets = append(targets, id)
		}
	}
	br.stateMu.RUnlock()

	have := &SidecarsHave{BlockRoot: root}
	for _, id := range targets {
		if peer := br.getPeer(id); peer != nil {
			br.trySendToPeer(peer, MessageTypeHave, have)
		}
	}
}

// handlePush processes an unsolicited sidecars delivery. The sidecars are self-verified (structure, a slot
// just ahead of our head, inclusion proofs against their own header, KZG proofs against their own commitments)
// before being cached, but they are NOT forwarded here. Self-verification proves internal consistency, not
// that a real proposal exists: an attacker can fabricate a consistent set for a made-up header. So re-gossip
// only happens once the set is bound to a verified proposal (NotifySidecarsObtained), fabricated sets are
// never amplified through honest nodes, and per-peer caching limits mean they only compete with their own
// sender's entries for cache space.
func (br *BlobReactor) handlePush(src p2p.Peer, push *SidecarsPush) {
	// Whatever the outcome, the peer claims to have this root; remember that so we do not push it back and can fetch from it by root.
	br.markPeerHasRoot(src.ID(), push.BlockRoot)

	br.pushMu.RLock()
	known := br.pushCache.get(push.BlockRoot) != nil
	br.pushMu.RUnlock()
	if known {
		// Counted so the duplication factor of the push flood is measurable; each of these is a full payload
		// received for data we already hold.
		br.metrics.observePush("duplicate")
		return
	}

	sidecars, err := br.decodeSidecarsChunk(push.SidecarData)
	if err != nil {
		br.logger.Warn("Rejecting invalid sidecars push", "peer", src.ID(), "error", err)
		br.adjustScore(src.ID(), scoreJunk)
		br.metrics.observePush("invalid")
		return
	}

	// Only cache pushes for slots just ahead of our finalized head: that is the only region proposals can
	// live in, and it stops stale or far-future fabrications before the expensive KZG check. No score penalty,
	// an honest peer may simply be ahead of us while we catch up.
	if head := br.getHeadSlot(); head > 0 {
		slot := sidecars[0].GetBeaconBlockHeader().GetSlot()
		if slot <= head || slot > head+pushSlotTolerance {
			br.logger.Debug("Ignoring sidecars push outside the tip window",
				"peer", src.ID(), "slot", slot.Unwrap(), "head", head.Unwrap())
			br.metrics.observePush("slot_out_of_window")
			return
		}
	}

	// The push must be internally consistent with its claimed root.
	root := sidecars[0].GetBeaconBlockHeader().HashTreeRoot()
	if root != push.BlockRoot {
		br.logger.Warn("Sidecars push root mismatch", "peer", src.ID(), "claimed", push.BlockRoot, "actual", root)
		br.adjustScore(src.ID(), scoreJunk)
		br.metrics.observePush("invalid")
		return
	}

	// Self-contained verification: inclusion proofs bind every commitment to the (shared) header, KZG proofs bind every blob to its
	// commitment.
	ctx, cancel := context.WithTimeout(context.Background(), pushVerifyTimeout)
	defer cancel()
	if err = br.verifier.VerifySidecars(ctx, sidecars, sidecars[0].GetBeaconBlockHeader(), commitmentsOf(sidecars)); err != nil {
		br.logger.Warn("Sidecars push failed verification", "peer", src.ID(), "block_root", root, "error", err)
		br.adjustScore(src.ID(), scoreJunk)
		br.metrics.observePush("invalid")
		return
	}

	entry := &pushEntry{root: root, raw: push.SidecarData, sidecars: sidecars, src: src.ID()}
	br.pushMu.Lock()
	br.pushCache.add(entry)
	br.pushMu.Unlock()

	// Announce before verification: peers finishing their own verification are about to re-push the full
	// payload to everyone not known to hold it, and this announcement arriving first is what turns those
	// duplicates into no-ops. The early claim is safe. An honest push toward us would be dropped anyway (the
	// root is cached), and by-root requests are served from verified entries only, so a requester acting on
	// the claim sees at worst a miss and tries the next peer.
	br.announceHave(root, src.ID())

	br.metrics.observePush("accepted")
	br.logger.Info("Accepted blob sidecars push",
		"peer", src.ID(), "block_root", root,
		"slot", sidecars[0].GetBeaconBlockHeader().GetSlot().Unwrap(), "count", len(sidecars))
}
