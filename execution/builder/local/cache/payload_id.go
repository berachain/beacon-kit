// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package cache

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// historicalPayloadIDCacheSize defines the maximum number of slots to retain in the cache.
// Beyond this number, older slots will be pruned to manage memory usage.
const historicalPayloadIDCacheSize = 2

// PayloadIDCache provides a mechanism to store and retrieve payload IDs based on slot and
// parent block hash. It is designed to improve the efficiency of payload ID retrieval by
// caching recent entries.
type PayloadIDCache struct {
	// mu protects access to the slotToEth1HashToPayloadID map.
	mu sync.RWMutex
	// slotToEth1HashToPayloadID maps a slot to a map of eth1 block hashes to payload IDs.
	slotToEth1HashToPayloadID map[primitives.Slot]map[common.Hash]primitives.PayloadID
}

// NewPayloadIDCache initializes and returns a new instance of PayloadIDCache.
// It prepares the internal data structures for storing payload ID mappings.
func NewPayloadIDCache() *PayloadIDCache {
	return &PayloadIDCache{
		mu:                        sync.RWMutex{},
		slotToEth1HashToPayloadID: make(map[primitives.Slot]map[common.Hash]primitives.PayloadID),
	}
}

// Get retrieves the payload ID associated with a given slot and eth1 hash.
// It returns the found payload ID and a boolean indicating whether the lookup was successful.
func (p *PayloadIDCache) Get(
	slot primitives.Slot, eth1Hash common.Hash,
) (primitives.PayloadID, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	innerMap, ok := p.slotToEth1HashToPayloadID[slot]
	if !ok {
		return primitives.PayloadID{}, false
	}
	pid, ok := innerMap[eth1Hash]
	if !ok {
		return primitives.PayloadID{}, false
	}
	return pid, true
}

// Set updates or inserts a payload ID for a given slot and eth1 hash.
// It also prunes entries in the cache that are older than the historicalPayloadIDCacheSize limit.
func (p *PayloadIDCache) Set(
	slot primitives.Slot, eth1Hash common.Hash, pid primitives.PayloadID,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Prune older slots to maintain the cache size limit.
	if slot >= historicalPayloadIDCacheSize {
		p.prunePrior(slot - historicalPayloadIDCacheSize)
	}

	// Update the cache with the new payload ID.
	innerMap, exists := p.slotToEth1HashToPayloadID[slot]
	if !exists {
		innerMap = make(map[common.Hash]primitives.PayloadID)
		p.slotToEth1HashToPayloadID[slot] = innerMap
	}
	innerMap[eth1Hash] = pid
}

// UnsafePrunePrior removes payload IDs from the cache for slots less than
// the specified slot. Only used for testing.
func (p *PayloadIDCache) UnsafePrunePrior(slot primitives.Slot) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.prunePrior(slot)
}

// Prune removes payload IDs from the cache for slots less than the specified slot.
// This method helps in managing the memory usage of the cache by discarding outdated entries.
func (p *PayloadIDCache) prunePrior(slot primitives.Slot) {
	for s := range p.slotToEth1HashToPayloadID {
		if s < slot {
			delete(p.slotToEth1HashToPayloadID, s)
		}
	}
}
