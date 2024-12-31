// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.

package chain

import (
	"sync/atomic"
)

// heightTracker tracks the current block height.
type heightTracker struct {
	currentHeight atomic.Uint64
}

// newHeightTracker creates a new instance of heightTracker
func newHeightTracker() *heightTracker {
	return &heightTracker{}
}

// GetCurrentHeight returns the current block height
func (ht *heightTracker) GetCurrentHeight() uint64 {
	return ht.currentHeight.Load()
}

// SetCurrentHeight sets the current block height
func (ht *heightTracker) SetCurrentHeight(height uint64) {
	ht.currentHeight.Store(height)
} 
