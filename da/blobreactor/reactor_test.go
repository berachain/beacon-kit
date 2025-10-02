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

package blobreactor_test

import (
	"errors"
	"testing"
	"time"

	"cosmossdk.io/log"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/blobreactor"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/p2p"
	"github.com/stretchr/testify/require"
)

type stubBlobStore struct {
	blobs map[uint64][][]byte // slot -> blob data
	delay time.Duration       // optional delay to simulate slow responses
}

func newStubBlobStore() *stubBlobStore {
	return &stubBlobStore{
		blobs: make(map[uint64][][]byte),
	}
}

func newSlowStubBlobStore(delay time.Duration) *stubBlobStore {
	return &stubBlobStore{
		blobs: make(map[uint64][][]byte),
		delay: delay,
	}
}

func (m *stubBlobStore) Has(_ uint64, _ []byte) (bool, error) {
	return false, nil
}

func (m *stubBlobStore) GetByIndex(index uint64) ([][]byte, error) {
	// Simulate slow response if delay is set
	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	if blobs, ok := m.blobs[index]; ok {
		return blobs, nil
	}

	return nil, errors.New("blobs not found")
}

func (m *stubBlobStore) setBlobs(slot uint64, sidecars datypes.BlobSidecars) error {
	sidecarBzs := make([][]byte, len(sidecars))
	for i, sidecar := range sidecars {
		data, err := sidecar.MarshalSSZ()
		if err != nil {
			return err
		}
		sidecarBzs[i] = data
	}
	m.blobs[slot] = sidecarBzs
	return nil
}

func createTestSidecars(t *testing.T, count int) datypes.BlobSidecars {
	t.Helper()

	sidecars := make([]*datypes.BlobSidecar, count)
	for i := range count {
		sidecars[i] = &datypes.BlobSidecar{
			Index: uint64(i),
			SignedBeaconBlockHeader: &ctypes.SignedBeaconBlockHeader{
				Header: &ctypes.BeaconBlockHeader{},
			},
			InclusionProof: make([]common.Root, ctypes.KZGInclusionProofDepth),
		}
	}

	return sidecars
}

func newTestReactor(t *testing.T, store blobreactor.BlobStore, config blobreactor.Config) *blobreactor.BlobReactor {
	t.Helper()
	logger := log.NewTestLogger(t)
	reactor := blobreactor.NewBlobReactor(store, logger, config)
	return reactor
}

func makeTestP2PConfig(t *testing.T) *config.P2PConfig {
	t.Helper()
	p2pConfig := config.DefaultP2PConfig()
	p2pConfig.ListenAddress = "tcp://127.0.0.1:0" // Use random port
	return p2pConfig
}

func makeConnectedReactors(
	t *testing.T, n int, stores []*stubBlobStore, configs []blobreactor.Config,
) ([]*blobreactor.BlobReactor, []*p2p.Switch) {
	t.Helper()

	tempDir := t.TempDir()

	p2pConfig := makeTestP2PConfig(t)
	p2pConfig.RootDir = tempDir

	reactors := make([]*blobreactor.BlobReactor, n)
	for i := range n {
		reactors[i] = newTestReactor(t, stores[i], configs[i])
	}

	initSwitch := func(i int, sw *p2p.Switch) *p2p.Switch {
		sw.AddReactor(blobreactor.ReactorName, reactors[i])
		return sw
	}
	switches := p2p.MakeConnectedSwitches(p2pConfig, n, initSwitch, p2p.Connect2Switches)

	return reactors, switches
}

func stopSwitches(switches []*p2p.Switch) {
	for _, s := range switches {
		_ = s.Stop()
	}
}

