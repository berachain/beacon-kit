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

	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// RootToPayloadIDMap is a map with keys the head root and values the
// corresponding PayloadID.
type RootToPayloadIDMap map[[32]byte]enginev1.PayloadIDBytes

// PayloadIDCache is a cache that keeps track of the prepared payload ID for the
// given slot and with the given head root.
type PayloadIDCache struct {
	slotToPayloadID map[primitives.Slot]RootToPayloadIDMap
	sync.Mutex
}

// NewPayloadIDCache returns a new payload ID cache.
func NewPayloadIDCache() *PayloadIDCache {
	return &PayloadIDCache{slotToPayloadID: make(map[primitives.Slot]RootToPayloadIDMap)}
}

// PayloadID returns the payload ID for the given slot and parent block root.
func (p *PayloadIDCache) PayloadID(
	slot primitives.Slot, root [32]byte,
) (enginev1.PayloadIDBytes, bool) {
	p.Lock()
	defer p.Unlock()
	inner, ok := p.slotToPayloadID[slot]
	if !ok {
		return enginev1.PayloadIDBytes{}, false
	}
	pid, ok := inner[root]
	if !ok {
		return enginev1.PayloadIDBytes{}, false
	}
	return pid, true
}

// SetPayloadID updates the payload ID for the given slot and head root
// it also prunes older entries in the cache.
func (p *PayloadIDCache) Set(slot primitives.Slot, root [32]byte, pid enginev1.PayloadIDBytes) {
	p.Lock()
	defer p.Unlock()
	if slot > 1 {
		p.prune(slot - 2) //nolint:gomnd // from prysm.
	}
	inner, ok := p.slotToPayloadID[slot]
	if !ok {
		inner = make(RootToPayloadIDMap)
		p.slotToPayloadID[slot] = inner
	}
	inner[root] = pid
}

// Prune prunes old payload IDs. Requires a Lock in the cache.
func (p *PayloadIDCache) prune(slot primitives.Slot) {
	for key := range p.slotToPayloadID {
		if key < slot {
			delete(p.slotToPayloadID, key)
		}
	}
}
