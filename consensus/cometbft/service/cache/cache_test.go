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

package cache_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/berachain/beacon-kit/consensus/cometbft/service/cache"
)

const maxCachedStates = 10

// mockTelemetrySink is a simple mock for testing.
type mockTelemetrySink struct {
	lastGaugeKey   string
	lastGaugeValue int64
}

func (m *mockTelemetrySink) SetGauge(key string, value int64, _ ...string) {
	m.lastGaugeKey = key
	m.lastGaugeValue = value
}

func TestCandidateStates_BasicOperations(t *testing.T) {
	t.Parallel()
	sink := &mockTelemetrySink{}
	c := cache.New(maxCachedStates, sink)

	// Test SetCached and GetCached
	elem := &cache.Element{State: &cache.State{}}
	c.SetCached("test-hash", elem)

	got, err := c.GetCached("test-hash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != elem {
		t.Error("retrieved element does not match stored element")
	}

	// Test GetCached for non-existent key
	_, err = c.GetCached("non-existent")
	if !errors.Is(err, cache.ErrStateNotFound) {
		t.Errorf("expected ErrStateNotFound, got: %v", err)
	}

	// Test MarkAsFinal
	err = c.MarkAsFinal("test-hash")
	if err != nil {
		t.Fatalf("unexpected error marking as final: %v", err)
	}

	// Test MarkAsFinal for non-existent key
	err = c.MarkAsFinal("non-existent")
	if !errors.Is(err, cache.ErrFinalizingUnknownState) {
		t.Errorf("expected ErrFinalizingUnknownState, got: %v", err)
	}

	// Test GetFinal
	hash, state, err := c.GetFinal()
	if err != nil {
		t.Fatalf("unexpected error getting final: %v", err)
	}
	if hash != "test-hash" {
		t.Errorf("expected hash 'test-hash', got: %s", hash)
	}
	if state != elem.State {
		t.Error("final state does not match expected state")
	}

	// Test Reset
	c.Reset()

	// Verify the metric was recorded with correct size (1 element was in cache)
	if sink.lastGaugeKey != "beacon_kit.comet.cached_states_size_at_reset" {
		t.Errorf("expected gauge key 'beacon_kit.comet.cached_states_size_at_reset', got: %s", sink.lastGaugeKey)
	}
	if sink.lastGaugeValue != 1 {
		t.Errorf("expected gauge value 1, got: %d", sink.lastGaugeValue)
	}

	_, err = c.GetCached("test-hash")
	if !errors.Is(err, cache.ErrStateNotFound) {
		t.Error("expected cache to be empty after reset")
	}

	_, _, err = c.GetFinal()
	if !errors.Is(err, cache.ErrNoFinalState) {
		t.Errorf("expected ErrNoFinalState after reset, got: %v", err)
	}
}

func TestCandidateStates_BoundedSize(t *testing.T) {
	t.Parallel()
	sink := &mockTelemetrySink{}
	c := cache.New(maxCachedStates, sink)

	// Add more entries than maxCachedStates
	for i := range maxCachedStates + 5 {
		hash := fmt.Sprintf("hash-%d", i)
		c.SetCached(hash, &cache.Element{})
	}

	// Verify oldest entries were evicted (LRU behavior)
	for i := range 5 {
		hash := fmt.Sprintf("hash-%d", i)
		_, err := c.GetCached(hash)
		if !errors.Is(err, cache.ErrStateNotFound) {
			t.Errorf("expected hash-%d to be evicted, but it was found", i)
		}
	}

	// Verify newest entries are still present
	for i := 5; i < maxCachedStates+5; i++ {
		hash := fmt.Sprintf("hash-%d", i)
		_, err := c.GetCached(hash)
		if err != nil {
			t.Errorf("expected hash-%d to be present, got error: %v", i, err)
		}
	}
}