// Test basic connectivity and just request blobs from a single peer
func TestBlobReactor_BasicRequest(t *testing.T) {
	t.Parallel()

	slot := math.Slot(123)
	requestingStore := newStubBlobStore()
	servingStore := newStubBlobStore()

	// Serving store has blobs
	blobs := createTestSidecars(t, 2)
	err := servingStore.setBlobs(slot.Unwrap(), blobs)
	require.NoError(t, err)

	stores := []*stubBlobStore{requestingStore, servingStore}
	configs := []blobreactor.Config{
		{RequestTimeout: 5 * time.Second},
		{RequestTimeout: 5 * time.Second},
	}

	reactors, switches := makeConnectedReactors(t, 2, stores, configs)
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(slot + 10)
	}

	verifier := func(_ datypes.BlobSidecars) error { return nil }

	sidecars, err := reactors[0].RequestBlobs(t.Context(), slot, verifier)

	require.NoError(t, err)
	require.NotNil(t, sidecars)
	require.Len(t, sidecars, 2)
}

// Test peer retry when first peer has no blobs and second peer succeeds
func TestBlobReactor_PeerRetry(t *testing.T) {
	t.Parallel()

	slot := math.Slot(123)
	requestingStore := newStubBlobStore()
	unavailableStore := newStubBlobStore() // Empty - will return error
	validStore := newStubBlobStore()

	// Only valid store has blobs
	validBlobs := createTestSidecars(t, 2)
	err := validStore.setBlobs(slot.Unwrap(), validBlobs)
	require.NoError(t, err)

	stores := []*stubBlobStore{requestingStore, unavailableStore, validStore}
	configs := []blobreactor.Config{
		{RequestTimeout: 5 * time.Second},
		{RequestTimeout: 5 * time.Second},
		{RequestTimeout: 5 * time.Second},
	}

	reactors, switches := makeConnectedReactors(t, 3, stores, configs)
	defer stopSwitches(switches)

	// Set head slots so blobs are considered available
	for _, r := range reactors {
		r.SetHeadSlot(slot + 10)
	}

	// Verifier accepts valid blobs
	verifier := func(sidecars datypes.BlobSidecars) error {
		if len(sidecars) != 2 {
			return errors.New("expected 2 blobs")
		}
		return nil
	}

	sidecars, err := reactors[0].RequestBlobs(t.Context(), slot, verifier)

	// Should succeed despite one peer not having blobs
	require.NoError(t, err)
	require.NotNil(t, sidecars)
	require.Len(t, sidecars, 2)
}

// Test when all peers fail to provide valid blobs
func TestBlobReactor_AllPeersFailed(t *testing.T) {
	t.Parallel()

	slot := math.Slot(456)

	// Create stores: requesting (empty), peer (empty - no blobs)
	requestingStore := newStubBlobStore()
	peerStore := newStubBlobStore()

	stores := []*stubBlobStore{requestingStore, peerStore}
	configs := []blobreactor.Config{
		{RequestTimeout: 500 * time.Millisecond},
		{RequestTimeout: 500 * time.Millisecond},
	}

	reactors, switches := makeConnectedReactors(t, 2, stores, configs)
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(slot + 10)
	}

	verifier := func(_ datypes.BlobSidecars) error {
		return nil
	}

	sidecars, err := reactors[0].RequestBlobs(t.Context(), slot, verifier)

	require.Error(t, err)
	require.ErrorIs(t, err, blobreactor.ErrAllPeersFailed)
	require.Nil(t, sidecars)
}

// Test request timeout when peer responds too slowly
func TestBlobReactor_RequestTimeout(t *testing.T) {
	t.Parallel()

	slot := math.Slot(789)

	requestingStore := newStubBlobStore()
	slowPeerStore := newSlowStubBlobStore(500 * time.Millisecond)

	stores := []*stubBlobStore{requestingStore, slowPeerStore}
	configs := []blobreactor.Config{
		{RequestTimeout: 200 * time.Millisecond},
		{RequestTimeout: 200 * time.Millisecond},
	}

	reactors, switches := makeConnectedReactors(t, 2, stores, configs)
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(slot + 10)
	}

	verifier := func(_ datypes.BlobSidecars) error {
		return nil
	}

	start := time.Now()
	sidecars, err := reactors[0].RequestBlobs(t.Context(), slot, verifier)
	elapsed := time.Since(start)

	require.Error(t, err)
	require.ErrorIs(t, err, blobreactor.ErrAllPeersFailed)
	require.Nil(t, sidecars)

	// Elapsed time should be at least 200ms (request timeout) and less than 500ms
	require.Greater(t, elapsed, 190*time.Millisecond, "Should wait for request timeout")
	require.Less(t, elapsed, 490*time.Millisecond, "Should timeout before peer responds")
}

