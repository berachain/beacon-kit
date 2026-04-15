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

package preconf

import (
	"sync"

	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

// ProposerTracker tracks the expected proposer for the next slot.
type ProposerTracker interface {
	// SetExpectedProposer records the expected proposer for a slot.
	SetExpectedProposer(slot math.Slot, pubkey crypto.BLSPubkey)
	// ClearExpectedProposer clears the expected proposer if the slot matches.
	ClearExpectedProposer(slot math.Slot)
	// IsExpectedProposer returns true if pubkey matches the expected proposer for slot.
	IsExpectedProposer(slot math.Slot, pubkey crypto.BLSPubkey) bool
}

// proposerTracker is the concrete implementation of ProposerTracker.
type proposerTracker struct {
	mu       sync.RWMutex
	slot     math.Slot
	expected crypto.BLSPubkey
}

// NewProposerTracker creates a new ProposerTracker. A freshly constructed
// tracker has no expected proposer, so IsExpectedProposer returns false until
// SetExpectedProposer is called.
func NewProposerTracker() ProposerTracker {
	return &proposerTracker{}
}

func (pt *proposerTracker) SetExpectedProposer(slot math.Slot, pubkey crypto.BLSPubkey) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.slot = slot
	pt.expected = pubkey
}

func (pt *proposerTracker) ClearExpectedProposer(slot math.Slot) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if pt.slot == slot {
		pt.expected = crypto.BLSPubkey{}
	}
}

func (pt *proposerTracker) IsExpectedProposer(slot math.Slot, pubkey crypto.BLSPubkey) bool {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	return pt.slot == slot && pt.expected == pubkey
}
