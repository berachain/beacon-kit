// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

	"github.com/itsdevbear/bolaris/types/primitives"
)

// historicalPayloadIDCacheSize is the number of slots to keep in the cache.
const historicalPayloadIDCacheSize = 2

// PayloadIDCache is a cache that keeps track of the prepared payload ID for the
// given slot and with the given head root.
type PayloadIDCache struct {
	slotToEth1HashToPayloadID sync.Map
}

// NewPayloadIDCache returns a new payload ID cache.
func NewPayloadIDCache() *PayloadIDCache {
	return &PayloadIDCache{}
}

// PayloadID returns the payload ID for the given slot and parent block root.
func (p *PayloadIDCache) Get(slot primitives.Slot, eth1Hash [32]byte) (primitives.PayloadID, bool) {
	value, ok := p.slotToEth1HashToPayloadID.Load(slot)
	if !ok {
		return primitives.PayloadID{}, false
	}
	innerMap, ok := value.(*sync.Map)
	if !ok {
		return primitives.PayloadID{}, false
	}
	pid, ok := innerMap.Load(eth1Hash)
	if !ok {
		return primitives.PayloadID{}, false
	}
	return pid.(primitives.PayloadID), true
}

// SetPayloadID updates the payload ID for the given slot and head root
// it also prunes older entries in the cache.
func (p *PayloadIDCache) Set(slot primitives.Slot, eth1Hash [32]byte, pid primitives.PayloadID) {
	value, _ := p.slotToEth1HashToPayloadID.LoadOrStore(slot, &sync.Map{})
	innerMap := value.(*sync.Map)
	innerMap.Store(eth1Hash, pid)
	if slot > 1 {
		p.Prune(slot - historicalPayloadIDCacheSize)
	}
}

// Prune prunes old payload IDs.
func (p *PayloadIDCache) Prune(slot primitives.Slot) {
	p.slotToEth1HashToPayloadID.Range(func(key, value any) bool {
		if key.(primitives.Slot) < slot {
			p.slotToEth1HashToPayloadID.Delete(key)
		}
		return true
	})
}
