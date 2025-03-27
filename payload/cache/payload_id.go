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

package cache

import (
	"sync"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// historicalPayloadIDCacheSize defines the maximum number of slots to retain
// in the cache. Beyond this number, older slots will be pruned to manage
// memory usage.
const historicalPayloadIDCacheSize = 2

// PayloadIDCache provides a mechanism to store and retrieve payload IDs based
// on slot and parent block hash. It is designed to improve the efficiency of
// payload ID retrieval by caching recent entries.
type PayloadIDCache struct {
	// mu protects access to the slotToBlockRootToPayloadID map.
	mu sync.RWMutex
	// slotToBlockRootToPayloadID is used for storing payload ID mappings
	slotToBlockRootToPayloadID map[payloadIDCacheKey]PayloadIDCacheResult
}

// payloadIDCacheKey is the (slot, root) tuple that is used to access a
// payloadID from the cache.
type payloadIDCacheKey struct {
	slot math.Slot
	root common.Root
}

type PayloadIDCacheResult struct {
	PayloadID   engineprimitives.PayloadID
	ForkVersion common.Version
}

// NewPayloadIDCache initializes and returns a new instance of PayloadIDCache.
// It prepares the internal data structures for storing payload ID mappings.
func NewPayloadIDCache() *PayloadIDCache {
	return &PayloadIDCache{
		mu: sync.RWMutex{},
		slotToBlockRootToPayloadID: make(
			map[payloadIDCacheKey]PayloadIDCacheResult,
		),
	}
}

// Has retrieves the payload ID associated with a given slot and eth1 hash.
// Has checks if a payload ID exists for a given slot and eth1 hash.
func (p *PayloadIDCache) Has(
	slot math.Slot,
	blockRoot common.Root,
) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.slotToBlockRootToPayloadID[payloadIDCacheKey{slot, blockRoot}]
	return ok
}

// GetAndEvict retrieves the payloadID from the cache. If successfully retrieved,
// evict it from the cache.
func (p *PayloadIDCache) GetAndEvict(
	slot math.Slot,
	blockRoot common.Root,
) (PayloadIDCacheResult, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	key := payloadIDCacheKey{slot, blockRoot}
	pid, ok := p.slotToBlockRootToPayloadID[key]
	if !ok {
		return PayloadIDCacheResult{}, false
	}

	// Successfully retrieved. Remove from cache.
	delete(p.slotToBlockRootToPayloadID, key)
	return pid, true
}

// Set updates or inserts a payload ID for a given slot and eth1 hash.
// It also prunes entries in the cache that are older than the
// historicalPayloadIDCacheSize limit.
func (p *PayloadIDCache) Set(
	slot math.Slot, blockRoot common.Root,
	pid engineprimitives.PayloadID, version common.Version,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Prune older slots to maintain the cache size limit.
	if slot >= historicalPayloadIDCacheSize {
		p.prunePrior(slot - historicalPayloadIDCacheSize)
	}

	// Update the cache with the new payload ID.
	p.slotToBlockRootToPayloadID[payloadIDCacheKey{slot, blockRoot}] = PayloadIDCacheResult{
		PayloadID:   pid,
		ForkVersion: version,
	}
}

// prunePrior removes payload IDs from the cache for slots less than
// the specified slot. This method helps in managing the memory usage
// of the cache by discarding outdated entries.
func (p *PayloadIDCache) prunePrior(slot math.Slot) {
	for s := range p.slotToBlockRootToPayloadID {
		if s.slot < slot {
			delete(p.slotToBlockRootToPayloadID, s)
		}
	}
}
