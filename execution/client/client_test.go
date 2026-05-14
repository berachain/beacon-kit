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

package client

import (
	"sync"
	"testing"
)

// TestCapabilitiesRace verifies that concurrent writes (ExchangeCapabilities)
// and reads (HasCapability) on the capabilities map do not race.
func TestCapabilitiesRace(t *testing.T) {
	t.Parallel()

	c := &EngineClient{
		capabilities: make(map[string]struct{}),
	}

	const cap = "engine_newPayloadV3"
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: repeatedly write a capability (simulates ExchangeCapabilities).
	go func() {
		defer wg.Done()
		for range iterations {
			c.capabilitiesMu.Lock()
			c.capabilities[cap] = struct{}{}
			c.capabilitiesMu.Unlock()
		}
	}()

	// Goroutine 2: repeatedly read a capability (simulates HasCapability).
	go func() {
		defer wg.Done()
		for range iterations {
			c.HasCapability(cap)
		}
	}()

	wg.Wait()
}
