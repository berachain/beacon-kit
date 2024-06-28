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

package cache

import (
	"sync"
)

// historicalPayloadIDCacheSize defines the maximum number of slots to retain
// in the cache. Beyond this number, older slots will be pruned to manage
// memory usage.
const historicalPayloadIDCacheSize = 2

// PayloadIDCache provides a mechanism to store and retrieve payload IDs based
// on slot and parent block hash. It is designed to improve the efficiency of
// payload ID retrieval by caching recent entries.
type PayloadIDCache[
	PayloadIDT ~[8]byte, RootT ~[32]byte, SlotT ~uint64,
] struct {
	// mu protects access to the slotToStateRootToPayloadID map.
	mu sync.RWMutex
	// slotToStateRootToPayloadID is used for storing payload ID mappings
	slotToStateRootToPayloadID map[SlotT]map[RootT]PayloadIDT
}

// NewPayloadIDCache initializes and returns a new instance of PayloadIDCache.
// It prepares the internal data structures for storing payload ID mappings.
func NewPayloadIDCache[
	PayloadIDT ~[8]byte, RootT ~[32]byte, SlotT ~uint64,
]() *PayloadIDCache[PayloadIDT, RootT, SlotT] {
	return &PayloadIDCache[PayloadIDT, RootT, SlotT]{
		mu: sync.RWMutex{},
		slotToStateRootToPayloadID: make(
			map[SlotT]map[RootT]PayloadIDT,
		),
	}
}

// Has retrieves the payload ID associated with a given slot and eth1 hash.
// Has checks if a payload ID exists for a given slot and eth1 hash.
func (p *PayloadIDCache[_, RootT, SlotT]) Has(
	slot SlotT,
	stateRoot RootT,
) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.slotToStateRootToPayloadID[slot][stateRoot]
	return ok
}

// Get was successful.
func (p *PayloadIDCache[PayloadIDT, RootT, SlotT]) Get(
	slot SlotT,
	stateRoot RootT,
) (PayloadIDT, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	innerMap, ok := p.slotToStateRootToPayloadID[slot]
	if !ok {
		return PayloadIDT{}, false
	}
	pid, ok := innerMap[stateRoot]
	if !ok {
		return PayloadIDT{}, false
	}
	return pid, true
}

// Set updates or inserts a payload ID for a given slot and eth1 hash.
// It also prunes entries in the cache that are older than the
// historicalPayloadIDCacheSize limit.
func (p *PayloadIDCache[PayloadIDT, RootT, SlotT]) Set(
	slot SlotT, stateRoot RootT, pid PayloadIDT,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Prune older slots to maintain the cache size limit.
	if slot >= historicalPayloadIDCacheSize {
		p.prunePrior(slot - historicalPayloadIDCacheSize)
	}

	// Update the cache with the new payload ID.
	innerMap, exists := p.slotToStateRootToPayloadID[slot]
	if !exists {
		innerMap = make(map[RootT]PayloadIDT)
		p.slotToStateRootToPayloadID[slot] = innerMap
	}
	innerMap[stateRoot] = pid
}

// UnsafePrunePrior removes payload IDs from the cache for slots less than
// the specified slot. Only used for testing.
func (p *PayloadIDCache[_, _, SlotT]) UnsafePrunePrior(
	slot SlotT,
) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.prunePrior(slot)
}

// Prune removes payload IDs from the cache for slots less than the specified
// slot. This method helps in managing the memory usage of the cache by
// discarding outdated entries.
func (p *PayloadIDCache[_, _, SlotT]) prunePrior(slot SlotT) {
	for s := range p.slotToStateRootToPayloadID {
		if s < slot {
			delete(p.slotToStateRootToPayloadID, s)
		}
	}
}