// Test concurrent requests route responses correctly
func TestBlobReactor_ConcurrentRequests(t *testing.T) {
	t.Parallel()

	requestingStore := newStubBlobStore()
	servingStore := newStubBlobStore()

	// Set up blobs for multiple slots
	slot1, slot2, slot3 := math.Slot(123), math.Slot(234), math.Slot(345)

	blobs1 := createTestSidecars(t, 1)
	err := servingStore.setBlobs(slot1.Unwrap(), blobs1)
	require.NoError(t, err)

	blobs2 := createTestSidecars(t, 1)
	err = servingStore.setBlobs(slot2.Unwrap(), blobs2)
	require.NoError(t, err)

	blobs3 := createTestSidecars(t, 1)
	err = servingStore.setBlobs(slot3.Unwrap(), blobs3)
	require.NoError(t, err)

	stores := []*stubBlobStore{requestingStore, servingStore}
	configs := []blobreactor.Config{
		{RequestTimeout: 2 * time.Second},
		{RequestTimeout: 2 * time.Second},
	}

	reactors, switches := makeConnectedReactors(t, 2, stores, configs)
	defer stopSwitches(switches)

	maxSlot := slot3 + 10
	for _, r := range reactors {
		r.SetHeadSlot(maxSlot)
	}

	verifier := func(_ datypes.BlobSidecars) error {
		return nil
	}

	// Start 3 concurrent requests
	type requestResult struct {
		slot     math.Slot
		sidecars []*datypes.BlobSidecar
		err      error
	}
	results := make(chan requestResult, 3)

	for _, slot := range []math.Slot{slot1, slot2, slot3} {
		slot := slot
		go func() {
			sc, reqErr := reactors[0].RequestBlobs(t.Context(), slot, verifier)
			results <- requestResult{slot, sc, reqErr}
		}()
	}

	// Collect all results
	receivedSlots := make(map[math.Slot]bool)
	for range 3 {
		select {
		case result := <-results:
			require.NoError(t, result.err, "Request for slot %d failed", result.slot)
			require.NotNil(t, result.sidecars)
			receivedSlots[result.slot] = true
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent request results")
		}
	}

	// Verify all requests succeeded with correct slots
	require.True(t, receivedSlots[slot1])
	require.True(t, receivedSlots[slot2])
	require.True(t, receivedSlots[slot3])
}

// Test that verifier correctly rejects then accepts blobs
func TestBlobReactor_VerifierFunctionality(t *testing.T) {
	t.Parallel()

	slot := math.Slot(567)
	requestingStore := newStubBlobStore()
	servingStore := newStubBlobStore()

	validBlobs := createTestSidecars(t, 2)
	err := servingStore.setBlobs(slot.Unwrap(), validBlobs)
	require.NoError(t, err)

	stores := []*stubBlobStore{requestingStore, servingStore}
	configs := []blobreactor.Config{
		{RequestTimeout: 2 * time.Second},
		{RequestTimeout: 2 * time.Second},
	}

	reactors, switches := makeConnectedReactors(t, 2, stores, configs)
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(slot + 10)
	}

	// First request: verifier rejects all blobs
	verifierReject := func(_ datypes.BlobSidecars) error {
		return errors.New("verification failed")
	}

	sidecars, err := reactors[0].RequestBlobs(t.Context(), slot, verifierReject)

	require.Error(t, err)
	require.Nil(t, sidecars)

	// Second request: verifier accepts blobs
	verifierAccept := func(sidecars datypes.BlobSidecars) error {
		if len(sidecars) != 2 {
			return errors.New("expected 2 blobs")
		}
		return nil
	}

	sidecars, err = reactors[0].RequestBlobs(t.Context(), slot, verifierAccept)

	require.NoError(t, err)
	require.NotNil(t, sidecars)
	require.Len(t, sidecars, 2)
}
